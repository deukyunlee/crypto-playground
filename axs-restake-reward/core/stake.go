package core

import (
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UserRewardResult struct {
	DebitedRewards   *big.Int
	CreditedRewards  *big.Int
	LastClaimedBlock *big.Int
}

const StakingManagerContract = "0x8bd81a19420bad681b7bfc20e703ebd8e253782d"

func GetStakingAmount() (*big.Float, error) {
	configInfo := util.GetConfigInfo()

	accountAddressStr := configInfo.AccountAddress
	accountAddress := common.HexToAddress(accountAddressStr)
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

func GetTotalStaked() (*big.Float, error) {

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

func GetUserRewardInfo() (UserRewardResult, error) {
	configInfo := util.GetConfigInfo()

	var userReward UserRewardResult
	accountAddressStr := configInfo.AccountAddress
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
