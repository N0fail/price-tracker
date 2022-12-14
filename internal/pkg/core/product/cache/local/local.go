package local

import (
	"context"
	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"sort"
	"sync"
)

const poolSize = 10

func New() cachePkg.Interface {
	return &cache{
		muP:          sync.RWMutex{},
		product:      make(map[string]models.Product, 0),
		muH:          sync.RWMutex{},
		priceHistory: make(map[string]models.PriceHistory, 0),
		poolCh:       make(chan struct{}, poolSize),
	}
}

type cache struct {
	muP          sync.RWMutex
	product      map[string]models.Product
	muH          sync.RWMutex
	priceHistory map[string]models.PriceHistory
	poolCh       chan struct{}
}

func (c *cache) ProductList(ctx context.Context, pageNumber, resultsPerPage uint32, orderBy string) (models.ProductSnapshots, error) {
	c.poolCh <- struct{}{}
	defer func() {
		<-c.poolCh
	}()
	c.muP.RLock()
	defer c.muP.RUnlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	from := pageNumber * resultsPerPage
	to := (pageNumber + 1) * resultsPerPage
	if from >= uint32(len(c.product)) {
		return nil, error_codes.ErrNoEntries
	}
	allProducts := make([]models.Product, 0, len(c.product))

	for _, product := range c.product {
		allProducts = append(allProducts, product)
	}
	sort.SliceStable(allProducts, func(i, j int) bool {
		if orderBy == "name" {
			return allProducts[i].Name < allProducts[j].Name
		} else {
			return allProducts[i].Code < allProducts[j].Code
		}
	})

	res := make(models.ProductSnapshots, 0, resultsPerPage)
	for i := from; i < to; i++ {
		if i >= uint32(len(allProducts)) {
			break
		}
		product := allProducts[i]
		res = append(res, models.ProductSnapshot{
			Name:      product.Name,
			Code:      product.Code,
			LastPrice: c.priceHistory[product.Code].GetLast(),
		})
	}

	return res, nil
}

func (c *cache) ProductCreate(ctx context.Context, p models.Product) error {
	c.poolCh <- struct{}{}
	defer func() {
		<-c.poolCh
	}()
	c.muP.Lock()
	defer c.muP.Unlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.product[p.Code]; ok {
		return errors.Wrap(error_codes.ErrProductExists, p.Code)
	}
	c.product[p.Code] = p
	c.priceHistory[p.Code] = make(models.PriceHistory, 0)
	return nil
}

func (c *cache) ProductDelete(ctx context.Context, code string) error {
	c.poolCh <- struct{}{}
	defer func() {
		<-c.poolCh
	}()
	c.muP.Lock()
	defer c.muP.Unlock()
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.product[code]; !ok {
		return errors.Wrap(error_codes.ErrProductNotExist, code)
	}
	delete(c.product, code)
	delete(c.priceHistory, code)
	return nil
}

func (c *cache) PriceTimeStampAdd(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error {
	c.poolCh <- struct{}{}
	defer func() {
		<-c.poolCh
	}()
	c.muH.Lock()
	defer c.muH.Unlock()

	if _, ok := c.priceHistory[code]; !ok {
		return errors.Wrap(error_codes.ErrProductNotExist, code)
	}

	priceHistory := c.priceHistory[code]
	priceHistory = append(priceHistory, priceTimeStamp)
	sort.Stable(priceHistory)
	c.priceHistory[code] = priceHistory
	return nil
}

func (c *cache) PriceHistory(ctx context.Context, code string) (models.PriceHistory, error) {
	c.poolCh <- struct{}{}
	defer func() {
		<-c.poolCh
	}()
	c.muH.RLock()
	defer c.muH.RUnlock()

	if _, ok := c.priceHistory[code]; !ok {
		return nil, errors.Wrap(error_codes.ErrProductNotExist, code)
	}

	return c.priceHistory[code].Copy(), nil
}
