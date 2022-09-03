package models

import (
	"encoding/json"
	"strings"
)

type PriceHistory []PriceTimeStamp

func (p PriceHistory) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PriceHistory) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p PriceHistory) Copy() PriceHistory {
	res := make(PriceHistory, len(p), len(p))
	copy(res, p)
	return res
}

func (p PriceHistory) GetLast() PriceTimeStamp {
	if p.Len() == 0 {
		return EmptyPriceTimeStamp
	}

	return p[p.Len()-1]
}

func (p PriceHistory) String() string {
	if p.Len() == 0 {
		return EmptyPriceTimeStamp.String()
	}
	stringData := make([]string, 0, p.Len())
	for _, priceTimeStamp := range p {
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
