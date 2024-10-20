package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestParseAbi(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_abi_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	abiJSON := `[
		{
			"constant": true,
			"inputs": [],
			"name": "myFunction",
			"outputs": [{"name": "", "type": "uint256"}],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		}
	]`

	_, err = tempFile.Write([]byte(abiJSON))
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())

	fmt.Println(tempFile.Name())
	parsedABI := ParseAbi(tempFile.Name())

	assert.NotNil(t, parsedABI)

	_, ok := parsedABI.Methods["myFunction"]
	assert.True(t, ok, "Expected method 'myFunction' to be parsed")
}
