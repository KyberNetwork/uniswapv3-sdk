package utils

import (
	"math/big"
	"testing"

	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestGetSqrtRatioAtTick(t *testing.T) {
	_, err := GetSqrtRatioAtTick(MinTick - 1)
	assert.ErrorIs(t, err, ErrInvalidTick, "tick tool small")

	_, err = GetSqrtRatioAtTick(MaxTick + 1)
	assert.ErrorIs(t, err, ErrInvalidTick, "tick tool large")

	rmax, _ := GetSqrtRatioAtTick(MinTick)
	assert.Equal(t, rmax, MinSqrtRatio, "returns the correct value for min tick")

	var r Uint160
	_ = GetSqrtRatioAtTickV2(MinTick+1, &r)
	assert.Equal(t, uint256.NewInt(4295343490), &r, "returns the correct value for min tick + 1")

	r0, _ := GetSqrtRatioAtTick(0)
	assert.Equal(t, r0, new(big.Int).Lsh(constants.One, 96), "returns the correct value for tick 0")

	rmin, _ := GetSqrtRatioAtTick(MaxTick)
	assert.Equal(t, rmin, MaxSqrtRatio, "returns the correct value for max tick")

	_ = GetSqrtRatioAtTickV2(MaxTick-1, &r)
	assert.Equal(t, uint256.MustFromDecimal("1461373636630004318706518188784493106690254656249"), &r, "returns the correct value for max tick - 1")

	_ = GetSqrtRatioAtTickV2(MaxTick, &r)
	assert.Equal(t, uint256.MustFromDecimal("1461446703485210103287273052203988822378723970342"), &r, "returns the correct value for max tick")
}

func TestGetTickAtSqrtRatio(t *testing.T) {
	tmin, _ := GetTickAtSqrtRatio(MinSqrtRatio)
	assert.Equal(t, tmin, MinTick, "returns the correct value for sqrt ratio at min tick")

	_, err := GetTickAtSqrtRatioV2(new(uint256.Int).SubUint64(MinSqrtRatioU256, 1))
	assert.ErrorIs(t, ErrInvalidSqrtRatio, err)

	_, err = GetTickAtSqrtRatioV2(MaxSqrtRatioU256)
	assert.ErrorIs(t, ErrInvalidSqrtRatio, err)

	tmax, _ := GetTickAtSqrtRatio(new(big.Int).Sub(MaxSqrtRatio, constants.One))
	assert.Equal(t, tmax, MaxTick-1, "returns the correct value for sqrt ratio at max tick")

	tt, _ := GetTickAtSqrtRatio(big.NewInt(4295343490))
	assert.Equal(t, MinTick+1, tt)
}
