package util

import (
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/logging"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"os"
	"strings"
)

var (
	logger = logging.GetLogger()
)

func ParseAbi(path string) abi.ABI {
	jsonFile, err := os.Open(path)
	if err != nil {
		logger.Errorf("Failed to open ABI file: %v", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			logger.Errorf("err: %s", err)
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		logger.Errorf("Failed to read ABI file: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(byteValue)))
	if err != nil {
		logger.Errorf("Failed to parse ABI: %v", err)
	}

	return parsedABI
}
