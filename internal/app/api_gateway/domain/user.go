package domain

import "time"

type User struct {
	Id        int
	IsStaff   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Active    bool
}

type Settings struct {
	Id        int
	UserId    int
	Interval  time.Duration
	Important []string
	AvatarUrl string
	FirstName string
	LastName  string
	UpdatedAt time.Time
}
