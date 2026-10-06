package domain

import "time"

type Block struct {
	Id     int
	UserId int
	Active bool
	Start  time.Time
	End    time.Time
}

type BlockPart struct {
	Id          int
	BlockId     int
	ClusterType string
	Text        string
	Start       time.Time
	End         time.Time
}
