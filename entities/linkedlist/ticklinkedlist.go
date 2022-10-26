package linkedlist

import (
	"container/list"
	"math/big"

	"github.com/KyberNetwork/uniswapv3-sdk/entities"
)

func ValidateList(ticks []entities.Tick, tickSpacing int) error {
	if tickSpacing <= 0 {
		return entities.ErrZeroTickSpacing
	}

	// ensure ticks are spaced appropriately
	for _, t := range ticks {
		if t.Index%tickSpacing != 0 {
			return entities.ErrInvalidTickSpacing
		}
	}

	// ensure tick liquidity deltas sum to 0
	sum := big.NewInt(0)
	for _, tick := range ticks {
		sum.Add(sum, tick.LiquidityNet)
	}
	if sum.Cmp(big.NewInt(0)) != 0 {
		return entities.ErrZeroNet
	}

	if !isTicksSorted(ticks) {
		return entities.ErrSorted
	}

	return nil
}

func IsBelowSmallest(ticksList *list.List, tick int) bool {
	if ticksList.Len() == 0 {
		panic("empty tick list")
	}

	return tick < ticksList.Front().Value.(entities.Tick).Index
}

func IsAtOrAboveLargest(ticksList *list.List, tick int) bool {
	if ticksList.Len() == 0 {
		panic("empty tick list")
	}

	return tick >= ticksList.Back().Value.(entities.Tick).Index
}

func GetTick(ticksList *list.List, index int) entities.Tick {
	for e := ticksList.Front(); e != nil; e = e.Next() {
		value := e.Value.(entities.Tick)
		if value.Index == index {
			return value
		}
	}

	panic("invalid tick index")
}

func NextInitializedTick(ticksList *list.List, tick int, lte bool) entities.Tick {
	if lte {
		if IsBelowSmallest(ticksList, tick) {
			panic("below smallest")
		}
		if IsAtOrAboveLargest(ticksList, tick) {
			return ticksList.Back().Value.(entities.Tick)
		}

		for e := ticksList.Front(); e != nil; e = e.Next() {
			value := e.Value.(entities.Tick)

			if value.Index >= tick {
				return e.Prev().Value.(entities.Tick)
			}
		}

		return ticksList.Back().Value.(entities.Tick)
	} else {
		if IsAtOrAboveLargest(ticksList, tick) {
			panic("at or above largest")
		}
		if IsBelowSmallest(ticksList, tick) {
			return ticksList.Front().Value.(entities.Tick)
		}

		for e := ticksList.Back(); e != nil; e = e.Prev() {
			value := e.Value.(entities.Tick)

			if value.Index <= tick {
				return e.Next().Value.(entities.Tick)
			}
		}

		return ticksList.Front().Value.(entities.Tick)
	}
}

// utils

func isTicksSorted(ticks []entities.Tick) bool {
	for i := 0; i < len(ticks)-1; i++ {
		if ticks[i].Index > ticks[i+1].Index {
			return false
		}
	}
	return true
}
