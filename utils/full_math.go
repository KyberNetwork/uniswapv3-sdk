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

func MulDivRoundingUpV2(a, b, denominator, result *uint256.Int) error {
	var remainder Uint256
	err := MulDivV2(a, b, denominator, result, &remainder)
	if err != nil {
		return err
	}

	if !remainder.IsZero() {
		if result.Cmp(MaxUint256) == 0 {
			return ErrInvariant
		}
		result.AddUint64(result, 1)
	}
	return nil
}

// result=floor(a×b÷denominator), remainder=a×b%denominator
// (pass remainder=nil if not required)
func MulDivV2(a, b, denominator, result, remainder *uint256.Int) error {
	// https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/FullMath.sol
	// 512-bit multiply [prod1 prod0] = a * b
	// Compute the product mod 2**256 and mod 2**256 - 1
	// then use the Chinese Remainder Theorem to reconstruct
	// the 512 bit result. The result is stored in two 256
	// variables such that product = prod1 * 2**256 + prod0
	var prod0 Uint256 // Least significant 256 bits of the product
	var prod1 Uint256 // Most significant 256 bits of the product

	var denominatorTmp Uint256 // temp var (need to modify denominator along the way)
	denominatorTmp.Set(denominator)

	var mm Uint256
	mm.MulMod(a, b, MaxUint256)
	prod0.Mul(a, b)
	prod1.Sub(&mm, &prod0)
	if mm.Cmp(&prod0) < 0 {
		prod1.SubUint64(&prod1, 1)
	}

	// Handle non-overflow cases, 256 by 256 division
	if prod1.IsZero() {
		if denominatorTmp.IsZero() {
			return ErrInvariant
		}

		if remainder != nil {
			remainder.MulMod(a, b, &denominatorTmp)
		}
		result.Div(&prod0, &denominatorTmp)
		return nil
	}

	// Make sure the result is less than 2**256.
	// Also prevents denominator == 0
	if denominatorTmp.Cmp(&prod1) <= 0 {
		return ErrInvariant
	}

	///////////////////////////////////////////////
	// 512 by 256 division.
	///////////////////////////////////////////////

	// Make division exact by subtracting the remainder from [prod1 prod0]
	// Compute remainder using mulmod
	if remainder == nil {
		var remainderTmp Uint256
		remainder = &remainderTmp
	}
	remainder.MulMod(a, b, &denominatorTmp)
	// Subtract 256 bit number from 512 bit number
	if remainder.Cmp(&prod0) > 0 {
		prod1.SubUint64(&prod1, 1)
	}
	prod0.Sub(&prod0, remainder)

	// Factor powers of two out of denominator
	// Compute largest power of two divisor of denominator.
	// Always >= 1.
	var twos, tmp, tmp1, zero, two, three Uint256
	twos.And(tmp.Neg(&denominatorTmp), &denominatorTmp)
	// Divide denominator by power of two
	denominatorTmp.Div(&denominatorTmp, &twos)

	// Divide [prod1 prod0] by the factors of two
	prod0.Div(&prod0, &twos)
	// Shift in bits from prod1 into prod0. For this we need
	// to flip `twos` such that it is 2**256 / twos.
	// If twos is zero, then it becomes one
	zero.Clear()
	twos.AddUint64(tmp.Div(tmp1.Sub(&zero, &twos), &twos), 1)
	prod0.Or(&prod0, tmp.Mul(&prod1, &twos))

	// Invert denominator mod 2**256
	// Now that denominator is an odd number, it has an inverse
	// modulo 2**256 such that denominator * inv = 1 mod 2**256.
	// Compute the inverse by starting with a seed that is correct
	// correct for four bits. That is, denominator * inv = 1 mod 2**4
	var inv Uint256
	two.SetUint64(2)
	three.SetUint64(3)
	inv.Xor(tmp.Mul(&denominatorTmp, &three), &two)
	// Now use Newton-Raphson iteration to improve the precision.
	// Thanks to Hensel's lifting lemma, this also works in modular
	// arithmetic, doubling the correct bits in each step.
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**8
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**16
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**32
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**64
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**128
	inv.Mul(&inv, tmp.Sub(&two, tmp1.Mul(&denominatorTmp, &inv))) // inverse mod 2**256

	// Because the division is now exact we can divide by multiplying
	// with the modular inverse of denominator. This will give us the
	// correct result modulo 2**256. Since the precoditions guarantee
	// that the outcome is less than 2**256, this is the final result.
	// We don't need to compute the high bits of the result and prod1
	// is no longer required.
	result.Mul(&prod0, &inv)
	return nil
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
func DivRoundingUp(a, denominator, result *uint256.Int) {
	var rem uint256.Int
	result.DivMod(a, denominator, &rem)
	if !rem.IsZero() {
		result.AddUint64(result, 1)
	}
}
