package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Gin struct {
	srv *http.Server
}

func New(url string, handler Handler) Gin {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// TODO add endpoints
	router.GET("/test/:id", handler.Test)

	return Gin{srv: &http.Server{Addr: url, Handler: router.Handler()}}
}

func (g *Gin) Start() error {

	if err := g.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (g *Gin) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	if err = g.srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("web server was shutdown incorrectly: %w", err)
	}

	return err
}
