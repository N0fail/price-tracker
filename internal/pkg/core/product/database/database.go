package database

import (
	productApiPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/api"
)

type Interface interface {
	productApiPkg.Interface
}
