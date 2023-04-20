package service

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/index/pkg/db"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

// https://etherscan.io/tx/0x08484c11c73f8bdab01f2f172eb5dece1f7601be9da3c96c072312fe275e6e39#eventlog
func TestGetAcountBalance(t *testing.T) {
	dsn := "root:@tcp(localhost:3306)/test_dev?parseTime=true"
	db.NewDBEngine(dsn)

	service := NewAcountService()

	balance := []TokenBalance{}
	addrs := hex.HexstrToHex("0x8D8890235639AA0715AFB141FD2943F19E101E78")
	service.GetAcountBalance(context.TODO(), addrs, &balance)

	t.Log("debug balance", "balance is", balance)
}

func TestGetAcountTx(t *testing.T) {
	dsn := "root:@tcp(localhost:3306)/test_dev?parseTime=true"
	db.NewDBEngine(dsn)

	service := NewAcountService()

	tx := []Tx{}
	addrs := hex.HexstrToHex("0x8D8890235639AA0715AFB141FD2943F19E101E78")
	service.GetAcountTx(context.TODO(), addrs, 1, 2, &tx)

	for i, _ := range tx {
		t.Log("debug tx", " height", tx[i].Height, " hash", tx[i].TxHash, " from ", tx[i].From, " to ", tx[i].To, " token", tx[i].Symbol)
	}
	t.Log("debug tx", " height", tx)
}

func TestGetAcountTxCnt(t *testing.T) {
	dsn := "root:@tcp(localhost:3306)/test_dev?parseTime=true"
	db.NewDBEngine(dsn)

	service := NewAcountService()

	addrs := hex.HexstrToHex("0x8D8890235639AA0715AFB141FD2943F19E101E78")

	cnt, err := service.GetAcountTxCnt(context.TODO(), addrs)

	t.Log("debug tx", " total", cnt, " err", err)
}
