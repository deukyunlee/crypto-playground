package core

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/ethClient"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const axsContractAddress = "0x97a9107c1793bc407d6f527b77e7fff4d812bece"

func GetBalance() (*big.Float, error) {
	configInfo := util.GetConfigInfo()

	accountAddressStr := configInfo.AccountAddress
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
	weiPerEther := new(big.Float).SetFloat64(1e18)

	balanceAmountInEther := new(big.Float).Quo(new(big.Float).SetInt(balanceAmount), weiPerEther)

	return balanceAmountInEther, nil
}
