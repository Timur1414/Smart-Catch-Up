package usecase

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/brianvoe/gofakeit/v7"
)

func GenerateNotification(notificationType string, text string, recipientId int, actors []domain.NotificationActor) domain.Notification {
	actor := actors[rand.Intn(len(actors))]
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
		ActorId:          actor.Id,
		ActorName:        actor.Name,
		ObjectId:         objectId,
		ObjectType:       "post",
		CreatedAt:        createdAt,
		ReadAt:           time.Time{},
		Payload:          payload,
		Actions:          nil,
	}
}

type PayloadGenerator func(notificationType string, objectId int) string

var payloadStrategies = map[string]PayloadGenerator{
	"order": func(notifType string, objectId int) string {
		templates := []string{
			"Заказ #{id} передан в службу доставки ({city})",
			"Статус заказа #{id} изменен: «В пути»",
			"Товар «{product}» из заказа #{id} готов к выдаче",
		}
		tpl := pickRandom(templates)
		res, err := gofakeit.Generate(tpl)
		if err != nil {
			return ""
		}
		return strings.ReplaceAll(res, "{id}", fmt.Sprintf("%d", objectId))
	},

	"payment": func(notifType string, objectId int) string {
		templates := []string{
			"Перевод на сумму {price:300,10000} ₽ успешно получен",
			"С вашего счета списано {price:100,2000} ₽ за оплату подписки",
			"Баланс пополнен на {price:500,5000} ₽ через СБП",
		}
		res, err := gofakeit.Generate(pickRandom(templates))
		if err != nil {
			return ""
		}
		return res
	},

	"security": func(notifType string, objectId int) string {
		templates := []string{
			"Выполнен новый вход с устройства ({appname}, г. {city}, IP: {ipv4address})",
			"Пароль от аккаунта был успешно изменен",
			"Попытка входа с подозрительного IP-адреса {ipv4address} заблокирована",
		}
		res, err := gofakeit.Generate(pickRandom(templates))
		if err != nil {
			return ""
		}
		return res
	},

	"comment": func(notifType string, objectId int) string {
		return fmt.Sprintf("оставил комментарий: «%s»", gofakeit.Comment())
	},

	"reaction": func(notifType string, objectId int) string {
		reactions := []string{
			"поставил отметку «Нравится» вашей публикации",
			"оценил вашу фотографию 🔥",
			"поделился вашей записью на своей странице",
		}
		return pickRandom(reactions)
	},

	"friend": func(notifType string, objectId int) string {
		actions := []string{
			"отправил вам запрос на добавление в друзья",
			"принял вашу заявку в друзья",
			"подписался на ваши обновления",
			"празднует сегодня день рождения 🎉. Не забудьте поздравить!",
		}
		return pickRandom(actions)
	},

	"media": func(notifType string, objectId int) string {
		templates := []string{
			"опубликовал новый пост: «{buzzword} в современной разработке»",
			"загрузил новое видео в альбом",
			"запустил прямую трансляцию: «Обсуждаем {buzzword}»",
		}
		res, err := gofakeit.Generate(pickRandom(templates))
		if err != nil {
			return ""
		}
		return res
	},
}

func GeneratePayloadByType(notificationType string, objectId int) string {
	category := resolveCategory(notificationType)
	generator, exists := payloadStrategies[category]
	if !exists {
		return gofakeit.Sentence(6)
	}
	return generator(notificationType, objectId)
}

func resolveCategory(t string) string {
	switch {
	case strings.Contains(t, "order"), strings.Contains(t, "cart"), strings.Contains(t, "delivery"), strings.Contains(t, "product"):
		return "order"
	case strings.Contains(t, "payment"), strings.Contains(t, "money"), strings.Contains(t, "balance"), strings.Contains(t, "subscription"):
		return "payment"
	case strings.Contains(t, "security"), strings.Contains(t, "access"), strings.Contains(t, "account"), strings.Contains(t, "verification"):
		return "security"
	case strings.Contains(t, "comment"), strings.Contains(t, "reply"):
		return "comment"
	case strings.Contains(t, "like"), strings.Contains(t, "reaction"), strings.Contains(t, "repost"):
		return "reaction"
	case strings.Contains(t, "friend"), strings.Contains(t, "subscriber"), strings.Contains(t, "birthday"):
		return "friend"
	case strings.Contains(t, "post"), strings.Contains(t, "video"), strings.Contains(t, "stream"), strings.Contains(t, "photo"):
		return "media"
	default:
		return "default"
	}
}

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}
