package entities

// A data provider for ticks that is backed by an in-memory array of ticks.
type TickListDataProvider struct {
	ticks []Tick
}

func NewTickListDataProvider(ticks []Tick, tickSpacing int) (*TickListDataProvider, error) {
	if err := ValidateList(ticks, tickSpacing); err != nil {
		return nil, err
	}
	return &TickListDataProvider{ticks: ticks}, nil
}

func (p *TickListDataProvider) GetTick(tick int) (Tick, error) {
	return GetTick(p.ticks, tick)
}

func (p *TickListDataProvider) NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool, error) {
	return NextInitializedTickWithinOneWord(p.ticks, tick, lte, tickSpacing)
}

func (p *TickListDataProvider) NextInitializedTickIndex(tick int, lte bool) (int, bool, error) {
	return NextInitializedTickIndex(p.ticks, tick, lte)
}
