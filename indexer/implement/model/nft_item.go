package model

type NFTItem struct {
	ImageUrl    string `xorm:"image_url"`
	TraitsScore uint64 `xorm:"traits_score"`
	TraitsRank  uint64 `xorm:"traits_rank"`
}

func (n NFTItem) TableName() string {
	return "opensea_items"
}
