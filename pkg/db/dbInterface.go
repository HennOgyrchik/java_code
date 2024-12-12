package db

import (
	"context"
	"github.com/gofrs/uuid"
)

type DB interface {
	Start() error
	Stop()

	Update(context.Context, Wallets) error
	Balance(context.Context, uuid.UUID) (float64, error)
}

type Wallets struct {
	Id      uuid.UUID
	Balance float64
}
