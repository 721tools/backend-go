package model

type NFTCollection struct {
	FloorPrice  string `xorm:"floor_price"`
	Verified    uint64 `xorm:"verified"`
	TotalSupply uint64 `xorm:"total_supply"`
}

func (n NFTCollection) TableName() string {
	return "opensea_collections"
}
