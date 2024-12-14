package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"java_code/pkg/db"
	"net/http"
	"strings"
)

type Service struct {
	ctx context.Context
	db  db.DB
}

type Request struct {
	WalletId      uuid.UUID
	OperationType string
	Amount        float64
}

func New(ctx context.Context, db db.DB) Service {
	return Service{ctx: ctx, db: db}
}

func (s *Service) Wallet(c *gin.Context) {
	var req Request

	if err := c.BindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	wal := db.Wallets{
		Id:      req.WalletId,
		Balance: 0.0,
	}

	req.OperationType = strings.ToLower(req.OperationType)

	switch req.OperationType {
	case "deposit":
		wal.Balance = req.Amount
	case "withdraw":
		balance, err := s.db.GetBalance(s.ctx, req.WalletId)
		if err != nil {
			sendError(c, http.StatusInternalServerError, err)
			return
		}
		if balance < req.Amount {
			sendError(c, http.StatusBadRequest, fmt.Errorf("insufficient funds"))
			return
		}
		wal.Balance = -(req.Amount)
	default:
		sendError(c, http.StatusBadRequest, fmt.Errorf("unknown type of operation"))
		return
	}

	err := s.db.Update(s.ctx, wal)

	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		sendError(c, http.StatusBadRequest, err)
	default:
		sendError(c, http.StatusInternalServerError, err)
	}

}

func (s *Service) Wallets(c *gin.Context) {
	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	balance, err := s.db.GetBalance(s.ctx, id)

	switch {
	case err == nil:
		c.JSON(http.StatusOK, db.Wallets{
			Id:      id,
			Balance: balance,
		})
	case errors.Is(err, pgx.ErrNoRows):
		sendError(c, http.StatusBadRequest, err)
	default:
		sendError(c, http.StatusInternalServerError, err)
	}

}

func (s *Service) NewWallet(c *gin.Context) {
	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	err = s.db.Create(s.ctx, id)

	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		sendError(c, http.StatusBadRequest, err)
	default:
		sendError(c, http.StatusInternalServerError, err)
	}
}

func sendError(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, struct {
		Error string
	}{Error: err.Error()})
}
