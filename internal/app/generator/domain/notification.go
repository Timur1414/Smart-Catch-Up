package domain

import "time"

type NotificationAction struct {
	ActionType   string
	ActionTarget string
}

type Notification struct {
	Id               int
	NotificationType string
	ActorId          int
	ActorName        string
	ObjectId         int
	ObjectType       string
	CreatedAt        time.Time
	ReadAt           time.Time
	Payload          string
	Actions          []NotificationAction
}
