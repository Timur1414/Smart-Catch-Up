package domain

import "time"

type Digest struct {
	Id        int
	UserId    int
	BlockId   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DigestCategoryCache struct {
	ClusterType string    `json:"cluster_type"`
	Text        string    `json:"text"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Rank        int       `json:"rank"`
}

type DigestCache struct {
	Start      time.Time             `json:"start"`
	End        time.Time             `json:"end"`
	Categories []DigestCategoryCache `json:"categories"`
}

func NewDigestCache(block Block) *DigestCache {
	categories := make([]DigestCategoryCache, len(block.Parts))
	for i, part := range block.Parts {
		categories[i] = DigestCategoryCache{
			ClusterType: part.ClusterType,
			Text:        part.Text,
			Start:       part.Start,
			End:         part.End,
			Rank:        part.Rank,
		}
	}
	return &DigestCache{
		Start:      block.Start,
		End:        block.End,
		Categories: categories,
	}
}
