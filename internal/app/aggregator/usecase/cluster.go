package usecase

import (
	"context"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
)

type ClusterUseCase interface {
	GetAllClusters(ctx context.Context) ([]string, error)
}

type Cluster struct {
	repository repository.ClusterRepository
}

func NewCluster(repository repository.ClusterRepository) *Cluster {
	return &Cluster{repository: repository}
}

func (obj *Cluster) GetAllClusters(ctx context.Context) ([]string, error) {
	return obj.repository.GetAllClusters(ctx)
}

func GetClusterForNotificationType(notificationType string) string {
	var NotificationTypeToCluster = map[string]string{
		// friends_and_subscriptions
		"friend_request":           "friends_and_subscriptions",
		"friend_request_accepted":  "friends_and_subscriptions",
		"new_subscriber":           "friends_and_subscriptions",
		"friend_invitation":        "friends_and_subscriptions",
		"found_acquaintance":       "friends_and_subscriptions",
		"recommended_acquaintance": "friends_and_subscriptions",
		"friend_birthday":          "friends_and_subscriptions",
		"gift":                     "friends_and_subscriptions",

		// content_reactions
		"like":                "content_reactions",
		"reaction_on_post":    "content_reactions",
		"reaction_on_photo":   "content_reactions",
		"reaction_on_video":   "content_reactions",
		"reaction_on_story":   "content_reactions",
		"reaction_on_comment": "content_reactions",
		"reaction_on_product": "content_reactions",
		"repost":              "content_reactions",

		// comments_and_replies
		"comment_on_post":    "comments_and_replies",
		"comment_on_photo":   "comments_and_replies",
		"comment_on_video":   "comments_and_replies",
		"comment_on_article": "comments_and_replies",
		"comment_on_product": "comments_and_replies",
		"reply_to_comment":   "comments_and_replies",

		// mentions_and_tags
		"author_mention":     "mentions_and_tags",
		"user_mention":       "mentions_and_tags",
		"photo_tag":          "mentions_and_tags",
		"post_mention":       "mentions_and_tags",
		"comment_mention":    "mentions_and_tags",
		"story_mention":      "mentions_and_tags",
		"discussion_mention": "mentions_and_tags",
		"story_question":     "mentions_and_tags",
		"story_answer":       "mentions_and_tags",

		// new_publications_and_media
		"new_post":           "new_publications_and_media",
		"new_photo":          "new_publications_and_media",
		"new_story":          "new_publications_and_media",
		"new_video":          "new_publications_and_media",
		"new_clip":           "new_publications_and_media",
		"new_audio":          "new_publications_and_media",
		"new_album":          "new_publications_and_media",
		"new_podcast":        "new_publications_and_media",
		"new_playlist":       "new_publications_and_media",
		"stream_started":     "new_publications_and_media",
		"new_author_content": "new_publications_and_media",

		// communities_and_management
		"new_community_content":    "communities_and_management",
		"community_invitation":     "communities_and_management",
		"join_request_accepted":    "communities_and_management",
		"role_changed":             "communities_and_management",
		"co_owner_status_changed":  "communities_and_management",
		"new_community_subscriber": "communities_and_management",
		"suggested_post":           "communities_and_management",
		"community_comment":        "communities_and_management",
		"community_reaction":       "communities_and_management",

		// shop_and_orders
		"product_comment":         "shop_and_orders",
		"product_review":          "shop_and_orders",
		"review_published":        "shop_and_orders",
		"review_rejected":         "shop_and_orders",
		"order_status_changed":    "shop_and_orders",
		"abandoned_cart_reminder": "shop_and_orders",
		"delivery_info":           "shop_and_orders",

		// payments_and_subscriptions
		"money_request":          "payments_and_subscriptions",
		"money_transfer":         "payments_and_subscriptions",
		"balance_changed":        "payments_and_subscriptions",
		"payment_state_changed":  "payments_and_subscriptions",
		"subscription_purchased": "payments_and_subscriptions",
		"subscription_renewed":   "payments_and_subscriptions",
		"payment_error":          "payments_and_subscriptions",

		// security_and_account_state
		"account_changed":   "security_and_account_state",
		"data_verification": "security_and_account_state",
		"phone_linked":      "security_and_account_state",
		"phone_unlinked":    "security_and_account_state",
		"email_linked":      "security_and_account_state",
		"email_unlinked":    "security_and_account_state",
		"security_warning":  "security_and_account_state",
		"access_restricted": "security_and_account_state",
		"access_restored":   "security_and_account_state",

		// events_reminders_and_recommendations
		"event_starting_soon":         "events_reminders_and_recommendations",
		"call_starting_soon":          "events_reminders_and_recommendations",
		"birthday_reminder":           "events_reminders_and_recommendations",
		"saved_wish_reminder":         "events_reminders_and_recommendations",
		"people_recommendation":       "events_reminders_and_recommendations",
		"communities_recommendation":  "events_reminders_and_recommendations",
		"publications_recommendation": "events_reminders_and_recommendations",
		"return_to_content_reminder":  "events_reminders_and_recommendations",

		// service_and_promo_messages
		"support_reply":        "service_and_promo_messages",
		"app_notification":     "service_and_promo_messages",
		"service_notification": "service_and_promo_messages",
		"game_achievement":     "service_and_promo_messages",
		"promo_selection":      "service_and_promo_messages",
		"new_feature_offer":    "service_and_promo_messages",
	}
	if cluster, ok := NotificationTypeToCluster[notificationType]; ok {
		return cluster
	}
	return "service_and_promo_messages"
}
