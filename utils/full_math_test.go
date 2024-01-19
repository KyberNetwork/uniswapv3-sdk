package utils

import (
	"fmt"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMulDiv(t *testing.T) {
	// https://github.com/Uniswap/v3-core/blob/main/test/FullMath.spec.ts

	tests := []struct {
		a         string
		b         string
		deno      string
		expResult string
	}{
		{MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Dec()},
		{"0x100000000000000000000000000000000", "0x80000000000000000000000000000000", "0x180000000000000000000000000000000", "113427455640312821154458202477256070485"},
		{"0x100000000000000000000000000000000", "0x2300000000000000000000000000000000", "0x800000000000000000000000000000000", "1488735355279105777652263907513985925120"},
		{"0x100000000000000000000000000000000", "0x3e800000000000000000000000000000000", "0xbb800000000000000000000000000000000", "113427455640312821154458202477256070485"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			r, err := MulDiv(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, r.Dec())
		})
	}

	failTests := []struct {
		a    string
		b    string
		deno string
	}{
		// {"0x100000000000000000000000000000000", "0x5", "0x0"}, // we don't catch div by zero here
		// {"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x0"},
		{"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x1"},
		{MaxUint256.Hex(), MaxUint256.Hex(), new(Uint256).SubUint64(MaxUint256, 1).Hex()},
	}
	for i, tt := range failTests {
		t.Run(fmt.Sprintf("fail test %d", i), func(t *testing.T) {
			_, err := MulDiv(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.NotNil(t, err)
		})
	}
}

func TestMulDivRoundingUp(t *testing.T) {
	// https://github.com/Uniswap/v3-core/blob/main/test/FullMath.spec.ts

	tests := []struct {
		a         string
		b         string
		deno      string
		expResult string
	}{
		{MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Dec()},
		{"0x100000000000000000000000000000000", "0x80000000000000000000000000000000", "0x180000000000000000000000000000000", "113427455640312821154458202477256070486"},
		{"0x100000000000000000000000000000000", "0x2300000000000000000000000000000000", "0x800000000000000000000000000000000", "1488735355279105777652263907513985925120"},
		{"0x100000000000000000000000000000000", "0x3e800000000000000000000000000000000", "0xbb800000000000000000000000000000000", "113427455640312821154458202477256070486"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			r, err := MulDivRoundingUp(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, r.Dec())
		})
	}

	failTests := []struct {
		a    string
		b    string
		deno string
	}{
		// {"0x100000000000000000000000000000000", "0x5", "0x0"}, // we don't catch div by zero here
		// {"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x0"},
		{"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x1"},
		{MaxUint256.Hex(), MaxUint256.Hex(), new(Uint256).SubUint64(MaxUint256, 1).Hex()},
		{"0x1e695d2db4f97", "0x10d5effea103c44aaf18a26b449186a7de3dd6c1ce3d26d03dfd9", "0x2"}, // mulDiv overflows 256 bits after rounding up
		{"0xffffffffffffffffffffffffffffffffffffffb07f6d608e4dcc38020b140b35", "0xffffffffffffffffffffffffffffffffffffffb07f6d608e4dcc38020b140b36", "0xffffffffffffffffffffffffffffffffffffff60fedac11c9b9870041628166c"}, // mulDiv overflows 256 bits after rounding up case 2
	}
	for i, tt := range failTests {
		t.Run(fmt.Sprintf("fail test %d", i), func(t *testing.T) {
			x, err := MulDivRoundingUp(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.NotNil(t, err, x)
		})
	}
}
