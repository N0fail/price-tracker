package postgres

import (
	"context"
	"github.com/pashagolub/pgxmock"
	databasePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/database"
	"log"
	"testing"
)

type postgresFixture struct {
	postgres   databasePkg.Interface
	dbConnMock pgxmock.PgxConnIface
}

func setUp(t *testing.T) postgresFixture {
	var fixture postgresFixture

	mock, err := pgxmock.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	fixture.dbConnMock = mock
	fixture.postgres = New(mock)
	return fixture
}

func (f *postgresFixture) tearDown(ctx context.Context) {
	f.dbConnMock.Close(ctx)
}
