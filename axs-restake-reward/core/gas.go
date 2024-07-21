package core

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

const (
	EstimateGasSelector      = "0x3d8527ba"
	EstimateGasMaxRetryCount = 5
	EstimateGasRetryDelay    = 1 * time.Minute
)

func EstimateGasWithRetry(ctx context.Context, ethCli *ethclient.Client, msg ethereum.CallMsg) (gas uint64, err error) {

	for i := 0; i < EstimateGasMaxRetryCount; i++ {
		gas, err = ethCli.EstimateGas(ctx, msg)
		if err == nil {
			return gas, nil // Success
		}

		log.Printf("EstimateGas failed (attempt %d/%d): %v", i+1, EstimateGasMaxRetryCount, err)

		// Delay before retrying
		time.Sleep(EstimateGasRetryDelay)
	}

	return 0, err
}
