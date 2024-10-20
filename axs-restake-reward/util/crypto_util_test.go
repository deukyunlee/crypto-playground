package util

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetAddressFromPrivateKey(t *testing.T) {
	privateKeyHex := "4c0883a69102937d6231471b5ecb4c6f44c0c1d6bcf2404b5e06de63480d710d"

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Cannot convert public key to ECDSA")
	}
	expectedAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	address := GetAddressFromPrivateKey(privateKeyHex)

	assert.Equal(t, expectedAddress, address)
}
