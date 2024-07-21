package core

import (
	"context"
	"fmt"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const axsContractAddress = "0x97a9107c1793bc407d6f527b77e7fff4d812bece"

func GetBalance(ctx context.Context) (*big.Int, error) {
	v := util.GetViper()

	accountAddressStr := v.GetString("accountAddress")
	accountAddress := common.HexToAddress(accountAddressStr)
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
		fmt.Println("Error unpacking output:", err)
	}
	balanceAmount.Div(balanceAmount, big.NewInt(1000000000000000000))

	return balanceAmount, nil
}

func GetStakingAmount(ctx context.Context) (*big.Int, error) {
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

	stakingAmount.Div(stakingAmount, big.NewInt(1000000000000000000))
	return stakingAmount, nil
}
