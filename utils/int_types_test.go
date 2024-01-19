package utils

import (
	"fmt"
	"testing"

	"github.com/KyberNetwork/int256"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToInt256(t *testing.T) {
	// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/692dbc560f48b2a5160e6e4f78302bb93314cd88/test/utils/math/SafeCast.test.js#L124

	successCases := []string{
		"0x0",
		"0x1",
		"0x18fe",
		"0x9234bbe",
		"0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // INT256_MAX
	}

	var res int256.Int
	for _, tc := range successCases {
		t.Run(fmt.Sprintf("test %s", tc), func(t *testing.T) {
			ui := uint256.MustFromHex(tc)
			err := ToInt256(ui, &res)
			require.Nil(t, err)

			// should be equal to the original value
			assert.Equal(t, ui.Dec(), res.Dec())
		})
	}

	// INT256_MAX+1
	assert.ErrorIs(t, ErrExceedMaxInt256, ToInt256(uint256.MustFromHex("0x8000000000000000000000000000000000000000000000000000000000000000"), &res))
	// INT256_MAX+2
	assert.ErrorIs(t, ErrExceedMaxInt256, ToInt256(uint256.MustFromHex("0x8000000000000000000000000000000000000000000000000000000000000001"), &res))
	// UINT256_MAX
	assert.ErrorIs(t, ErrExceedMaxInt256, ToInt256(uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), &res))
}

func TestAddDeltaInPlace(t *testing.T) {
	//https://github.com/Uniswap/v3-core/blob/main/test/LiquidityMath.spec.ts

	successCases := []struct {
		x    *Uint128
		y    *Int128
		expX *Uint128
	}{
		{uint256.NewInt(1), int256.NewInt(0), uint256.NewInt(1)},
		{uint256.NewInt(1), int256.NewInt(-1), uint256.NewInt(0)},
		{uint256.NewInt(1), int256.NewInt(1), uint256.NewInt(2)},
	}

	for i, tc := range successCases {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			err := AddDeltaInPlace(tc.x, tc.y)
			require.Nil(t, err)

			// should be equal to the original value
			assert.Equal(t, tc.expX.Dec(), tc.x.Dec())
		})
	}

	// 2**128-15 + 15 overflows
	tmp := new(uint256.Int).SubUint64(new(uint256.Int).Exp(uint256.NewInt(2), uint256.NewInt(128)), 15)
	assert.ErrorIs(t, ErrOverflowUint128, AddDeltaInPlace(tmp, int256.NewInt(15)))
	// 0 + -1 underflows
	assert.ErrorIs(t, ErrOverflowUint128, AddDeltaInPlace(uint256.NewInt(0), int256.NewInt(-1)))
	// 3 + -4 underflows underflows
	assert.ErrorIs(t, ErrOverflowUint128, AddDeltaInPlace(uint256.NewInt(3), int256.NewInt(-4)))
}
