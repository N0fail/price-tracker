package api

import (
	"context"
	"github.com/golang/mock/gomock"
	mockProductApi "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/api/mocks"
	"testing"
)

type BaseFixture struct {
	Ctx context.Context
}

type apiFixture struct {
	Ctx context.Context
	api *implementation
}

type productFixture struct {
	Ctx     context.Context
	product *mockProductApi.MockInterface
	api     *implementation
}

func productSetUp(t *testing.T) productFixture {
	t.Parallel()

	p := productFixture{}
	p.Ctx = context.Background()
	p.product = mockProductApi.NewMockInterface(gomock.NewController(t))
	p.api = New(p.product)

	return p
}
