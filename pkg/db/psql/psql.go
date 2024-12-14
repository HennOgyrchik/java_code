package psql

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"java_code/pkg/db"
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

	pool, err := pgxpool.Connect(ctxTimeout, p.url)
	if err != nil {
		return fmt.Errorf("Create PSQL connect: %w", err)
	}

	p.pool = pool

	_, err = p.pool.Exec(ctx, "create table if not exists  wallets (id uuid unique, balance float)")

	return err
}

func (p *PSQL) Stop() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *PSQL) Update(ctx context.Context, wal db.Wallets) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	x, err := p.pool.Exec(ctxWithTimeout, "update wallets set balance = balance + $1 where id = $2", wal.Balance, wal.Id)

	if x.RowsAffected() < 1 {
		err = pgx.ErrNoRows
	}

	return err

}

func (p *PSQL) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	row := p.pool.QueryRow(ctxWithTimeout, "select balance from wallets where id = $1", id)

	var balance float64
	err := row.Scan(&balance)

	return balance, err
}

func (p *PSQL) Create(ctx context.Context, id uuid.UUID) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctxWithTimeout, "insert into wallets(id,balance) values ($1,$2)", id, 0.0)

	return err
}
