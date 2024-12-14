package db

import (
	"context"
	"github.com/gofrs/uuid"
)

type DB interface {
	Start(ctx context.Context) error
	Stop()

	Update(context.Context, Wallets) error
	GetBalance(context.Context, uuid.UUID) (float64, error)
	Create(context.Context, uuid.UUID) error
}

type Wallets struct {
	Id      uuid.UUID
	Balance float64
}
