package utils

import (
	"errors"

	"github.com/KyberNetwork/int256"
	"github.com/holiman/uint256"
)

// define placeholders for these types, in case we need to customize them later
// (for example to add boundary check...)

type Uint256 = uint256.Int
type Uint160 = uint256.Int
type Uint128 = uint256.Int

type Int256 = int256.Int
type Int128 = int256.Int

var (
	ErrExceedMaxInt256 = errors.New("exceed max int256")
	ErrOverflowUint128 = errors.New("overflow uint128")
	ErrOverflowUint160 = errors.New("overflow uint160")

	Uint128Max = uint256.MustFromHex("0xffffffffffffffffffffffffffffffff")
	Uint160Max = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffff")
)

// https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/SafeCast.sol
func ToInt256(value *Uint256, result *Int256) error {
	// if value (interpreted as a two's complement signed number) is negative -> it must be larger than max int256
	if value.Sign() < 0 {
		return ErrExceedMaxInt256
	}
	var ba [32]byte
	value.WriteToArray32(&ba)
	result.SetBytes32(ba[:])
	return nil
}

func ToUInt256(value *Int256, result *Uint256) error {
	var ba [32]byte
	value.WriteToArray32(&ba)
	result.SetBytes32(ba[:])
	return nil
}

// https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/SafeCast.sol
func CheckToUint160(value *Uint256) error {
	// we're using same type for Uint256 and Uint160, so use the original for now
	if value.Cmp(Uint160Max) > 0 {
		return ErrOverflowUint160
	}
	return nil
}

// x = x + y
func AddDeltaInPlace(x *Uint128, y *Int128) error {
	// for now we're using int256 for Int128, and uint256 for Uint128
	// and both of them is using two's complement internally
	// so just cast `y` to uint256 and add them together
	var ba [32]byte
	y.WriteToArray32(&ba)
	var yuint Uint128
	yuint.SetBytes32(ba[:])
	x.Add(x, &yuint)

	if x.Cmp(Uint128Max) > 0 {
		// could be overflow or underflow
		return ErrOverflowUint128
	}
	return nil
}
