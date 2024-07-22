package util

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"log"
	"os"
	"strings"
)

func ParseAbi(path string) abi.ABI {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open ABI file: %v", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Printf("err: %s", err)
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(byteValue)))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
	}

	return parsedABI
}
