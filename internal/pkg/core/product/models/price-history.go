package models

import (
	"sort"
	"strings"
)

type PriceHistory []PriceTimeStamp

func (p PriceHistory) GetSorted() PriceHistory {
	result := make(PriceHistory, len(p), len(p))
	copy(result, p)
	sort.Stable(result)
	return result
}

func (p PriceHistory) String() string {
	if len(p) == 0 {
		return "no price records"
	}
	sorted := p.GetSorted()
	return "last price: " + sorted[len(sorted)-1].String()
}

func (p PriceHistory) FullHistoryString() string {
	if p.Len() == 0 {
		return "no price records"
	}
	sorted := p.GetSorted()
	stringData := make([]string, 0, len(sorted))
	for _, priceTimeStamp := range sorted {
		stringData = append(stringData, priceTimeStamp.String())
	}
	return strings.Join(stringData, "\n")
}

func (p PriceHistory) Len() int {
	return len(p)
}

func (p PriceHistory) Less(i, j int) bool {
	return p[j].Date.After(p[i].Date)
}

func (p PriceHistory) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
