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
	res, err := GetAmount0DeltaV2(
		uint256.MustFromBig(sqrtRatioAX96),
		uint256.MustFromBig(sqrtRatioBX96),
		uint256.MustFromBig(liquidity),
		roundUp,
	)
	if err != nil {
		panic(err)
	}
	return res.ToBig()
}

func GetAmount0DeltaV2(sqrtRatioAX96, sqrtRatioBX96 *Uint160, liquidity *Uint128, roundUp bool) (*uint256.Int, error) {
	// https://github.com/Uniswap/v3-core/blob/d8b1c635c275d2a9450bd6a78f3fa2484fef73eb/contracts/libraries/SqrtPriceMath.sol#L159
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	var numerator1, numerator2 uint256.Int
	numerator1.Lsh(liquidity, 96)
	numerator2.Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		deno, err := MulDivRoundingUp(&numerator1, &numerator2, sqrtRatioBX96)
		if err != nil {
			return nil, err
		}
		return DivRoundingUp(deno, sqrtRatioAX96), nil
	}
	// : FullMath.mulDiv(numerator1, numerator2, sqrtRatioBX96) / sqrtRatioAX96;
	tmp, err := MulDiv(&numerator1, &numerator2, sqrtRatioBX96)
	if err != nil {
		return nil, err
	}
	result := new(uint256.Int).Div(tmp, sqrtRatioAX96)
	return result, nil
}

// deprecated
func GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	res, err := GetAmount1DeltaV2(
		uint256.MustFromBig(sqrtRatioAX96),
		uint256.MustFromBig(sqrtRatioBX96),
		uint256.MustFromBig(liquidity),
		roundUp,
	)
	if err != nil {
		panic(err)
	}
	return res.ToBig()
}

func GetAmount1DeltaV2(sqrtRatioAX96, sqrtRatioBX96 *Uint160, liquidity *Uint128, roundUp bool) (*uint256.Int, error) {
	// https://github.com/Uniswap/v3-core/blob/d8b1c635c275d2a9450bd6a78f3fa2484fef73eb/contracts/libraries/SqrtPriceMath.sol#L188
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	var diff uint256.Int
	diff.Sub(sqrtRatioBX96, sqrtRatioAX96)
	if roundUp {
		return MulDivRoundingUp(liquidity, &diff, constants.Q96U256)
	}
	// : FullMath.mulDiv(liquidity, sqrtRatioBX96 - sqrtRatioAX96, FixedPoint96.Q96);
	return MulDiv(liquidity, &diff, constants.Q96U256)
}

func GetNextSqrtPriceFromInput(sqrtPX96 *Uint160, liquidity *Uint128, amountIn *uint256.Int, zeroForOne bool) (*Uint160, error) {
	if sqrtPX96.Sign() <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Sign() <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true)
	}
	return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true)
}

func GetNextSqrtPriceFromOutput(sqrtPX96 *Uint160, liquidity *Uint128, amountOut *uint256.Int, zeroForOne bool) (*Uint160, error) {
	if sqrtPX96.Sign() <= 0 {
		return nil, ErrSqrtPriceLessThanZero
	}
	if liquidity.Sign() <= 0 {
		return nil, ErrLiquidityLessThanZero
	}
	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false)
	}
	return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false)
}

func getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96 *Uint160, liquidity *Uint128, amount *uint256.Int, add bool) (*Uint160, error) {
	if amount.IsZero() {
		return sqrtPX96, nil
	}

	var numerator1, denominator, product, tmp uint256.Int
	numerator1.Lsh(liquidity, 96)
	multiplyIn256(amount, sqrtPX96, &product)
	if add {
		if tmp.Div(&product, amount).Cmp(sqrtPX96) == 0 {
			addIn256(&numerator1, &product, &denominator)
			if denominator.Cmp(&numerator1) >= 0 {
				return MulDivRoundingUp(&numerator1, sqrtPX96, &denominator)
			}
		}
		tmp.Div(&numerator1, sqrtPX96)
		tmp.Add(&tmp, amount)
		return DivRoundingUp(&numerator1, &tmp), nil
	} else {
		if tmp.Div(&product, amount).Cmp(sqrtPX96) != 0 {
			return nil, ErrInvariant
		}
		if numerator1.Cmp(&product) <= 0 {
			return nil, ErrInvariant
		}
		denominator.Sub(&numerator1, &product)
		return MulDivRoundingUp(&numerator1, sqrtPX96, &denominator)
	}
}

func getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96 *Uint160, liquidity *Uint128, amount *uint256.Int, add bool) (*Uint160, error) {
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
			return nil, ErrAddOverflow
		}
		err := CheckToUint160(&quotient)
		if err != nil {
			return nil, err
		}
		return &quotient, nil
	}

	quotient, err := MulDivRoundingUp(amount, constants.Q96U256, liquidity)
	if err != nil {
		return nil, err
	}
	if sqrtPX96.Cmp(quotient) <= 0 {
		return nil, ErrInvariant
	}
	quotient.Sub(sqrtPX96, quotient)
	// always fits 160 bits
	return quotient, nil
}
