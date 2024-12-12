package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
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

	balance, err := s.db.Balance(s.ctx, req.WalletId)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err)
		return
	}

	req.OperationType = strings.ToLower(req.OperationType)

	switch req.OperationType {
	case "deposit":
		balance += req.Amount
	case "withdraw":
		if balance < req.Amount {
			sendError(c, http.StatusBadRequest, fmt.Errorf("insufficient funds"))
			return
		}
		balance -= req.Amount
	default:
		sendError(c, http.StatusBadRequest, fmt.Errorf("unknown type of operation"))
		return
	}

	wal := db.Wallets{
		Id:      req.WalletId,
		Balance: balance,
	}
	if err := s.db.Update(s.ctx, wal); err != nil {
		sendError(c, http.StatusInternalServerError, err)
	}
}

func (s *Service) Wallets(c *gin.Context) {
	id, err := uuid.FromString(c.Param("uuid"))
	if err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	balance, err := s.db.Balance(s.ctx, id)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, db.Wallets{
		Id:      id,
		Balance: balance,
	})
}

func sendError(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, struct {
		Error string
	}{Error: err.Error()})
}
