package psql

import (
	"context"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"java_code/pkg/db"
	"time"
)

type PSQL struct {
	timeout time.Duration
	url     string
	db      *gorm.DB
}

func New(url string, timeout time.Duration) PSQL {
	return PSQL{
		timeout: timeout,
		url:     url,
		db:      nil,
	}
}

func (p *PSQL) Start() error {
	var err error

	if p.db, err = gorm.Open(postgres.Open(p.url), &gorm.Config{}); err != nil {
		return err
	}

	q, err := p.db.DB()
	if err != nil {
		return err
	}

	q.SetMaxIdleConns(1000)
	q.SetMaxOpenConns(1000)

	err = p.db.AutoMigrate(&db.Wallets{})
	return err
}

func (p *PSQL) Stop() {
	if p.db != nil {
		dbInstance, _ := p.db.DB()
		_ = dbInstance.Close()
	}
}

func (p *PSQL) Update(ctx context.Context, wal db.Wallets) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.db.WithContext(ctxWithTimeout).Save(wal).Error

}

func (p *PSQL) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	data := db.Wallets{}

	err := p.db.WithContext(ctxWithTimeout).Select("balance").Where("id = ?", id).Find(&data).Error

	return data.Balance, err
}
