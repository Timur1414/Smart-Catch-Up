package domain

import "time"

type User struct {
	Id        int
	IsStaff   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	VkId      int
}
