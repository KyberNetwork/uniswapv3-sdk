package utils

import (
	"fmt"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMostSignificantBit(t *testing.T) {
	tests := []struct {
		value     string
		expResult uint
	}{
		{"0x1", 0},
		{"0x100000000000000000000000000000000", 128},
		{"0x10000000000000000", 64},
		{"0x100000000", 32},
		{"0x10000", 16},
		{"0x100", 8},
		{"0x10", 4},
		{"0x4", 2},
		{"0x2", 1},
		{"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 255}, // 2^256 - 1
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			r, err := MostSignificantBit(uint256.MustFromHex(tt.value))
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, r)
		})
	}
}
