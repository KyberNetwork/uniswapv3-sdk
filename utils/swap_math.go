package utils

import (
	"github.com/KyberNetwork/int256"
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/holiman/uint256"
)

const MaxFeeInt = 1000000

var MaxFeeUint256 = uint256.NewInt(MaxFeeInt)

func ComputeSwapStep(
	sqrtRatioCurrentX96,
	sqrtRatioTargetX96 *Uint160,
	liquidity *Uint128,
	amountRemaining *int256.Int,
	feePips constants.FeeAmount,
) (sqrtRatioNextX96 *Uint160, amountIn, amountOut, feeAmount *uint256.Int, err error) {
	zeroForOne := sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0
	exactIn := amountRemaining.Sign() >= 0

	var amountRemainingU uint256.Int
	if exactIn {
		ToUInt256(amountRemaining, &amountRemainingU)
	} else {
		ToUInt256(amountRemaining, &amountRemainingU)
		amountRemainingU.Neg(&amountRemainingU)
	}

	var maxFeeMinusFeePips uint256.Int
	maxFeeMinusFeePips.SetUint64(MaxFeeInt - uint64(feePips))
	if exactIn {
		var amountRemainingLessFee, tmp uint256.Int
		tmp.Mul(&amountRemainingU, &maxFeeMinusFeePips)
		amountRemainingLessFee.Div(&tmp, MaxFeeUint256)
		if zeroForOne {
			amountIn, err = GetAmount0DeltaV2(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, true)
			if err != nil {
				return
			}
		} else {
			amountIn, err = GetAmount1DeltaV2(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			sqrtRatioNextX96, err = GetNextSqrtPriceFromInput(sqrtRatioCurrentX96, liquidity, &amountRemainingLessFee, zeroForOne)
			if err != nil {
				return
			}
		}
	} else {
		if zeroForOne {
			amountOut, err = GetAmount1DeltaV2(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, false)
			if err != nil {
				return
			}
		} else {
			amountOut, err = GetAmount0DeltaV2(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, false)
			if err != nil {
				return
			}
		}
		if amountRemainingU.Cmp(amountOut) >= 0 {
			sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			sqrtRatioNextX96, err = GetNextSqrtPriceFromOutput(sqrtRatioCurrentX96, liquidity, &amountRemainingU, zeroForOne)
			if err != nil {
				return
			}
		}
	}

	max := sqrtRatioTargetX96.Cmp(sqrtRatioNextX96) == 0

	if zeroForOne {
		if !(max && exactIn) {
			amountIn, err = GetAmount0DeltaV2(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if !(max && !exactIn) {
			amountOut, err = GetAmount1DeltaV2(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, false)
			if err != nil {
				return
			}
		}
	} else {
		if !(max && exactIn) {
			amountIn, err = GetAmount1DeltaV2(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, true)
			if err != nil {
				return
			}
		}
		if !(max && !exactIn) {
			amountOut, err = GetAmount0DeltaV2(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, false)
			if err != nil {
				return
			}
		}
	}

	if !exactIn && amountOut.Cmp(&amountRemainingU) > 0 {
		amountOut = &amountRemainingU
	}

	if exactIn && sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		// we didn't reach the target, so take the remainder of the maximum input as fee
		feeAmount = new(uint256.Int).Sub(&amountRemainingU, amountIn)
	} else {
		if feeAmount == nil {
			feeAmount = new(Uint256)
		}
		err = MulDivRoundingUpV2(amountIn, uint256.NewInt(uint64(feePips)), &maxFeeMinusFeePips, feeAmount)
		if err != nil {
			return
		}
	}

	return
}
