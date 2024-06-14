package periphery

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDecodeFirstPool(t *testing.T) {
	input := "0x9e32b13ce7f2e80a01932b42553652e053d6ed8e000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	input = strings.TrimPrefix(input, "0x")

	// Convert the binary data to bytes
	pathBytes, err := hex.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}

	token0, token1, fee, err := DecodeFirstPool(pathBytes)
	if err != nil {
		t.Fatal(err)
	}

	expected := struct {
		Token0 string
		Token1 string
		Fee    string
	}{
		Token0: strings.ToLower("0x9E32b13ce7f2E80A01932B42553652E053D6ed8e"),
		Token1: strings.ToLower("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		Fee:    "3000",
	}

	assert.Equal(t, false, HasMultiplePools(pathBytes))
	assert.Equal(t, expected.Token0, strings.ToLower(token0.String()))
	assert.Equal(t, expected.Token1, strings.ToLower(token1.String()))
	assert.Equal(t, expected.Fee, fee.String())
}
