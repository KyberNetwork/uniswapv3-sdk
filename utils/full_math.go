package utils

import (
	"errors"
	"math/big"

	"github.com/holiman/uint256"
)

var (
	ErrMulDivOverflow = errors.New("muldiv overflow")
	One               = big.NewInt(1)
)

// Calculates ceil(a×b÷denominator) with full precision
func MulDivRoundingUp(a, b, denominator *uint256.Int) (*uint256.Int, error) {
	// the product can overflow so need to use big.Int here
	// TODO: optimize this
	var product, rem, result big.Int
	product.Mul(a.ToBig(), b.ToBig())
	result.DivMod(&product, denominator.ToBig(), &rem)
	if rem.Sign() != 0 {
		result.Add(&result, One)
	}

	resultU, overflow := uint256.FromBig(&result)
	if overflow {
		return nil, ErrMulDivOverflow
	}
	return resultU, nil
}

// Calculates floor(a×b÷denominator) with full precision
func MulDiv(a, b, denominator *uint256.Int) (*uint256.Int, error) {
	// the product can overflow so need to use big.Int here
	// TODO: optimize this follow univ3 code
	var product, result big.Int
	product.Mul(a.ToBig(), b.ToBig())
	result.Div(&product, denominator.ToBig())

	resultU, overflow := uint256.FromBig(&result)
	if overflow {
		return nil, ErrMulDivOverflow
	}
	return resultU, nil
}

// Returns ceil(x / y)
func DivRoundingUp(a, denominator *uint256.Int) *uint256.Int {
	var result, rem uint256.Int
	result.DivMod(a, denominator, &rem)
	if !rem.IsZero() {
		result.AddUint64(&result, 1)
	}
	return &result
}
