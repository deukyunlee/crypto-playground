package core

import (
	"context"
	"errors"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

const (
	StakingContractAddress   = "0x05b0bb3c1c320b280501b86706c3551995bc8571"
	RestakeRewardsMethodName = "restakeRewards"
)

func RestakeRewards() {
	v := util.GetViper()

	chainId := v.GetInt64("chainId")
	gasLimit := v.GetUint64("gasLimit")
	accountAddressStr := v.GetString("accountAddress")
	PK := v.GetString("pk")

	accountAddress := common.HexToAddress(accountAddressStr)

	ethCli, ctx := ethClient.GetEthClient()

	contractAddress := common.HexToAddress(StakingContractAddress)

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	contract := bind.NewBoundContract(contractAddress, parsedABI, ethCli, ethCli, ethCli)

	nonce, err := GetPendingNonceWithRetry(ctx, ethCli, accountAddress)
	if err != nil {
		log.Printf("err: %s", err)
	}
	log.Printf("Nonce: %d\n", nonce)

	log.Println("Restaking...")

	gasPrice, err := ethCli.SuggestGasPrice(ctx)
	if err != nil {
		log.Printf("err: %s", err)
	}

	msg := ethereum.CallMsg{
		From: accountAddress,
		To:   &contractAddress,
		Data: common.FromHex(EstimateGasSelector),
	}
	gas, err := EstimateGasWithRetry(ctx, ethCli, msg)

	if err != nil {
		log.Printf("err: %s", err)
		return
	}
	log.Printf("gas: %d\n", gas)

	// Create the transaction
	tx, err := contract.Transact(&bind.TransactOpts{
		From:     accountAddress,
		Nonce:    big.NewInt(int64(nonce)),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			privateKey, err := crypto.HexToECDSA(PK[2:])
			if err != nil {
				return nil, err
			}
			signer := types.NewEIP155Signer(big.NewInt(chainId))
			return types.SignTx(tx, signer, privateKey)
		},
	}, RestakeRewardsMethodName)
	if err != nil {
		log.Printf("err: %s", err)
	}

	txHash := tx.Hash()
	log.Printf("Hash: %s - Explorer: https://explorer.roninchain.com/tx/%s", txHash.Hex(), txHash.Hex())

	finalReceipt := waitForTransactionReceipt(ctx, ethCli, txHash, 5*time.Minute)
	if finalReceipt != nil {
		log.Println("Sleep for 60 seconds")
		time.Sleep(60 * time.Second)
		finalReceipt = getTransactionReceipt(ctx, ethCli, txHash)
		if finalReceipt != nil {
			txStatus := finalReceipt.Status
			if txStatus == 1 {
				log.Println("Restake tx status Ok")
				log.Println("DONE")
			} else {
				log.Println("Restake tx status Not Ok")
				log.Println("DONE & FAILED")
			}
		} else {
			log.Println("Restake receipt is None")
		}
		log.Println(finalReceipt)
	}
}

func getTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) *types.Receipt {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			log.Printf("Transaction %s not found.\n", txHash.Hex())
			return nil
		}
		log.Printf("err: %s", err)
	}
	return receipt
}

func waitForTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) *types.Receipt {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ticker.C:
			receipt := getTransactionReceipt(ctx, client, txHash)
			if receipt != nil {
				return receipt
			}
		case <-timeoutTimer.C:
			log.Printf("Time exhausted while waiting for Transaction %s.\n", txHash.Hex())
			return nil
		}
	}
}
