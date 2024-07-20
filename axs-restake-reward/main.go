package main

import (
	"context"
	"errors"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Configurations from config package
const (
	STAKE_LOGS                  = "./logs/staking_logs.log"
	CONTRACT_ADDRESS            = "0x05b0bb3c1c320b280501b86706c3551995bc8571"
	RESTAKE_REWARDS_METHOD_NAME = "restakeRewards"
	ESTIMATE_GAS_SELECTOR       = "0x3d8527ba"
)

func getTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) *types.Receipt {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			log.Printf("Transaction %s not found.\n", txHash.Hex())
			return nil
		}
		log.Fatal(err)
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

func main() {
	// Setup logging
	logFile, err := os.OpenFile(STAKE_LOGS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	v := util.GetViper()

	chainId := v.GetInt64("chainId")
	gasLimit := v.GetUint64("gasLimit")
	accountAddressStr := v.GetString("accountAddress")

	accountAddress := common.HexToAddress(accountAddressStr)
	PK := v.GetString("pk")

	ethCli, ctx := ethClient.GetEthClient()

	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)

	parsedABI := util.ParseAbi()

	contract := bind.NewBoundContract(contractAddress, parsedABI, ethCli, ethCli, ethCli)

	nonce, err := ethCli.PendingNonceAt(ctx, accountAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Nonce: %d\n", nonce)

	log.Println("Restaking...")

	gasPrice, err := ethCli.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		From: accountAddress,
		To:   &contractAddress,
		Data: common.FromHex(ESTIMATE_GAS_SELECTOR),
	}
	gas, err := ethCli.EstimateGas(ctx, msg)
	if err != nil {
		log.Fatal(err)
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
	}, RESTAKE_REWARDS_METHOD_NAME)
	if err != nil {
		log.Fatal(err)
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
