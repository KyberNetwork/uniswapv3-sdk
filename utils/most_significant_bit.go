package utils

import (
	"errors"

	"github.com/holiman/uint256"
)

var ErrInvalidInput = errors.New("invalid input")

type powerOf2 struct {
	power uint
	value *uint256.Int
}

var powersOf2 = []powerOf2{
	{128, uint256.MustFromHex("0x100000000000000000000000000000000")},
	{64, uint256.MustFromHex("0x10000000000000000")},
	{32, uint256.MustFromHex("0x100000000")},
	{16, uint256.MustFromHex("0x10000")},
	{8, uint256.MustFromHex("0x100")},
	{4, uint256.MustFromHex("0x10")},
	{2, uint256.MustFromHex("0x4")},
	{1, uint256.MustFromHex("0x2")},
}

func MostSignificantBit(x *uint256.Int) (uint, error) {
	if x.Sign() == 0 {
		return 0, ErrInvalidInput
	}

	var tmpX uint256.Int
	tmpX.Set(x)
	var msb uint
	for _, p := range powersOf2 {
		if tmpX.Cmp(p.value) >= 0 {
			tmpX.Rsh(&tmpX, p.power)
			msb += p.power
		}
	}
	return msb, nil
}
