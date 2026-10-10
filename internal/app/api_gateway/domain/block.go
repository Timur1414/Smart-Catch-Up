package domain

import "time"

type Block struct {
	Id     int
	UserId int
	Status string
	Start  time.Time
	End    time.Time
	Parts  []BlockPart
}

type BlockPart struct {
	Id          int
	BlockId     int
	ClusterType string
	Text        string
	Start       time.Time
	End         time.Time
}
