package entities

import (
	"math/big"
	"testing"

	"github.com/KyberNetwork/uniswapv3-sdk-uint256/utils"
	"github.com/stretchr/testify/assert"
)

var (
	lowTick = Tick{
		Index:          utils.MinTick + 1,
		LiquidityNet:   big.NewInt(10),
		LiquidityGross: big.NewInt(10),
	}
	midTick = Tick{
		Index:          0,
		LiquidityNet:   big.NewInt(-5),
		LiquidityGross: big.NewInt(5),
	}
	highTick = Tick{
		Index:          utils.MaxTick - 1,
		LiquidityNet:   big.NewInt(-5),
		LiquidityGross: big.NewInt(5),
	}
)

func TestValidateList(t *testing.T) {
	assert.ErrorIs(t, ValidateList([]Tick{lowTick}, 1), ErrZeroNet, "panics for incomplete lists")
	assert.ErrorIs(t, ValidateList([]Tick{highTick, lowTick, midTick}, 1), ErrSorted, "panics for unsorted lists")
	assert.ErrorIs(t, ValidateList([]Tick{highTick, midTick, lowTick}, 1337), ErrInvalidTickSpacing, "errors if ticks are not on multiples of tick spacing")
}

func TestIsBelowSmallest(t *testing.T) {
	result := []Tick{lowTick, midTick, highTick}
	isBelowSmallest1, _ := IsBelowSmallest(result, utils.MinTick)
	assert.True(t, isBelowSmallest1)

	isBelowSmallest2, _ := IsBelowSmallest(result, utils.MinTick+1)
	assert.False(t, isBelowSmallest2)
}

func TestIsAtOrAboveSmallest(t *testing.T) {
	result := []Tick{lowTick, midTick, highTick}

	isAtOrAboveLargest1, _ := IsAtOrAboveLargest(result, utils.MaxTick-2)
	assert.False(t, isAtOrAboveLargest1)

	isAtOrAboveLargest2, _ := IsAtOrAboveLargest(result, utils.MaxTick-1)
	assert.True(t, isAtOrAboveLargest2)
}

func TestNextInitializedTick(t *testing.T) {
	ticks := []Tick{lowTick, midTick, highTick}

	type args struct {
		ticks []Tick
		tick  int
		lte   bool
	}
	tests := []struct {
		name string
		args args
		want Tick
	}{
		{name: "low - lte = true 0", args: args{ticks: ticks, tick: utils.MinTick + 1, lte: true}, want: lowTick},
		{name: "low - lte = true 1", args: args{ticks: ticks, tick: utils.MinTick + 2, lte: true}, want: lowTick},
		{name: "low - lte = false 0", args: args{ticks: ticks, tick: utils.MinTick, lte: false}, want: lowTick},
		{name: "low - lte = false 1", args: args{ticks: ticks, tick: utils.MinTick + 1, lte: false}, want: midTick},
		{name: "mid - lte = true 0", args: args{ticks: ticks, tick: 0, lte: true}, want: midTick},
		{name: "mid - lte = true 1", args: args{ticks: ticks, tick: 1, lte: true}, want: midTick},
		{name: "mid - lte = false 0", args: args{ticks: ticks, tick: -1, lte: false}, want: midTick},
		{name: "mid - lte = false 1", args: args{ticks: ticks, tick: 0 + 1, lte: false}, want: highTick},
		{name: "high - lte = true 0", args: args{ticks: ticks, tick: utils.MaxTick - 1, lte: true}, want: highTick},
		{name: "high - lte = true 1", args: args{ticks: ticks, tick: utils.MaxTick, lte: true}, want: highTick},
		{name: "high - lte = false 0", args: args{ticks: ticks, tick: utils.MaxTick - 2, lte: false}, want: highTick},
		{name: "high - lte = false 1", args: args{ticks: ticks, tick: utils.MaxTick - 3, lte: false}, want: highTick},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextInitializedTick, _ := NextInitializedTick(tt.args.ticks, tt.args.tick, tt.args.lte)
			assert.Equal(t, tt.want, nextInitializedTick)
		})
	}

	nextInitializedTick1, err1 := NextInitializedTick(ticks, utils.MinTick, true)
	assert.Zero(t, nextInitializedTick1, "below smallest")
	assert.ErrorIs(t, err1, ErrBelowSmallest)

	nextInitializedTick2, err2 := NextInitializedTick(ticks, utils.MaxTick-1, false)
	assert.Zero(t, nextInitializedTick2, "at or above largest")
	assert.ErrorIs(t, err2, ErrAtOrAboveLargest)
}

func TestNextInitializedTickWithinOneWord(t *testing.T) {
	ticks := []Tick{lowTick, midTick, highTick}

	// words around 0, lte = true
	type args struct {
		ticks       []Tick
		tick        int
		lte         bool
		tickSpacing int
	}
	tests := []struct {
		name  string
		args  args
		want0 int
		want1 bool
	}{
		// words around 0, lte = true
		{name: "lte = true  0", args: args{ticks: ticks, tick: -257, lte: true, tickSpacing: 1}, want0: -512, want1: false},
		{name: "lte = true  1", args: args{ticks: ticks, tick: -256, lte: true, tickSpacing: 1}, want0: -256, want1: false},
		{name: "lte = true  2", args: args{ticks: ticks, tick: -1, lte: true, tickSpacing: 1}, want0: -256, want1: false},
		{name: "lte = true  3", args: args{ticks: ticks, tick: 0, lte: true, tickSpacing: 1}, want0: 0, want1: true},
		{name: "lte = true  4", args: args{ticks: ticks, tick: 1, lte: true, tickSpacing: 1}, want0: 0, want1: true},
		{name: "lte = true  5", args: args{ticks: ticks, tick: 255, lte: true, tickSpacing: 1}, want0: 0, want1: true},
		{name: "lte = true  6", args: args{ticks: ticks, tick: 256, lte: true, tickSpacing: 1}, want0: 256, want1: false},
		{name: "lte = true  7", args: args{ticks: ticks, tick: 257, lte: true, tickSpacing: 1}, want0: 256, want1: false},

		// words around 0, lte = false
		{name: "lte = false 0", args: args{ticks: ticks, tick: -258, lte: false, tickSpacing: 1}, want0: -257, want1: false},
		{name: "lte = false 1", args: args{ticks: ticks, tick: -257, lte: false, tickSpacing: 1}, want0: -1, want1: false},
		{name: "lte = false 2", args: args{ticks: ticks, tick: -256, lte: false, tickSpacing: 1}, want0: -1, want1: false},
		{name: "lte = false 3", args: args{ticks: ticks, tick: -2, lte: false, tickSpacing: 1}, want0: -1, want1: false},
		{name: "lte = false 4", args: args{ticks: ticks, tick: -1, lte: false, tickSpacing: 1}, want0: 0, want1: true},
		{name: "lte = false 5", args: args{ticks: ticks, tick: 0, lte: false, tickSpacing: 1}, want0: 255, want1: false},
		{name: "lte = false 6", args: args{ticks: ticks, tick: 1, lte: false, tickSpacing: 1}, want0: 255, want1: false},
		{name: "lte = false 7", args: args{ticks: ticks, tick: 254, lte: false, tickSpacing: 1}, want0: 255, want1: false},
		{name: "lte = false 8", args: args{ticks: ticks, tick: 255, lte: false, tickSpacing: 1}, want0: 511, want1: false},
		{name: "lte = false 9", args: args{ticks: ticks, tick: 256, lte: false, tickSpacing: 1}, want0: 511, want1: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got0, got1, _ := NextInitializedTickWithinOneWord(tt.args.ticks, tt.args.tick, tt.args.lte, tt.args.tickSpacing)
			assert.Equal(t, tt.want0, got0)
			assert.Equal(t, tt.want1, got1)
		})
	}

}
