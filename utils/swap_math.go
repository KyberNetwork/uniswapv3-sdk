package utils

import (
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/holiman/uint256"
)

const MaxFeeInt = 1000000

var MaxFeeUint256 = uint256.NewInt(MaxFeeInt)

func ComputeSwapStep(
	sqrtRatioCurrentX96,
	sqrtRatioTargetX96 *Uint160,
	liquidity *Uint128,
	amountRemaining *Int256,
	feePips constants.FeeAmount,

	sqrtRatioNextX96 *Uint160, amountIn, amountOut, feeAmount *Uint256,
) error {
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
			err := GetAmount0DeltaV2(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, true, amountIn)
			if err != nil {
				return err
			}
		} else {
			err := GetAmount1DeltaV2(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, true, amountIn)
			if err != nil {
				return err
			}
		}
		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX96.Set(sqrtRatioTargetX96)
		} else {
			err := GetNextSqrtPriceFromInput(sqrtRatioCurrentX96, liquidity, &amountRemainingLessFee, zeroForOne, sqrtRatioNextX96)
			if err != nil {
				return err
			}
		}
	} else {
		if zeroForOne {
			err := GetAmount1DeltaV2(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, false, amountOut)
			if err != nil {
				return err
			}
		} else {
			err := GetAmount0DeltaV2(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, false, amountOut)
			if err != nil {
				return err
			}
		}
		if amountRemainingU.Cmp(amountOut) >= 0 {
			sqrtRatioNextX96.Set(sqrtRatioTargetX96)
		} else {
			err := GetNextSqrtPriceFromOutput(sqrtRatioCurrentX96, liquidity, &amountRemainingU, zeroForOne, sqrtRatioNextX96)
			if err != nil {
				return err
			}
		}
	}

	max := sqrtRatioTargetX96.Cmp(sqrtRatioNextX96) == 0

	if zeroForOne {
		if !(max && exactIn) {
			err := GetAmount0DeltaV2(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, true, amountIn)
			if err != nil {
				return err
			}
		}
		if !(max && !exactIn) {
			err := GetAmount1DeltaV2(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, false, amountOut)
			if err != nil {
				return err
			}
		}
	} else {
		if !(max && exactIn) {
			err := GetAmount1DeltaV2(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, true, amountIn)
			if err != nil {
				return err
			}
		}
		if !(max && !exactIn) {
			err := GetAmount0DeltaV2(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, false, amountOut)
			if err != nil {
				return err
			}
		}
	}

	if !exactIn && amountOut.Cmp(&amountRemainingU) > 0 {
		amountOut.Set(&amountRemainingU)
	}

	if exactIn && sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		// we didn't reach the target, so take the remainder of the maximum input as fee
		feeAmount.Sub(&amountRemainingU, amountIn)
	} else {
		err := MulDivRoundingUpV2(amountIn, uint256.NewInt(uint64(feePips)), &maxFeeMinusFeePips, feeAmount)
		if err != nil {
			return err
		}
	}

	return nil
}
