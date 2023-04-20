package boot

import (
	"context"
	"errors"

	"github.com/721tools/backend-go/indexer/implement/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/pulling"
	"github.com/721tools/backend-go/indexer/pkg/types/v1beta1"
	"github.com/721tools/backend-go/indexer/pkg/utils/files"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

func Refill(fileIO files.FileIO) (err error) {
	refill := new(v1beta1.Refill)
	if err = fileIO.Read(refill); err != nil {
		return
	}
	if refill.Status.RunHeight >= refill.Spec.BlockEnd {
		return errors.New("please check the yaml file, task is finished! ")
	}

	if refill.Spec.BlockStart <= 0 {
		return errors.New("please check the yaml file, block start must greater than zero! ")
	}
	if refill.Spec.BlockStart >= refill.Spec.BlockEnd {
		return errors.New("please check the yaml file, block start are greater than block end! ")
	}
	if refill.Status.RunHeight > refill.Spec.BlockEnd {
		return
	}
	if refill.Status.RunHeight < refill.Spec.BlockStart {
		refill.Status.RunHeight = refill.Spec.BlockStart
		if err = fileIO.Write(refill); err != nil {
			return
		}
	}
	contractAddress := hex.HexstrToHex(refill.Spec.ContractAddress)
	svc := service.NewBlockService()
	puller := pulling.NewPuller(3, 30, client.GetClient())
	for height := refill.Status.RunHeight; height <= refill.Spec.BlockEnd; height++ {
		if err = findOrRefillBlock(puller, svc, height, contractAddress); err != nil {
			return
		}
		refill.Status.RunHeight = height + 1
		if err = fileIO.Write(refill); err != nil {
			return
		}
	}
	return
}

func findOrRefillBlock(puller *pulling.Puller, svc service.BlockIface, height uint64, contract hex.Hex) (err error) {
	exists := svc.FindContractInBlock(context.Background(), contract, height)
	if !exists {
		return nil
	}
	if err = svc.ClearBlocks(context.Background(), []uint64{height}); err != nil {
		return
	}
	return puller.Fix(height)
}
