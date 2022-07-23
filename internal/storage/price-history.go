package storage

import (
	"sort"
	"strings"
)

type PriceHistory struct {
	data   []*PriceTimeStamp
	sorted bool
}

func NewPriceHistory() *PriceHistory {
	return &PriceHistory{
		data:   make([]*PriceTimeStamp, 0),
		sorted: true,
	}
}

func (p PriceHistory) Len() int {
	return len(p.data)
}

func (p PriceHistory) Less(i, j int) bool {
	return p.data[j].date.After(p.data[i].date)
}

func (p PriceHistory) Swap(i, j int) {
	p.data[i], p.data[j] = p.data[j], p.data[i]
}

func (p *PriceHistory) AddPriceTimeStamp(stamp *PriceTimeStamp) {
	p.sorted = false
	p.data = append(p.data, stamp)
}

func (p *PriceHistory) GetData() []*PriceTimeStamp {
	if !p.sorted {
		sort.Stable(p)
		p.sorted = true
	}
	return p.data
}

func (p *PriceHistory) String() string {
	if len(p.data) == 0 {
		return "no price records"
	}
	data := p.GetData()
	return "last price: " + data[len(data)-1].String()
}

func (p *PriceHistory) FullHistoryString() string {
	data := p.GetData()
	stringData := make([]string, 0, len(data))
	for _, priceTimeStamp := range data {
		stringData = append(stringData, priceTimeStamp.String())
	}
	return strings.Join(stringData, "\n")
}
