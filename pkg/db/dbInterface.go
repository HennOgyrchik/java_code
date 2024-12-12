package db

import "context"

type DB interface {
	Start(context.Context) error
	Stop()

	Test(context.Context)
}
