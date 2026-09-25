package domain

import "time"

type Notification struct {
	Id         int
	Payload    string
	ReceivedAt time.Time
	UserId     int
	ClusterId  int
}

type NotificationType struct {
	Id   int
	Type string
}
