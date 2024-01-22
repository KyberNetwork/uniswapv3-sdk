package utils

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMulDiv(t *testing.T) {
	// https://github.com/Uniswap/v3-core/blob/main/test/FullMath.spec.ts

	tests := []struct {
		a         string
		b         string
		deno      string
		expResult string
	}{
		{MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Dec()},
		{"0x100000000000000000000000000000000", "0x80000000000000000000000000000000", "0x180000000000000000000000000000000", "113427455640312821154458202477256070485"},
		{"0x100000000000000000000000000000000", "0x2300000000000000000000000000000000", "0x800000000000000000000000000000000", "1488735355279105777652263907513985925120"},
		{"0x100000000000000000000000000000000", "0x3e800000000000000000000000000000000", "0xbb800000000000000000000000000000000", "113427455640312821154458202477256070485"},

		{"0x61ae64157b363469ec1e000000000000000000000000", "0x5d5502f19f7baee2e5fa2", "0x69b797741ba66bda48a81e9", "126036350226489723925526476841950279379016090973169"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			r, err := MulDiv(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, r.Dec())

			// v2
			var rv2 Uint256
			err = MulDivV2(uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno), &rv2, nil)
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, rv2.Dec())
		})
	}

	failTests := []struct {
		a    string
		b    string
		deno string
	}{
		// {"0x100000000000000000000000000000000", "0x5", "0x0"}, // we don't catch div by zero here
		// {"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x0"},
		{"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x1"},
		{MaxUint256.Hex(), MaxUint256.Hex(), new(Uint256).SubUint64(MaxUint256, 1).Hex()},
	}
	for i, tt := range failTests {
		t.Run(fmt.Sprintf("fail test %d", i), func(t *testing.T) {
			_, err := MulDiv(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.NotNil(t, err)

			// v2
			var rv2 Uint256
			err = MulDivV2(uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno), &rv2, nil)
			require.NotNil(t, err)
		})
	}
}

func RandNumberHexString(maxLen int) string {
	sLen := rand.Intn(maxLen) + 1
	var s string
	for i := 0; i < sLen; i++ {
		var c int
		if i == 0 {
			c = rand.Intn(15) + 1
		} else {
			c = rand.Intn(16)
		}
		s = fmt.Sprintf("%s%x", s, c)
	}
	return s
}

func RandUint256() *Uint256 {
	s := RandNumberHexString(64)
	return uint256.MustFromHex("0x" + s)
}

func TestMulDivV2(t *testing.T) {
	for i := 0; i < 500; i++ {
		a := RandUint256()
		b := RandUint256()
		deno := RandUint256()

		t.Run(fmt.Sprintf("test %s %s %s", a.Hex(), b.Hex(), deno.Hex()), func(t *testing.T) {
			r, err := MulDiv(a, b, deno)

			var rv2 Uint256
			errv2 := MulDivV2(a, b, deno, &rv2, nil)

			if err != nil {
				require.NotNil(t, errv2)
			} else {
				require.Nil(t, errv2)
				assert.Equal(t, r.Dec(), rv2.Dec())
			}
		})
	}
}

func TestMulDivRoundingUp(t *testing.T) {
	// https://github.com/Uniswap/v3-core/blob/main/test/FullMath.spec.ts

	tests := []struct {
		a         string
		b         string
		deno      string
		expResult string
	}{
		{MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Hex(), MaxUint256.Dec()},
		{"0x100000000000000000000000000000000", "0x80000000000000000000000000000000", "0x180000000000000000000000000000000", "113427455640312821154458202477256070486"},
		{"0x100000000000000000000000000000000", "0x2300000000000000000000000000000000", "0x800000000000000000000000000000000", "1488735355279105777652263907513985925120"},
		{"0x100000000000000000000000000000000", "0x3e800000000000000000000000000000000", "0xbb800000000000000000000000000000000", "113427455640312821154458202477256070486"},

		{"0x2a60f4810d72e89eaee06f20122f1de80adc64777e121", "0xfd21718acef075500c6395ba922064220", "0xd195e7433221b9e4b6ef3f19b457c9c9797ae6b5eaacb402113dce147e97979f", "14406918379743960"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			r, err := MulDivRoundingUp(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, r.Dec())

			// v2
			var rv2 Uint256
			err = MulDivRoundingUpV2(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno), &rv2)
			require.Nil(t, err)
			assert.Equal(t, tt.expResult, rv2.Dec())
		})
	}

	failTests := []struct {
		a    string
		b    string
		deno string
	}{
		// {"0x100000000000000000000000000000000", "0x5", "0x0"}, // we don't catch div by zero here
		// {"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x0"},
		{"0x100000000000000000000000000000000", "0x100000000000000000000000000000000", "0x1"},
		{MaxUint256.Hex(), MaxUint256.Hex(), new(Uint256).SubUint64(MaxUint256, 1).Hex()},
		{"0x1e695d2db4f97", "0x10d5effea103c44aaf18a26b449186a7de3dd6c1ce3d26d03dfd9", "0x2"}, // mulDiv overflows 256 bits after rounding up
		{"0xffffffffffffffffffffffffffffffffffffffb07f6d608e4dcc38020b140b35", "0xffffffffffffffffffffffffffffffffffffffb07f6d608e4dcc38020b140b36", "0xffffffffffffffffffffffffffffffffffffff60fedac11c9b9870041628166c"}, // mulDiv overflows 256 bits after rounding up case 2
	}
	for i, tt := range failTests {
		t.Run(fmt.Sprintf("fail test %d", i), func(t *testing.T) {
			x, err := MulDivRoundingUp(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno))
			require.NotNil(t, err, x)

			// v2
			var rv2 Uint256
			err = MulDivRoundingUpV2(
				uint256.MustFromHex(tt.a), uint256.MustFromHex(tt.b),
				uint256.MustFromHex(tt.deno), &rv2)
			require.NotNil(t, err)
		})
	}
}

func TestMulDivRoundingUpV2(t *testing.T) {
	for i := 0; i < 500; i++ {
		a := RandUint256()
		b := RandUint256()
		deno := RandUint256()

		t.Run(fmt.Sprintf("test %s %s %s", a.Hex(), b.Hex(), deno.Hex()), func(t *testing.T) {
			r, err := MulDivRoundingUp(a, b, deno)

			var rv2 Uint256
			errv2 := MulDivRoundingUpV2(a, b, deno, &rv2)

			if err != nil {
				require.NotNil(t, errv2)
			} else {
				require.Nil(t, errv2)
				assert.Equal(t, r.Dec(), rv2.Dec())
			}
		})
	}
}
