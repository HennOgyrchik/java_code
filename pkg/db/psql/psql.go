package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PSQL struct {
	timeout time.Duration
	url     string
	pool    *pgxpool.Pool
}

func New(url string, timeout time.Duration) PSQL {
	return PSQL{
		timeout: timeout,
		url:     url,
		pool:    nil,
	}
}

func (p *PSQL) Start(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var err error

	p.pool, err = pgxpool.Connect(ctxTimeout, p.url)
	return err
}

func (p *PSQL) Stop() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *PSQL) RowsFunc(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	rows, err := p.pool.Query(ctxTimeout, "select * from test_table")
	if err != nil {
		return err
	}

	for rows.Next() {
		var x int

		if err = rows.Scan(&x); err != nil {
			return err
		}
	}

	return nil
}

func (p *PSQL) ExecFunc(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctxTimeout, "insert into test_table(id) values ($1)", 51)

	return err
}

func (p *PSQL) Test(ctx context.Context) {
	fmt.Println(" PSQL TEST FUNCTION")
}
