package core

import (
	"context"
	"errors"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

var (
	logger = logging.GetLogger()
)

const (
	StakingContractAddress        = "0x05b0bb3c1c320b280501b86706c3551995bc8571"
	AutoCompoundRewardsMethodName = "restakeRewards"
)

func AutoCompoundRewards() string {
	configInfo := util.GetConfigInfo()

	chainId := configInfo.ChainID
	gasLimit := configInfo.GasLimit
	accountAddressStr := configInfo.AccountAddress
	PK := configInfo.PK

	accountAddress := common.HexToAddress(accountAddressStr)

	ethCli, ctx := ethClient.GetEthClient()

	contractAddress := common.HexToAddress(StakingContractAddress)

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	contract := bind.NewBoundContract(contractAddress, parsedABI, ethCli, ethCli, ethCli)

	nonce, err := GetPendingNonceWithRetry(ctx, ethCli, accountAddress)
	if err != nil {
		logger.Errorf("err: %s", err)
	}
	logger.Info("Nonce: %d\n", nonce)

	logger.Info("Restaking...")

	gasPrice, err := ethCli.SuggestGasPrice(ctx)
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	msg := ethereum.CallMsg{
		From: accountAddress,
		To:   &contractAddress,
		Data: common.FromHex(EstimateGasSelector),
	}
	gas, err := EstimateGasWithRetry(ctx, ethCli, msg)

	if err != nil {
		logger.Errorf("err: %s", err)
		return ""
	}
	logger.Info("gas: %d\n", gas)

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
	}, AutoCompoundRewardsMethodName)
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	txHash := tx.Hash()
	logger.Infof("Hash: %s - Explorer: https://explorer.roninchain.com/tx/%s", txHash.Hex(), txHash.Hex())

	finalReceipt := waitForTransactionReceipt(ctx, ethCli, txHash, 5*time.Minute)
	if finalReceipt != nil {
		logger.Info("Sleep for 60 seconds")
		time.Sleep(60 * time.Second)
		finalReceipt = getTransactionReceipt(ctx, ethCli, txHash)
		if finalReceipt != nil {
			txStatus := finalReceipt.Status
			if txStatus == 1 {
				logger.Info("Auto Compounded tx status Ok")
			} else {
				logger.Errorf("Auto Compounded tx status Not Ok")
			}
		} else {
			logger.Errorf("Auto Compounded receipt is None")
		}
		logger.Info(finalReceipt)
	}

	return txHash.Hex()
}

func getTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) *types.Receipt {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			logger.Errorf("Transaction %s not found.\n", txHash.Hex())
			return nil
		}
		logger.Errorf("err: %s", err)
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
			logger.Errorf("Time exhausted while waiting for Transaction %s.\n", txHash.Hex())
			return nil
		}
	}
}
