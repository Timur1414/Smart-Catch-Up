package domain

import "time"

type User struct {
	Id     int
	Email  string
	Active bool
}

type Settings struct {
	Id        int
	UserId    int
	AvatarUrl string
	UpdatedAt time.Time
}
