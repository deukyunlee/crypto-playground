package core

import (
	"context"
	"errors"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/ethClient"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/logging"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

const axsContractAddress = "0x97a9107c1793bc407d6f527b77e7fff4d812bece"

type EvmManager struct {
}

const (
	StakingContractAddress        = "0x05b0bb3c1c320b280501b86706c3551995bc8571"
	AutoCompoundRewardsMethodName = "restakeRewards"
	StakingManagerContract        = "0x8bd81a19420bad681b7bfc20e703ebd8e253782d"
	EstimateGasSelector           = "0x3d8527ba"
)

var (
	logger = logging.GetLogger()
)

func NewClientManager(client *ethclient.Client) *ethClient.ClientManger {
	return &ethClient.ClientManger{Client: client}
}

func (m *EvmManager) GetBalance(accountAddress string) (*big.Float, error) {

	contractAddress := common.HexToAddress(axsContractAddress)

	ethCli, ctx := ethClient.GetEthClient()

	parsedABI := util.ParseAbi("abi/axs_balance_of_abi.json")

	data, err := parsedABI.Pack("balanceOf", accountAddress)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	output, err := ethCli.CallContract(ctx, callMsg, nil)
	if err != nil {
		return nil, err
	}

	var balanceAmount *big.Int
	err = parsedABI.UnpackIntoInterface(&balanceAmount, "balanceOf", output)
	if err != nil {
		logger.Errorf("Error unpacking output: %s", err)
		return nil, err
	}

	weiPerEther := new(big.Float).SetFloat64(1e18)
	balanceAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(balanceAmount), weiPerEther)

	return balanceAmountInEther, nil
}

func (m *EvmManager) GetStakingAmount(accountAddress string) (*big.Float, error) {
	contractAddress := common.HexToAddress(StakingContractAddress)

	ethCli, ctx := ethClient.GetEthClient()

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	data, err := parsedABI.Pack("getStakingAmount", accountAddress)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	output, err := ethCli.CallContract(ctx, callMsg, nil)
	if err != nil {

		return nil, err
	}

	var stakingAmount *big.Int
	err = parsedABI.UnpackIntoInterface(&stakingAmount, "getStakingAmount", output)
	if err != nil {
		logger.Errorf("Error unpacking output: %s", err)
	}
	weiPerEther := new(big.Float).SetFloat64(1e18)

	stakingAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(stakingAmount), weiPerEther)

	return stakingAmountInEther, nil
}

func (m *EvmManager) GetTotalStaked() (*big.Float, error) {
	contractAddress := common.HexToAddress(StakingContractAddress)

	ethCli, ctx := ethClient.GetEthClient()

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	data, err := parsedABI.Pack("getStakingTotal")
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	output, err := ethCli.CallContract(ctx, callMsg, nil)
	if err != nil {

		return nil, err
	}

	var stakingAmount *big.Int
	err = parsedABI.UnpackIntoInterface(&stakingAmount, "getStakingTotal", output)
	if err != nil {
		logger.Errorf("Error unpacking output: %s", err)
	}
	weiPerEther := new(big.Float).SetFloat64(1e18)

	stakingAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(stakingAmount), weiPerEther)

	return stakingAmountInEther, nil
}

func (m *EvmManager) GetUserRewardInfo() (UserRewardResult, error) {

	var userReward UserRewardResult
	pk := util.GetConfigInfo().PK

	accountAddressStr := util.GetAddressFromPrivateKey(pk)

	stakingManagerContractAddress := common.HexToAddress(StakingManagerContract)
	stakingContractAddress := common.HexToAddress(StakingContractAddress)
	accountAddress := common.HexToAddress(accountAddressStr)

	ethCli, ctx := ethClient.GetEthClient()

	parsedABI := util.ParseAbi("abi/staking_manager_abi.json")

	data, err := parsedABI.Pack("userRewardInfo", stakingContractAddress, accountAddress)
	if err != nil {
		return userReward, err
	}
	callMsg := ethereum.CallMsg{
		To:   &stakingManagerContractAddress,
		Data: data,
	}

	output, err := ethCli.CallContract(ctx, callMsg, nil)
	if err != nil {
		return userReward, err
	}

	err = parsedABI.UnpackIntoInterface(&userReward, "userRewardInfo", output)
	if err != nil {
		logger.Errorf("Error unpacking output: %s", err)
	}

	return userReward, nil
}

func (m *EvmManager) AutoCompoundRewards() (transactionHash string, err error) {
	configInfo := util.GetConfigInfo()

	chainId := configInfo.ChainID
	gasLimit := configInfo.GasLimit
	pk := util.GetConfigInfo().PK
	accountAddressStr := util.GetAddressFromPrivateKey(pk)

	accountAddress := common.HexToAddress(accountAddressStr)

	ethCli, ctx := ethClient.GetEthClient()

	contractAddress := common.HexToAddress(StakingContractAddress)

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	contract := bind.NewBoundContract(contractAddress, parsedABI, ethCli, ethCli, ethCli)

	cliManager := NewClientManager(ethCli)
	nonce, err := cliManager.GetPendingNonceWithRetry(accountAddress, ctx)

	if err != nil {
		logger.Errorf("err: %s", err)
		return "", err
	}
	logger.Infof("Nonce: %d\n", nonce)

	logger.Info("Auto Compound Rewards...")

	gasPrice, err := ethCli.SuggestGasPrice(ctx)
	if err != nil {
		logger.Errorf("err: %s", err)
		return "", err
	}

	msg := ethereum.CallMsg{
		From: accountAddress,
		To:   &contractAddress,
		Data: common.FromHex(EstimateGasSelector),
	}
	gas, err := cliManager.EstimateGasWithRetry(ctx, msg)

	if err != nil {
		logger.Errorf("err: %s", err)
		return "", err
	}
	logger.Infof("gas: %d\n", gas)

	// Create the transaction
	tx, err := contract.Transact(&bind.TransactOpts{
		From:     accountAddress,
		Nonce:    big.NewInt(int64(nonce)),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			privateKey, err := crypto.HexToECDSA(pk[2:])
			if err != nil {
				return nil, err
			}
			signer := types.NewEIP155Signer(big.NewInt(chainId))
			return types.SignTx(tx, signer, privateKey)
		},
	}, AutoCompoundRewardsMethodName)
	if err != nil {
		logger.Errorf("err: %s", err)
		return "", err
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

	return txHash.Hex(), err
}

func (m *EvmManager) GetLastClaimedTime() time.Time {
	userRewardInfo, err := m.GetUserRewardInfo()
	if err != nil {
		logger.Errorf("err: %s", err)
		return time.Unix(0, 0)
	}
	lastClaimedTimestampUnix := userRewardInfo.LastClaimedBlock.Int64()
	lastClaimedTime := time.Unix(lastClaimedTimestampUnix, 0).UTC()
	logger.Infof("lastClaimedTime: %s\n", lastClaimedTime.In(util.Location))

	return lastClaimedTime
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
