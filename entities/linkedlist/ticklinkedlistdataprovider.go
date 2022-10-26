package linkedlist

import (
	"container/list"

	"github.com/KyberNetwork/uniswapv3-sdk/entities"
)

// TickListDataProvider A data provider for ticks that is backed by a doubly linked list of ticks.
type TickListDataProvider struct {
	ticksList *list.List
}

func NewTickLinkedListDataProvider(ticks []entities.Tick, tickSpacing int) (*TickListDataProvider, error) {
	if err := ValidateList(ticks, tickSpacing); err != nil {
		return nil, err
	}

	ticksList := list.New()

	for _, v := range ticks {
		_ = ticksList.PushBack(v)
	}

	return &TickListDataProvider{ticksList: ticksList}, nil
}

func (p *TickListDataProvider) GetTick(tick int) entities.Tick {
	return GetTick(p.ticksList, tick)
}

func (p *TickListDataProvider) NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool) {
	t := NextInitializedTick(p.ticksList, tick, lte)

	return t.Index, true
}
