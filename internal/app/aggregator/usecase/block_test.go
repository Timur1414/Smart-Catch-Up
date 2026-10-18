package usecase_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/usecase"
)

type mockBlockRepo struct {
	createdBlocks []domain.Block
}

func (m *mockBlockRepo) Create(ctx context.Context, block domain.Block) (int, error) {
	m.createdBlocks = append(m.createdBlocks, block)
	return len(m.createdBlocks), nil
}

type mockNotificationRepo struct {
	notifications []domain.Notification
}

func (m *mockNotificationRepo) GetAllByUserRange(ctx context.Context, startUserID, endUserID int) ([]domain.Notification, error) {
	var result []domain.Notification
	for _, n := range m.notifications {
		if n.RecipientId >= startUserID && n.RecipientId < endUserID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (m *mockNotificationRepo) GetById(ctx context.Context, id int) (domain.Notification, error) {
	return domain.Notification{}, nil
}

func (m *mockNotificationRepo) GetByUser(ctx context.Context, userId int) ([]domain.Notification, error) {
	return nil, nil
}

type mockDigestRepo struct {
	updatedBlocks map[int]int
}

func (m *mockDigestRepo) UpdateBlock(ctx context.Context, userId, blockId int) error {
	if m.updatedBlocks == nil {
		m.updatedBlocks = make(map[int]int)
	}
	m.updatedBlocks[userId] = blockId
	return nil
}

type mockSettingsRepo struct {
	settingsByUser map[int]domain.Settings
}

func (m *mockSettingsRepo) GetByUser(ctx context.Context, id int) (domain.Settings, error) {
	if s, ok := m.settingsByUser[id]; ok {
		return s, nil
	}
	return domain.Settings{UserId: id}, nil
}

func (m *mockSettingsRepo) GetByUsers(ctx context.Context, ids []int) (map[int]domain.Settings, error) {
	res := make(map[int]domain.Settings, len(ids))
	for _, id := range ids {
		if s, ok := m.settingsByUser[id]; ok {
			res[id] = s
		} else {
			res[id] = domain.Settings{UserId: id}
		}
	}
	return res, nil
}

func TestAggregate_WithImportantNotifications(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	blockRepo := &mockBlockRepo{}
	digestRepo := &mockDigestRepo{}
	settingsRepo := &mockSettingsRepo{
		settingsByUser: map[int]domain.Settings{
			1: {
				UserId:    1,
				Important: []string{"friend_request", "security_warning"},
			},
		},
	}

	notifRepo := &mockNotificationRepo{
		notifications: []domain.Notification{
			{
				Id:               1,
				RecipientId:      1,
				NotificationType: "like",
				ActorName:        "Alice",
				Payload:          "понравилась ваша запись",
				CreatedAt:        now.Add(-10 * time.Minute),
			},
			{
				Id:               2,
				RecipientId:      1,
				NotificationType: "friend_request",
				ActorName:        "Bob",
				Payload:          "отправил(-а) вам заявку в друзья",
				CreatedAt:        now.Add(-5 * time.Minute),
			},
			{
				Id:               3,
				RecipientId:      1,
				NotificationType: "security_warning",
				ActorName:        "System",
				Payload:          "Вход с нового устройства",
				CreatedAt:        now.Add(-1 * time.Minute),
			},
		},
	}

	uc := usecase.NewBlock(blockRepo, notifRepo, digestRepo, settingsRepo)
	err := uc.Aggregate(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blockRepo.createdBlocks) != 1 {
		t.Fatalf("expected 1 block created, got %d", len(blockRepo.createdBlocks))
	}

	block := blockRepo.createdBlocks[0]
	if block.UserId != 1 {
		t.Errorf("expected block user_id 1, got %d", block.UserId)
	}

	// Should have 2 parts: 1 for content_reactions cluster ("like"), 1 for "important"
	if len(block.Parts) != 2 {
		t.Fatalf("expected 2 block parts, got %d", len(block.Parts))
	}

	var importantPart *domain.BlockPart
	var clusterPart *domain.BlockPart

	for i := range block.Parts {
		if block.Parts[i].ClusterType == "important" {
			importantPart = &block.Parts[i]
		} else if block.Parts[i].ClusterType == "content_reactions" {
			clusterPart = &block.Parts[i]
		}
	}

	if importantPart == nil {
		t.Fatal("expected an 'important' block part, but none found")
	}
	if clusterPart == nil {
		t.Fatal("expected a 'content_reactions' block part, but none found")
	}

	// Check important part content: original payloads only (not ActorName: Payload)
	expectedImportantText := "отправил(-а) вам заявку в друзья\nВход с нового устройства"
	if importantPart.Text != expectedImportantText {
		t.Errorf("expected important text %q, got %q", expectedImportantText, importantPart.Text)
	}

	if len(importantPart.NotificationIds) != 2 || importantPart.NotificationIds[0] != 2 || importantPart.NotificationIds[1] != 3 {
		t.Errorf("expected important notification IDs [2, 3], got %v", importantPart.NotificationIds)
	}

	expectedImpStart := now.Add(-5 * time.Minute)
	expectedImpEnd := now.Add(-1 * time.Minute)
	if !importantPart.Start.Equal(expectedImpStart) || !importantPart.End.Equal(expectedImpEnd) {
		t.Errorf("expected important time range [%v, %v], got [%v, %v]", expectedImpStart, expectedImpEnd, importantPart.Start, importantPart.End)
	}

	// Check cluster part content: formatted with ActorName: Payload
	expectedClusterText := "Alice: понравилась ваша запись"
	if clusterPart.Text != expectedClusterText {
		t.Errorf("expected cluster text %q, got %q", expectedClusterText, clusterPart.Text)
	}

	// Digest update checked
	if digestRepo.updatedBlocks[1] != 1 {
		t.Errorf("expected digest updated with blockId 1, got %d", digestRepo.updatedBlocks[1])
	}
}

func TestAggregate_NoImportantNotifications(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	blockRepo := &mockBlockRepo{}
	digestRepo := &mockDigestRepo{}
	settingsRepo := &mockSettingsRepo{
		settingsByUser: map[int]domain.Settings{
			1: {
				UserId:    1,
				Important: []string{"security_warning"},
			},
		},
	}

	notifRepo := &mockNotificationRepo{
		notifications: []domain.Notification{
			{
				Id:               1,
				RecipientId:      1,
				NotificationType: "like",
				ActorName:        "Alice",
				Payload:          "понравилась ваша запись",
				CreatedAt:        now,
			},
		},
	}

	uc := usecase.NewBlock(blockRepo, notifRepo, digestRepo, settingsRepo)
	err := uc.Aggregate(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blockRepo.createdBlocks) != 1 {
		t.Fatalf("expected 1 block created, got %d", len(blockRepo.createdBlocks))
	}

	block := blockRepo.createdBlocks[0]
	for _, part := range block.Parts {
		if part.ClusterType == "important" {
			t.Errorf("did not expect important block part, but found one")
		}
	}
}

func TestAggregate_OnlyImportantNotifications(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	blockRepo := &mockBlockRepo{}
	digestRepo := &mockDigestRepo{}
	settingsRepo := &mockSettingsRepo{
		settingsByUser: map[int]domain.Settings{
			1: {
				UserId:    1,
				Important: []string{"friend_request"},
			},
		},
	}

	notifRepo := &mockNotificationRepo{
		notifications: []domain.Notification{
			{
				Id:               1,
				RecipientId:      1,
				NotificationType: "friend_request",
				ActorName:        "Bob",
				Payload:          "отправил заявку",
				CreatedAt:        now,
			},
		},
	}

	uc := usecase.NewBlock(blockRepo, notifRepo, digestRepo, settingsRepo)
	err := uc.Aggregate(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blockRepo.createdBlocks) != 1 {
		t.Fatalf("expected 1 block created, got %d", len(blockRepo.createdBlocks))
	}

	block := blockRepo.createdBlocks[0]
	if len(block.Parts) != 1 {
		t.Fatalf("expected exactly 1 block part, got %d", len(block.Parts))
	}
	if block.Parts[0].ClusterType != "important" {
		t.Errorf("expected cluster 'important', got %q", block.Parts[0].ClusterType)
	}
	if block.Parts[0].Text != "отправил заявку" {
		t.Errorf("expected original payload 'отправил заявку', got %q", block.Parts[0].Text)
	}
	if strings.Contains(block.Parts[0].Text, "Bob:") {
		t.Errorf("important notification should be in original form without ActorName prefix")
	}
}
