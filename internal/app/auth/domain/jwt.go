package domain

import "time"

type RefreshToken struct {
	Uuid      string
	UserId    int
	ExpiredAt time.Time
}
