package local

import (
	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"sort"
	"sync"
)

var (
	ErrProductNotExist = errors.New("product does not exist")
	ErrProductExists   = errors.New("product exists")
)

func New() cachePkg.Interface {
	return &cache{
		muP:          sync.RWMutex{},
		product:      make(map[string]models.Product, 0),
		muH:          sync.RWMutex{},
		priceHistory: make(map[string]models.PriceHistory, 0),
	}
}

type cache struct {
	muP          sync.RWMutex
	product      map[string]models.Product
	muH          sync.RWMutex
	priceHistory map[string]models.PriceHistory
}

func (c *cache) ProductList() []models.ProductSnapshot {
	c.muP.RLock()
	defer c.muP.RUnlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	res := make([]models.ProductSnapshot, 0, len(c.product))
	for _, product := range c.product {
		res = append(res, models.ProductSnapshot{
			Name:      product.Name,
			Code:      product.Code,
			LastPrice: c.priceHistory[product.Code].GetLast(),
		})
	}
	return res
}

func (c *cache) ProductCreate(p models.Product) error {
	c.muP.Lock()
	defer c.muP.Unlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.product[p.Code]; ok {
		return errors.Wrap(ErrProductExists, p.Code)
	}
	c.product[p.Code] = p
	c.priceHistory[p.Code] = make(models.PriceHistory, 0)
	return nil
}

func (c *cache) ProductDelete(code string) error {
	c.muP.Lock()
	defer c.muP.Unlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.product[code]; !ok {
		return errors.Wrap(ErrProductNotExist, code)
	}
	delete(c.product, code)
	delete(c.priceHistory, code)
	return nil
}

func (c *cache) AddPriceTimeStamp(code string, priceTimeStamp models.PriceTimeStamp) error {
	c.muH.Lock()
	defer c.muH.Unlock()

	if _, ok := c.priceHistory[code]; !ok {
		return errors.Wrap(ErrProductNotExist, code)
	}

	priceHistory := c.priceHistory[code]
	priceHistory = append(priceHistory, priceTimeStamp)
	sort.Stable(priceHistory)
	c.priceHistory[code] = priceHistory
	return nil
}

func (c *cache) FullHistory(code string) (models.PriceHistory, error) {
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.priceHistory[code]; !ok {
		return nil, errors.Wrap(ErrProductNotExist, code)
	}

	return c.priceHistory[code].Copy(), nil
}
