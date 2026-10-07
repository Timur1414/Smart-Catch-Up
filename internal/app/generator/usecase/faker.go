package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/brianvoe/gofakeit/v7"
)

var actors = []struct {
	Id   int
	Name string
}{
	{Id: 1, Name: "System"},
	{Id: 2, Name: "Алексей Смирнов"},
	{Id: 3, Name: "Екатерина Кузнецова"},
	{Id: 4, Name: "Михаил Соколов"},
	{Id: 5, Name: "Мария Козлова"},
	{Id: 6, Name: "Иван Федоров"},
	{Id: 7, Name: "Анна Васильева"},
	{Id: 8, Name: "Сергей Лебедев"},
	{Id: 9, Name: "Ольга Новикова"},
	{Id: 10, Name: "Артем Морозов"},
	{Id: 11, Name: "Полина Волкова"},
}

func GenerateNotification(notificationType string, text string, recipientId int, userIds []int) domain.Notification {
	actor := actors[rand.Intn(len(actors))]
	actorName := actor.Name
	actorId := 1
	// TODO: switch rand actor to rand userId from request & db
	now := time.Now()
	createdAt := gofakeit.DateRange(now.Add(-5*24*time.Hour), now)
	objectId := gofakeit.Number(100, 10000)
	payload := text
	if payload == "" {
		switch notificationType {
		case "comment_on_post", "reply_to_comment":
			payload = fmt.Sprintf("оставил комментарий: «%s»", gofakeit.Sentence(gofakeit.Number(4, 10)))
		case "like", "reaction_on_post":
			payload = "поставил отметку «Нравится» вашей публикации"
		case "order_status_changed":
			payload = fmt.Sprintf("Статус заказа #%d: «%s»", objectId, gofakeit.Word())
		case "money_transfer":
			payload = fmt.Sprintf("Перевод на сумму %.0f ₽ получен", gofakeit.Price(100, 5000))
		default:
			payload = gofakeit.Sentence(6)
		}
	}

	return domain.Notification{
		NotificationType: notificationType,
		RecipientId:      recipientId,
		ActorId:          actorId,
		ActorName:        actorName,
		ObjectId:         objectId,
		ObjectType:       "post",
		CreatedAt:        createdAt,
		ReadAt:           time.Time{},
		Payload:          payload,
		Actions:          nil,
	}
}
