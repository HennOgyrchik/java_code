package service

import (
	"github.com/gin-gonic/gin"
	"java_code/pkg/db"
)

type Service struct {
	db db.DB
}

func New(db db.DB) Service {
	return Service{db: db}
}

func (s *Service) Test(c *gin.Context) {

}
