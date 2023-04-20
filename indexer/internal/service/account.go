package service

import (
	"context"

	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"xorm.io/xorm"
)

type TokenBalance struct {
	Address hex.Hex `xorm:"address"`
	From    hex.Hex `xorm:"from"`
	To      hex.Hex `xorm:"to"`
	Type    uint    `xorm:"type"`
	Value   hex.Hex `xorm:"value"`
	Name    string  `xorm:"name"`
	Symbol  string  `xorm:"symbol"`
	Decimal uint64  `xorm:"decimal"`
}

type AcountService struct {
	xorm *xorm.Engine
}

func NewAcountService() *AcountService {
	return &AcountService{xorm: db.GetDBEngine()}
}

func (b *AcountService) GetAcountBalance(ctx context.Context, address hex.Hex, balance *[]TokenBalance) error {
	b.xorm.SQL("select c.type, c.type, c.name, c.symbol, c.decimal, t.* from contract c right join token_flow t on c.address=t.address where t.`to`=?", address).Find(balance)
	return nil
}

func (b *AcountService) GetAcountTxCnt(ctx context.Context, address hex.Hex) (uint64, error) {
	counts, err := b.xorm.SQL("select count(1) as cnt from token_flow t where t.from = ? or t.to= ?", address, address).Count()
	return uint64(counts), err
}

func (b *AcountService) GetAcountTx(ctx context.Context, address hex.Hex, page, limit uint64, tx *[]Tx) error {
	if page == 0 {
		page = 1
	}
	start := (page - 1) * limit
	step := limit

	b.xorm.SQL("select t.*, c.type, c.name, c.symbol, c.decimal, c.total_supply from token_flow t left join contract c on t.address = c.address where t.from = ? or t.to= ? limit ?,?", address, address, start, step).Find(tx)
	return nil
}
