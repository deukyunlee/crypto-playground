package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

const (
	PendingNonceMaxRetryCount = 5
	PendingNonceRetryDelay    = 1 * time.Minute
)

func GetPendingNonceWithRetry(ctx context.Context, ethCli *ethclient.Client, accountAddress common.Address) (nonce uint64, err error) {

	for i := 0; i < PendingNonceMaxRetryCount; i++ {
		nonce, err = ethCli.PendingNonceAt(ctx, accountAddress)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("EstimateGas failed (attempt %d/%d): %v", i+1, PendingNonceMaxRetryCount, err)

		// Delay before retrying
		time.Sleep(PendingNonceRetryDelay)
	}

	return 0, err
}
