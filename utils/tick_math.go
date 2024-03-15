package utils

import (
	"errors"
	"math/big"

	"github.com/KyberNetwork/int256"
	"github.com/holiman/uint256"
)

const (
	MinTick = -887272  // The minimum tick that can be used on any pool.
	MaxTick = -MinTick // The maximum tick that can be used on any pool.
)

var (
	Q32             = big.NewInt(1 << 32)
	MinSqrtRatio    = big.NewInt(4295128739)                                                          // The sqrt ratio corresponding to the minimum tick that could be used on any pool.
	MaxSqrtRatio, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 10) // The sqrt ratio corresponding to the maximum tick that could be used on any pool.

	Q32U256          = uint256.NewInt(1 << 32)
	MinSqrtRatioU256 = uint256.NewInt(4295128739)                                                   // The sqrt ratio corresponding to the minimum tick that could be used on any pool.
	MaxSqrtRatioU256 = uint256.MustFromDecimal("1461446703485210103287273052203988822378723970342") // The sqrt ratio corresponding to the maximum tick that could be used on any pool.
)

var (
	ErrInvalidTick      = errors.New("invalid tick")
	ErrInvalidSqrtRatio = errors.New("invalid sqrt ratio")
)

func mulShift(val *Uint256, mulBy *Uint256) {
	var tmp Uint256
	val.Rsh(tmp.Mul(val, mulBy), 128)
}

var (
	sqrtConst1  = uint256.MustFromHex("0xfffcb933bd6fad37aa2d162d1a594001")
	sqrtConst2  = uint256.MustFromHex("0x100000000000000000000000000000000")
	sqrtConst3  = uint256.MustFromHex("0xfff97272373d413259a46990580e213a")
	sqrtConst4  = uint256.MustFromHex("0xfff2e50f5f656932ef12357cf3c7fdcc")
	sqrtConst5  = uint256.MustFromHex("0xffe5caca7e10e4e61c3624eaa0941cd0")
	sqrtConst6  = uint256.MustFromHex("0xffcb9843d60f6159c9db58835c926644")
	sqrtConst7  = uint256.MustFromHex("0xff973b41fa98c081472e6896dfb254c0")
	sqrtConst8  = uint256.MustFromHex("0xff2ea16466c96a3843ec78b326b52861")
	sqrtConst9  = uint256.MustFromHex("0xfe5dee046a99a2a811c461f1969c3053")
	sqrtConst10 = uint256.MustFromHex("0xfcbe86c7900a88aedcffc83b479aa3a4")
	sqrtConst11 = uint256.MustFromHex("0xf987a7253ac413176f2b074cf7815e54")
	sqrtConst12 = uint256.MustFromHex("0xf3392b0822b70005940c7a398e4b70f3")
	sqrtConst13 = uint256.MustFromHex("0xe7159475a2c29b7443b29c7fa6e889d9")
	sqrtConst14 = uint256.MustFromHex("0xd097f3bdfd2022b8845ad8f792aa5825")
	sqrtConst15 = uint256.MustFromHex("0xa9f746462d870fdf8a65dc1f90e061e5")
	sqrtConst16 = uint256.MustFromHex("0x70d869a156d2a1b890bb3df62baf32f7")
	sqrtConst17 = uint256.MustFromHex("0x31be135f97d08fd981231505542fcfa6")
	sqrtConst18 = uint256.MustFromHex("0x9aa508b5b7a84e1c677de54f3e99bc9")
	sqrtConst19 = uint256.MustFromHex("0x5d6af8dedb81196699c329225ee604")
	sqrtConst20 = uint256.MustFromHex("0x2216e584f5fa1ea926041bedfe98")
	sqrtConst21 = uint256.MustFromHex("0x48a170391f7dc42444e8fa2")

	MaxUint256 = uint256.MustFromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

// deprecated
func GetSqrtRatioAtTick(tick int) (*big.Int, error) {
	var res Uint160
	err := GetSqrtRatioAtTickV2(tick, &res)
	return res.ToBig(), err
}

/**
 * Returns the sqrt ratio as a Q64.96 for the given tick. The sqrt ratio is computed as sqrt(1.0001)^tick
 * @param tick the tick for which to compute the sqrt ratio
 */
func GetSqrtRatioAtTickV2(tick int, result *Uint160) error {
	if tick < MinTick || tick > MaxTick {
		return ErrInvalidTick
	}
	absTick := tick
	if tick < 0 {
		absTick = -tick
	}
	var ratio Uint256
	if absTick&0x1 != 0 {
		ratio.Set(sqrtConst1)
	} else {
		ratio.Set(sqrtConst2)
	}
	if (absTick & 0x2) != 0 {
		mulShift(&ratio, sqrtConst3)
	}
	if (absTick & 0x4) != 0 {
		mulShift(&ratio, sqrtConst4)
	}
	if (absTick & 0x8) != 0 {
		mulShift(&ratio, sqrtConst5)
	}
	if (absTick & 0x10) != 0 {
		mulShift(&ratio, sqrtConst6)
	}
	if (absTick & 0x20) != 0 {
		mulShift(&ratio, sqrtConst7)
	}
	if (absTick & 0x40) != 0 {
		mulShift(&ratio, sqrtConst8)
	}
	if (absTick & 0x80) != 0 {
		mulShift(&ratio, sqrtConst9)
	}
	if (absTick & 0x100) != 0 {
		mulShift(&ratio, sqrtConst10)
	}
	if (absTick & 0x200) != 0 {
		mulShift(&ratio, sqrtConst11)
	}
	if (absTick & 0x400) != 0 {
		mulShift(&ratio, sqrtConst12)
	}
	if (absTick & 0x800) != 0 {
		mulShift(&ratio, sqrtConst13)
	}
	if (absTick & 0x1000) != 0 {
		mulShift(&ratio, sqrtConst14)
	}
	if (absTick & 0x2000) != 0 {
		mulShift(&ratio, sqrtConst15)
	}
	if (absTick & 0x4000) != 0 {
		mulShift(&ratio, sqrtConst16)
	}
	if (absTick & 0x8000) != 0 {
		mulShift(&ratio, sqrtConst17)
	}
	if (absTick & 0x10000) != 0 {
		mulShift(&ratio, sqrtConst18)
	}
	if (absTick & 0x20000) != 0 {
		mulShift(&ratio, sqrtConst19)
	}
	if (absTick & 0x40000) != 0 {
		mulShift(&ratio, sqrtConst20)
	}
	if (absTick & 0x80000) != 0 {
		mulShift(&ratio, sqrtConst21)
	}
	if tick > 0 {
		result.Div(MaxUint256, &ratio)
		ratio.Set(result)
	}

	// back to Q96
	var rem Uint256
	result.DivMod(&ratio, Q32U256, &rem)
	if !rem.IsZero() {
		result.AddUint64(result, 1)
		return nil
	} else {
		return nil
	}
}

var (
	magicSqrt10001 = int256.MustFromDec("255738958999603826347141")
	magicTickLow   = int256.MustFromDec("3402992956809132418596140100660247210")
	magicTickHigh  = int256.MustFromDec("291339464771989622907027621153398088495")
)

// deprecated
func GetTickAtSqrtRatio(sqrtRatioX96 *big.Int) (int, error) {
	return GetTickAtSqrtRatioV2(uint256.MustFromBig(sqrtRatioX96))
}

/**
 * Returns the tick corresponding to a given sqrt ratio, s.t. #getSqrtRatioAtTick(tick) <= sqrtRatioX96
 * and #getSqrtRatioAtTick(tick + 1) > sqrtRatioX96
 * @param sqrtRatioX96 the sqrt ratio as a Q64.96 for which to compute the tick
 */
func GetTickAtSqrtRatioV2(sqrtRatioX96 *Uint160) (int, error) {
	if sqrtRatioX96.Cmp(MinSqrtRatioU256) < 0 || sqrtRatioX96.Cmp(MaxSqrtRatioU256) >= 0 {
		return 0, ErrInvalidSqrtRatio
	}
	var sqrtRatioX128 Uint256
	sqrtRatioX128.Lsh(sqrtRatioX96, 32)
	msb, err := MostSignificantBit(&sqrtRatioX128)
	if err != nil {
		return 0, err
	}
	var r Uint256
	if msb >= 128 {
		r.Rsh(&sqrtRatioX128, msb-127)
	} else {
		r.Lsh(&sqrtRatioX128, 127-msb)
	}

	log2 := int256.NewInt(int64(msb - 128))
	log2.Lsh(log2, 64)

	var tmp, f Uint256
	for i := 0; i < 14; i++ {
		tmp.Mul(&r, &r)
		r.Rsh(&tmp, 127)
		f.Rsh(&r, 128)
		tmp.Lsh(&f, uint(63-i))

		// this is for Or, so we can cast the underlying words directly without copying
		tmpsigned := (*int256.Int)(&tmp)

		log2.Or(log2, tmpsigned)
		r.Rsh(&r, uint(f.Uint64()))
	}

	var logSqrt10001, tmp1, tmp2 Int256
	logSqrt10001.Mul(log2, magicSqrt10001)

	tickLow := tmp2.Rsh(tmp1.Sub(&logSqrt10001, magicTickLow), 128).Uint64()
	tickHigh := tmp2.Rsh(tmp1.Add(&logSqrt10001, magicTickHigh), 128).Uint64()

	if tickLow == tickHigh {
		return int(tickLow), nil
	}

	var sqrtRatio Uint160
	err = GetSqrtRatioAtTickV2(int(tickHigh), &sqrtRatio)
	if err != nil {
		return 0, err
	}
	if sqrtRatio.Cmp(sqrtRatioX96) <= 0 {
		return int(tickHigh), nil
	} else {
		return int(tickLow), nil
	}
}
