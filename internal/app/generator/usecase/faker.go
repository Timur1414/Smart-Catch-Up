package usecase

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/domain"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"
)

func registerLookup(name, display, category, desc string, list []string) {
	gofakeit.AddFuncLookup(name, gofakeit.Info{
		Display:     display,
		Category:    category,
		Description: desc,
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			return pickRandom(list), nil
		},
	})
}

func init() {
	registerLookup("city", "City", "address", "Случайный российский город", russianCities)
	registerLookup("buzzword", "Buzzword", "word", "IT-термины на русском", russianBuzzwords)
	registerLookup("productname", "ProductName", "product", "Названия товаров", russianProducts)
	registerLookup("community", "Community", "social", "Названия сообществ", russianCommunities)
	registerLookup("track_title", "TrackTitle", "media", "Названия треков", russianTracks)
	registerLookup("playlist_title", "PlaylistTitle", "media", "Названия плейлистов", russianPlaylists)
	registerLookup("album_title", "AlbumTitle", "media", "Названия альбомов", russianAlbums)
	registerLookup("podcast_topic", "PodcastTopic", "media", "Темы подкастов", russianPodcasts)
	registerLookup("event_title", "EventTitle", "event", "Названия мероприятий", russianEvents)
	registerLookup("order_status", "OrderStatus", "shop", "Статусы заказа", russianOrderStatuses)
	registerLookup("subscription_plan", "SubscriptionPlan", "payment", "Тарифные планы подписок", russianSubscriptionPlans)
	registerLookup("reaction_emoji", "ReactionEmoji", "social", "Эмодзи реакций", russianReactions)
	registerLookup("gift_name", "GiftName", "social", "Названия подарков", russianGifts)
	registerLookup("gift_wish", "GiftWish", "social", "Пожелания к подаркам", russianGiftWishes)
	registerLookup("community_role", "CommunityRole", "social", "Роли в сообществе", russianCommunityRoles)
	registerLookup("street_name", "StreetName", "address", "Улицы российских городов", russianStreets)
	registerLookup("transfer_comment", "TransferComment", "payment", "Комментарии к переводам", russianTransferComments)
	registerLookup("achievement_title", "AchievementTitle", "game", "Названия игровых достижений", russianAchievements)
	registerLookup("game_title", "GameTitle", "game", "Названия игр", russianGames)
	registerLookup("support_answer", "SupportAnswer", "support", "Ответы службы поддержки", russianSupportReplies)
	registerLookup("review_text", "ReviewText", "shop", "Тексты отзывов", russianReviews)
	registerLookup("question_text", "QuestionText", "social", "Вопросы для историй", russianQuestions)
	registerLookup("answer_text", "AnswerText", "social", "Ответы на вопросы историй", russianAnswers)
	registerLookup("comment_text", "CommentText", "social", "Комментарии пользователей", russianComments)
}

func GenerateNotification(notificationType string, text string, recipientId int, actors []domain.NotificationActor) domain.Notification {
	actor := actors[rand.Intn(len(actors))]
	systemTypes := []string{
		"account_changed", "data_verification", "phone_linked", "phone_unlinked",
		"email_linked", "email_unlinked", "security_warning", "access_restricted", "access_restored",
		"new_feature_offer", "service_notification", "app_notification", "support_reply", "payment_error",
		"subscription_renewed", "subscription_purchased", "payment_state_changed", "balance_changed",
		"money_transfer",
	}
	if slices.Contains(systemTypes, notificationType) {
		actor = actors[0]
	}
	now := time.Now().UTC()
	createdAt := gofakeit.DateRange(now.Add(-5*24*time.Hour), now)
	log := logger.GetLogger()
	log.Debug("generated time", zap.Time("created_at", createdAt))
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

var exactPayloadTemplates = map[string][]string{
	"friend_request": {
		"отправил(-а) вам заявку в друзья",
		"хочет добавить вас в друзья",
		"хочет стать вашим другом",
	},
	"friend_request_accepted": {
		"принял(-а) вашу заявку в друзья",
		"подтвердил(-а) ваш запрос на добавление в друзья",
		"теперь у вас в друзьях",
	},
	"new_subscriber": {
		"подписался(-ась) на ваши обновления",
		"стал(-а) вашим новым подписчиком",
		"теперь подписан(-а) на вашу страницу",
	},
	"friend_invitation": {
		"приглашает вас дружить и обмениваться новостями",
		"отправил(-а) вам приглашение в друзья",
	},
	"found_acquaintance": {
		"Ваш знакомый из контактов зарегистрировался в приложении",
		"Найден знакомый из вашей телефонной книги",
		"Возможно, вы знакомы: пользователь есть в ваших контактах",
	},
	"recommended_acquaintance": {
		"Возможно, вы знакомы: у вас {number:2,15} общих друзей",
		"Рекомендация: вы можете знать этого пользователя",
		"Пользователь может быть вашим знакомым",
	},
	"friend_birthday": {
		"празднует сегодня день рождения 🎉. Не забудьте поздравить!",
		"сегодня отмечает день рождения 🎂. Напишите поздравление!",
	},
	"gift": {
		"отправил(-а) вам подарок: «{gift_name}» 🎁",
		"сделал(-а) вам подарок «{gift_name}» с пожеланием: «{gift_wish}»",
	},
	"like": {
		"поставил(-а) отметку «Нравится» вашей записи",
		"оценил(-а) вашу публикацию 👍",
	},
	"reaction_on_post": {
		"поставил(-а) реакцию «{reaction_emoji}» на вашу публикацию",
		"отреагировал(-а) на ваш пост: {reaction_emoji}",
		"оставил(-а) реакцию на вашу запись",
	},
	"reaction_on_photo": {
		"оценил(-а) вашу фотографию: ❤️",
		"поставил(-а) реакцию «{reaction_emoji}» на ваше фото",
	},
	"reaction_on_video": {
		"оценил(-а) ваше видео: 👏",
		"поставил(-а) реакцию «{reaction_emoji}» на ваш видеоролик",
	},
	"reaction_on_story": {
		"отреагировал(-а) на вашу историю: 🔥",
		"отправил(-а) реакцию «{reaction_emoji}» на вашу историю",
	},
	"reaction_on_comment": {
		"оценил(-а) ваш комментарий: «{comment_text}»",
		"поставил(-а) лайк на ваш комментарий",
	},
	"reaction_on_product": {
		"добавил(-а) ваш товар «{productname}» в избранное ❤️",
		"оценил(-а) карточку товара «{productname}»",
	},
	"repost": {
		"поделился(-ась) вашей записью на своей странице",
		"сделал(-а) репост вашей публикации",
	},
	"comment_on_post": {
		"прокомментировал(-а) ваш пост: «{comment_text}»",
		"оставил(-а) комментарий к вашей записи: «{comment_text}»",
	},
	"comment_on_photo": {
		"прокомментировал(-а) вашу фотографию: «{comment_text}»",
		"оставил(-а) комментарий под вашим фото: «{comment_text}»",
	},
	"comment_on_video": {
		"оставил(-а) комментарий к вашему видео: «{comment_text}»",
		"прокомментировал(-а) видеоролик: «{comment_text}»",
	},
	"comment_on_article": {
		"оставил(-а) отзыв к статье: «{comment_text}»",
		"прокомментировал(-а) вашу статью про {buzzword}: «{comment_text}»",
	},
	"comment_on_product": {
		"задал(-а) вопрос о товаре «{productname}»: «{comment_text}»",
		"оставил(-а) комментарий к товару «{productname}»",
	},
	"reply_to_comment": {
		"ответил(-а) на ваш комментарий: «{comment_text}»",
		"написал(-а) ответ в ветке обсуждения: «{comment_text}»",
	},
	"author_mention": {
		"упомянул(-а) вас как автора в публикации",
		"указал(-а) ваше авторство в посте про {buzzword}",
	},
	"user_mention": {
		"упомянул(-а) вас в своем посте",
		"отметил(-а) вас в новой публикации",
	},
	"photo_tag": {
		"отметил(-а) вас на фотографии",
		"добавил(-а) отметку с вами на новом фото",
	},
	"post_mention": {
		"упомянул(-а) вас в тексте публикации",
		"добавил(-а) ссылку на ваш профиль в посте про {buzzword}",
	},
	"comment_mention": {
		"упомянул(-а) вас в комментарии: «{comment_text}»",
		"отметил(-а) вас в обсуждении ветки комментариев",
	},
	"story_mention": {
		"отметил(-а) вас в своей истории",
		"упомянул(-а) ваш аккаунт в новой истории",
	},
	"discussion_mention": {
		"упомянул(-а) вас в теме обсуждения «{buzzword}»",
		"обратился(-ась) к вам в ветке обсуждения сообщества",
	},
	"story_question": {
		"задал(-а) вам вопрос в истории: «{question_text}»",
		"отправил(-а) вопрос через стикер в истории: «{question_text}»",
	},
	"story_answer": {
		"ответил(-а) на ваш вопрос в истории: «{answer_text}»",
		"опубликовал(-а) ответ на ваш вопрос в новой истории",
	},
	"new_post": {
		"опубликовал(-а) новый пост: «{buzzword} в современной разработке»",
		"выложил(-а) новую запись: «Как применять {buzzword} на практике»",
	},
	"new_photo": {
		"опубликовал(-а) новую фотографию",
		"загрузил(-а) новые фотографии в альбом «{album_title}»",
	},
	"new_story": {
		"опубликовал(-а) новую историю",
		"выложил(-а) свежую историю — посмотрите, пока не исчезла!",
	},
	"new_video": {
		"загрузил(-а) новое видео: «{buzzword}: подробный разбор»",
		"опубликовал(-а) новый видеоролик длительностью {number:2,45} мин.",
	},
	"new_clip": {
		"опубликовал(-а) новый клип 🎬",
		"выложил(-а) короткое видео в ленту клипов",
	},
	"new_audio": {
		"добавил(-а) новый трек «{track_title}»",
		"поделился(-ась) аудиозаписью «{track_title}»",
	},
	"new_album": {
		"создал(-а) новый фотоальбом «{album_title}»",
		"опубликовал(-а) музыкальный альбом «{album_title}»",
	},
	"new_podcast": {
		"выпустил(-а) новый выпуск подкаста: «{podcast_topic}»",
		"опубликовал(-а) выпуск #{number:1,50} подкаста про {buzzword}",
	},
	"new_playlist": {
		"собрал(-а) новый плейлист «{playlist_title}» ({number:10,30} треков)",
		"поделился(-ась) плейлистом «{playlist_title}»",
	},
	"stream_started": {
		"запустил(-а) прямую трансляцию: «Обсуждаем {buzzword}» 🔴",
		"ведет прямой эфир прямо сейчас — присоединяйтесь к стриму!",
	},
	"new_author_content": {
		"Автор опубликовал эксклюзивный материал: «{buzzword}»",
		"Вышел новый закрытый пост для подписчиков от автора",
	},
	"new_community_content": {
		"В сообществе «{community}» опубликована новая запись: «{buzzword}»",
		"Новая публикация в группе «{community}»",
	},
	"community_invitation": {
		"приглашает вас вступить в сообщество «{community}»",
		"отправил(-а) вам приглашение в группу «{community}»",
	},
	"join_request_accepted": {
		"Ваша заявка на вступление в сообщество «{community}» одобрена",
		"Вы стали участником группы «{community}»",
	},
	"role_changed": {
		"Вам назначена роль «{community_role}» в сообществе «{community}»",
		"Ваши права администратора в сообществе «{community}» обновлены",
	},
	"co_owner_status_changed": {
		"Вы назначены совладельцем сообщества «{community}»",
		"Статус управления сообществом «{community}» был изменен",
	},
	"new_community_subscriber": {
		"Новый подписчик вступил в ваше сообщество «{community}»",
		"В вашей группе «{community}» прибавился участник",
	},
	"suggested_post": {
		"предложил(-а) новость в сообщество «{community}»: «{buzzword}»",
		"В предложенных записях сообщества «{community}» появился новый пост",
	},
	"community_comment": {
		"Сообщество «{community}» оставило комментарий: «{comment_text}»",
		"От имени сообщества «{community}» опубликован ответ: «{comment_text}»",
	},
	"community_reaction": {
		"Сообщество «{community}» оценило вашу запись 🔥",
		"Сообщество «{community}» поставило отметку «Нравится» вашей публикации",
	},
	"product_comment": {
		"задал(-а) вопрос к товару «{productname}»: «{comment_text}»",
		"Новый комментарий в карточке вашего товара «{productname}»",
	},
	"product_review": {
		"оставил(-а) отзыв ({number:3,5}★) на товар «{productname}»: «{review_text}»",
		"Покупатель оценил товар «{productname}»: «{review_text}»",
	},
	"review_published": {
		"Ваш отзыв на товар «{productname}» успешно прошел модерацию и опубликован",
		"Отзыв о товаре «{productname}» опубликован на сайте",
	},
	"review_rejected": {
		"Ваш отзыв на товар «{productname}» отклонен: нарушение правил публикации",
		"Отзыв о товаре «{productname}» не прошел модерацию",
	},
	"order_status_changed": {
		"Статус заказа #{id} изменен: «{order_status}»",
		"Заказ #{id} передан в службу доставки ({city})",
		"Товар «{productname}» из заказа #{id} готов к выдаче",
	},
	"abandoned_cart_reminder": {
		"Вы забыли товары в корзине! Оформите заказ #{id}, пока они есть в наличии",
		"Товар «{productname}» в вашей корзине ждет оформления. Скидка действует еще 24 часа ⏳",
	},
	"delivery_info": {
		"Курьер доставит заказ #{id} сегодня по адресу: г. {city}, ул. {street_name}",
		"Заказ #{id} ожидает в пункте выдачи: г. {city}, ул. {street_name}",
	},
	"money_request": {
		"запросил(-а) у вас перевод на сумму {price:100,5000} ₽",
		"отправил(-а) запрос средств на сумму {price:200,3000} ₽ («{transfer_comment}»)",
	},
	"money_transfer": {
		"Перевод на сумму {price:300,10000} ₽ успешно получен",
		"Вам поступил перевод на сумму {price:500,15000} ₽ («{transfer_comment}»)",
	},
	"balance_changed": {
		"Баланс пополнен на {price:500,5000} ₽ через СБП",
		"Баланс кошелька изменился: текущий остаток {price:1000,50000} ₽",
	},
	"payment_state_changed": {
		"Платеж #{id} на сумму {price:200,5000} ₽ успешно проведен",
		"Статус оплаты по заказу #{id} обновлен: подтверждено банком",
	},
	"subscription_purchased": {
		"Оформлена подписка «{subscription_plan}»",
		"С вашего счета списано {price:199,999} ₽ за оплату подписки «{subscription_plan}»",
	},
	"subscription_renewed": {
		"Подписка «{subscription_plan}» успешно продлена на следующий период",
		"Автопродление подписки «{subscription_plan}» выполнено ({price:199,999} ₽)",
	},
	"payment_error": {
		"Ошибка оплаты: не удалось списать {price:199,2990} ₽ по карте *{number:1000,9999}: недостаточно средств",
		"Не удалось провести платеж #{id}. Пожалуйста, проверьте реквизиты карты",
	},
	"account_changed": {
		"Данные вашего аккаунта были успешно обновлены",
		"Пароль от аккаунта был успешно изменен",
	},
	"data_verification": {
		"Ваш профиль успешно прошел верификацию ✓",
		"Подтверждение данных завершено: получен значок проверенного аккаунта",
	},
	"phone_linked": {
		"К аккаунту успешно привязан новый номер телефона: +7 (9**) ***-{number:10,99}-{number:10,99}",
		"Номер телефона подтвержден и сохранен в профиле",
	},
	"phone_unlinked": {
		"Номер телефона был отвязан от вашего аккаунта",
		"Резервный номер телефона успешно удален из профиля",
	},
	"email_linked": {
		"Email-адрес {email} успешно привязан к вашему аккаунту",
		"Электронная почта {email} подтверждена",
	},
	"email_unlinked": {
		"Email-адрес был отвязан от вашего аккаунта",
		"Адрес электронной почты удален из способов связи",
	},
	"security_warning": {
		"Выполнен новый вход с устройства ({appname}, г. {city}, IP: {ipv4address})",
		"Попытка входа с подозрительного IP-адреса {ipv4address} заблокирована",
	},
	"access_restricted": {
		"Действие некоторых функций аккаунта временно ограничено из-за подозрительной активности",
		"Доступ к отправке сообщений временно заблокирован",
	},
	"access_restored": {
		"Все ограничения сняты. Полный доступ к аккаунту восстановлен",
		"Безопасность профиля подтверждена, ограничения доступа сняты",
	},
	"event_starting_soon": {
		"Мероприятие «{event_title}» начнется сегодня в {number:12,20}:00",
		"Напоминание: онлайн-событие «{event_title}» стартует через 15 минут",
	},
	"call_starting_soon": {
		"Запланированный групповой звонок начнется через {number:5,15} минут",
		"Скоро начнется созвон: «{buzzword}»",
	},
	"birthday_reminder": {
		"У вашего друга завтра день рождения! Пора выбрать подарок 🎁",
		"Напоминание: скоро день рождения у близкого друга",
	},
	"saved_wish_reminder": {
		"Товар из вашего списка желаний «{productname}» стал доступен со скидкой {number:10,40}%!",
		"Напоминание: вы сохраняли публикацию про {buzzword} в закладки",
	},
	"people_recommendation": {
		"Возможно, вам интересны эти люди: рекомендации на основе интересов в сфере «{buzzword}»",
		"Рекомендация друзей: пользователи со схожими интересами",
	},
	"communities_recommendation": {
		"Рекомендуем сообщество «{community}»: здесь обсуждают {buzzword}",
		"Популярная группа недели: «{community}»",
	},
	"publications_recommendation": {
		"Рекомендованная запись: «{buzzword}: секреты и лайфхаки»",
		"Вам может понравиться публикация про {buzzword}",
	},
	"return_to_content_reminder": {
		"Вы не досмотрели видео про {buzzword} — продолжить просмотр?",
		"Вы остановились на чтении публикации: вернитесь к чтению в один клик",
	},
	"support_reply": {
		"Служба поддержки ответила на ваше обращение #{id}: «{support_answer}»",
		"По вашей заявке #{id} поступил ответ от специалиста техподдержки",
	},
	"app_notification": {
		"Доступно обновление приложения до версии {number:2,5}.{number:0,9}.{number:0,9} 🚀",
		"Приложение обновлено: добавлены новые возможности и повышена стабильность",
	},
	"service_notification": {
		"Плановые технические работы завершены, все сервисы работают в штатном режиме",
		"Сервисное уведомление: обновлены правила безопасности платформы",
	},
	"game_achievement": {
		"Получено новое достижение «{achievement_title}» в игре «{game_title}» 🏆 (+{number:50,500} очков)",
		"Разблокирован трофей «{achievement_title}» в игре «{game_title}»",
	},
	"promo_selection": {
		"Специально для вас: скидка {number:10,50}% на категорию «{buzzword}»!",
		"Подборка выгодных предложений недели: скидки до {number:30,70}%",
	},
	"new_feature_offer": {
		"Попробуйте новую функцию: персонализированный дайджест новостей уже доступен!",
		"Встречайте обновление: теперь вы можете настраивать {buzzword} в один клик",
	},
}

func GeneratePayloadByType(notificationType string, objectId int) string {
	if templates, exists := exactPayloadTemplates[notificationType]; exists && len(templates) > 0 {
		tpl := pickRandom(templates)
		res, err := gofakeit.Generate(tpl)
		if err == nil && res != "" {
			return strings.ReplaceAll(res, "{id}", fmt.Sprintf("%d", objectId))
		}
	}
	return generateComment()
}

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
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
	"асинхронное программирование", "распределенные базы данных",
}

var russianProducts = []string{
	"Беспроводные наушники Pro", "Умная колонка с голосовым помощником",
	"Механическая клавиатура RGB", "Смарт-часы Series 5",
	"Рюкзак для ноутбука", "Электронная книга 7''",
	"Кофеварка рожковая", "Фитнес-браслет",
}

var russianCommunities = []string{
	"Типичный Программист", "Клуб любителей Go", "Дизайн и интерфейсы",
	"Наука и технологии", "Музыкальный уголок", "Киномания",
	"Книжный клуб", "Geek News",
}

var russianGifts = []string{
	"Золотой кубок", "Букет цветов", "Коробка конфет",
	"Праздничный торт", "Супер-звезда", "Плюшевый мишка",
}

var russianGiftWishes = []string{
	"Отличного настроения!", "С праздником!", "Удачи во всем!",
	"Ты лучший!", "Спасибо за помощь!",
}

var russianReactions = []string{
	"🔥", "❤️", "👍", "👏", "🎉", "😮", "🚀", "😍",
}

var russianTracks = []string{
	"Midnight City Beat", "Echoes of Silence", "Summer Breeze",
	"Cyberpunk Drive", "Acoustic Sunset", "Neon Lights",
}

var russianPlaylists = []string{
	"Музыка для концентрации", "Вечерний чилл", "Топ чарт 2026",
	"Энергичный рок", "Кодинг под лоу-фай", "Акустический вечер",
}

var russianPodcasts = []string{
	"Архитектура IT-систем", "Будни разработчика",
	"Тренды искусственного интеллекта", "Истории стартапов", "Карьера в Tech",
}

var russianAlbums = []string{
	"Летнее путешествие", "Конференция 2026", "Рабочие моменты",
	"Выходные с друзьями", "Пейзажи и город",
}

var russianEvents = []string{
	"Митап разработчиков Go", "Онлайн-конференция TechConf",
	"Вебинар по микросервисам", "Хакатон по AI", "Круглый стол по безопасности",
}

var russianOrderStatuses = []string{
	"В обработке", "Передан в доставку", "В пути",
	"Прибыл в пункт выдачи", "Получен", "Ожидает оплаты",
}

var russianSubscriptionPlans = []string{
	"Премиум Плюс", "VK Музыка", "Combo Max", "Pro подписка", "Ультра Доступ",
}

var russianCommunityRoles = []string{
	"Администратор", "Модератор", "Редактор новостей", "Куратор тем",
}

var russianStreets = []string{
	"Ленина, д. 15", "Пушкина, д. 42", "Мира, д. 7",
	"Советская, д. 23", "Гагарина, д. 10",
}

var russianTransferComments = []string{
	"За обед", "Подарок на праздник", "Возврат долга",
	"За билеты в кино", "На кофе",
}

var russianAchievements = []string{
	"Первооткрыватель", "Знаток кода", "Душа компании",
	"Неудержимый", "Мастер обсуждений", "Активист недели",
}

var russianGames = []string{
	"Cyber Arena", "Magic Quest", "Space Explorer",
	"Pixel Runner", "Battle Tactics",
}

var russianSupportReplies = []string{
	"Вопрос решен, средства зачислены на баланс.",
	"Мы исправили проблему, проверьте работу функции.",
	"Спасибо за обращение! Ошибка передана разработчикам.",
}

var russianReviews = []string{
	"Отличное качество, полностью соответствует описанию!",
	"Быстрая доставка, покупкой очень доволен.",
	"Товар хороший, но упаковка была немного помята.",
	"Превзошло все ожидания, рекомендую к покупке!",
}

var russianQuestions = []string{
	"Какой стек технологий выбрать в 2026 году?",
	"Где лучше учиться Go?",
	"Как вы организуете рабочий день?",
	"Какую книгу по архитектуре посоветуете?",
}

var russianAnswers = []string{
	"Рекомендую начать с официальной документации и пет-проектов.",
	"Лучше всего практиковаться на реальных задачах.",
	"Мне больше всего помог опыт участия в open-source проектах.",
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
