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
		payload = GeneratePayloadByType(notificationType, objectId)
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
			"Товар «{productname}» из заказа #{id} готов к выдаче",
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
		return fmt.Sprintf("оставил комментарий: «%s»", generateComment())
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

	"default": func(notifType string, objectId int) string {
		return generateComment()
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

var russianComments = []string{
	"Отличный пост, спасибо за полезную информацию!",
	"Полностью согласен с автором 👍",
	"А есть примеры кода или ссылка на репозиторий?",
	"Интересная мысль, но на практике часто возникают нюансы.",
	"Очень вовремя, как раз сейчас разбираюсь с этой темой 🔥",
	"Не совсем понял один момент, можно подробнее?",
	"Круто расписано, сохранил себе в закладки 📌",
	"Давно искал внятное объяснение, спасибо большое!",
	"Жду продолжения, очень интересно 🚀",
	"Спорное утверждение, но аргументация хорошая.",
	"Это база 💯",
	"Полезно, переслал коллегам в рабочий чат.",
	"А как это решение покажет себя при высоких нагрузках?",
	"Поддерживаю! У нас на проекте была аналогичная ситуация.",
	"Супер, всё четко и по делу.",
}

func generateComment() string {
	return pickRandom(russianComments)
}

var russianCities = []string{
	"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань",
	"Нижний Новгород", "Челябинск", "Красноярск", "Самара", "Уфа",
	"Ростов-на-Дону", "Омск", "Краснодар", "Воронеж", "Пермь", "Волгоград",
}

var russianBuzzwords = []string{
	"микросервисы", "нейросети и AI", "высокие нагрузки",
	"чистую архитектуру", "Kubernetes", "DevOps-практики",
	"Go против Rust", "безопасность API", "рефакторинг легаси",
}

func init() {
	gofakeit.AddFuncLookup("city", gofakeit.Info{
		Display:     "City",
		Category:    "address",
		Description: "Случайный российский город",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			return pickRandom(russianCities), nil
		},
	})
	gofakeit.AddFuncLookup("buzzword", gofakeit.Info{
		Display:     "Buzzword",
		Category:    "word",
		Description: "IT-термины на русском",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			return pickRandom(russianBuzzwords), nil
		},
	})
}
