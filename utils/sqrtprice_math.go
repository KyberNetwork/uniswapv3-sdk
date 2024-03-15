package utils

import (
	"errors"
	"math/big"

	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/holiman/uint256"
)

var (
	ErrSqrtPriceLessThanZero = errors.New("sqrt price less than zero")
	ErrLiquidityLessThanZero = errors.New("liquidity less than zero")
	ErrInvariant             = errors.New("invariant violation")
	ErrAddOverflow           = errors.New("add overflow")
)

var MaxUint160 = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffff")

func multiplyIn256(x, y, product *uint256.Int) *uint256.Int {
	product.Mul(x, y)
	return product // no need to And with MaxUint256 here
}

func addIn256(x, y, sum *uint256.Int) *uint256.Int {
	sum.Add(x, y)
	return sum // no need to And with MaxUint256 here
}

// deprecated
func GetAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	var res uint256.Int
	err := GetAmount0DeltaV2(
		uint256.MustFromBig(sqrtRatioAX96),
		uint256.MustFromBig(sqrtRatioBX96),
		uint256.MustFromBig(liquidity),
		roundUp,
		&res,
	)
	if err != nil {
		panic(err)
	}
	return res.ToBig()
}

func GetAmount0DeltaV2(sqrtRatioAX96, sqrtRatioBX96 *Uint160, liquidity *Uint128, roundUp bool, result *Uint256) error {
	// https://github.com/Uniswap/v3-core/blob/d8b1c635c275d2a9450bd6a78f3fa2484fef73eb/contracts/libraries/SqrtPriceMath.sol#L159
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	var numerator1, numerator2 uint256.Int
	numerator1.Lsh(liquidity, 96)
	numerator2.Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		var deno Uint256
		err := MulDivRoundingUpV2(&numerator1, &numerator2, sqrtRatioBX96, &deno)
		if err != nil {
			return err
		}
		DivRoundingUp(&deno, sqrtRatioAX96, result)
		return nil
	}
	// : FullMath.mulDiv(numerator1, numerator2, sqrtRatioBX96) / sqrtRatioAX96;
	var tmp Uint256
	err := MulDivV2(&numerator1, &numerator2, sqrtRatioBX96, &tmp, nil)
	if err != nil {
		return err
	}
	result.Div(&tmp, sqrtRatioAX96)
	return nil
}

// deprecated
func GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	var res Uint256
	err := GetAmount1DeltaV2(
		uint256.MustFromBig(sqrtRatioAX96),
		uint256.MustFromBig(sqrtRatioBX96),
		uint256.MustFromBig(liquidity),
		roundUp,
		&res,
	)
	if err != nil {
		panic(err)
	}
	return res.ToBig()
}

func GetAmount1DeltaV2(sqrtRatioAX96, sqrtRatioBX96 *Uint160, liquidity *Uint128, roundUp bool, result *Uint256) error {
	// https://github.com/Uniswap/v3-core/blob/d8b1c635c275d2a9450bd6a78f3fa2484fef73eb/contracts/libraries/SqrtPriceMath.sol#L188
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	var diff uint256.Int
	diff.Sub(sqrtRatioBX96, sqrtRatioAX96)
	if roundUp {
		err := MulDivRoundingUpV2(liquidity, &diff, constants.Q96U256, result)
		if err != nil {
			return err
		}
		return nil
	}
	// : FullMath.mulDiv(liquidity, sqrtRatioBX96 - sqrtRatioAX96, FixedPoint96.Q96);
	err := MulDivV2(liquidity, &diff, constants.Q96U256, result, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetNextSqrtPriceFromInput(sqrtPX96 *Uint160, liquidity *Uint128, amountIn *uint256.Int, zeroForOne bool, result *Uint160) error {
	if sqrtPX96.Sign() <= 0 {
		return ErrSqrtPriceLessThanZero
	}
	if liquidity.Sign() <= 0 {
		return ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true, result)
	}
	return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true, result)
}

func GetNextSqrtPriceFromOutput(sqrtPX96 *Uint160, liquidity *Uint128, amountOut *uint256.Int, zeroForOne bool, result *Uint160) error {
	if sqrtPX96.Sign() <= 0 {
		return ErrSqrtPriceLessThanZero
	}
	if liquidity.Sign() <= 0 {
		return ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false, result)
	}
	return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false, result)
}

func getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96 *Uint160, liquidity *Uint128, amount *uint256.Int, add bool, result *Uint160) error {
	if amount.IsZero() {
		result.Set(sqrtPX96)
		return nil
	}

	var numerator1, denominator, product, tmp uint256.Int
	numerator1.Lsh(liquidity, 96)
	multiplyIn256(amount, sqrtPX96, &product)
	if add {
		if tmp.Div(&product, amount).Cmp(sqrtPX96) == 0 {
			addIn256(&numerator1, &product, &denominator)
			if denominator.Cmp(&numerator1) >= 0 {
				err := MulDivRoundingUpV2(&numerator1, sqrtPX96, &denominator, result)
				return err
			}
		}
		tmp.Div(&numerator1, sqrtPX96)
		tmp.Add(&tmp, amount)
		DivRoundingUp(&numerator1, &tmp, result)
		return nil
	} else {
		if tmp.Div(&product, amount).Cmp(sqrtPX96) != 0 {
			return ErrInvariant
		}
		if numerator1.Cmp(&product) <= 0 {
			return ErrInvariant
		}
		denominator.Sub(&numerator1, &product)
		err := MulDivRoundingUpV2(&numerator1, sqrtPX96, &denominator, result)
		return err
	}
}

func getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96 *Uint160, liquidity *Uint128, amount *uint256.Int, add bool, result *Uint160) error {
	if add {
		var quotient, tmp uint256.Int
		if amount.Cmp(MaxUint160) <= 0 {
			tmp.Lsh(amount, 96)
			quotient.Div(&tmp, liquidity)
		} else {
			tmp.Mul(amount, constants.Q96U256)
			quotient.Div(&tmp, liquidity)
		}
		_, overflow := quotient.AddOverflow(&quotient, sqrtPX96)
		if overflow {
			return ErrAddOverflow
		}
		err := CheckToUint160(&quotient)
		if err != nil {
			return err
		}
		result.Set(&quotient)
		return nil
	}

	var quotient Uint256
	err := MulDivRoundingUpV2(amount, constants.Q96U256, liquidity, &quotient)
	if err != nil {
		return err
	}
	if sqrtPX96.Cmp(&quotient) <= 0 {
		return ErrInvariant
	}
	quotient.Sub(sqrtPX96, &quotient)
	// always fits 160 bits
	result.Set(&quotient)
	return nil
}
