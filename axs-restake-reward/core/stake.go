package core

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func GetStakingAmount() (*big.Float, error) {
	v := util.GetViper()

	accountAddressStr := v.GetString("accountAddress")
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
		fmt.Println("Error unpacking output:", err)
	}
	weiPerEther := new(big.Float).SetFloat64(1e18)

	stakingAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(stakingAmount), weiPerEther)

	return stakingAmountInEther, nil
}

func GetTotalStaked() (*big.Float, error) {
	v := util.GetViper()

	accountAddressStr := v.GetString("accountAddress")
	accountAddress := common.HexToAddress(accountAddressStr)
	contractAddress := common.HexToAddress(StakingContractAddress)

	ethCli, ctx := ethClient.GetEthClient()

	parsedABI := util.ParseAbi("abi/axs_staking_abi.json")

	data, err := parsedABI.Pack("getStakingTotal", accountAddress)
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
		fmt.Println("Error unpacking output:", err)
	}
	weiPerEther := new(big.Float).SetFloat64(1e18)

	stakingAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(stakingAmount), weiPerEther)

	return stakingAmountInEther, nil
}
