package domain

import "time"

type Digest struct {
	Id        int
	UserId    int
	BlockId   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
