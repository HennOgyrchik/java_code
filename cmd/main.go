package main

import (
	"context"
	"java_code/pkg/config"
	"java_code/pkg/db/psql"
	"java_code/pkg/service"
	"java_code/pkg/web/ginSrvr"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	if err := config.LoadConfig("./app/config/config.env"); err != nil {
		slog.Error("read configuration", "error", err.Error())
		return
	}

	cfg := config.New()

	dbUrl, err := cfg.Postgres.ConnectionURL()

	if err != nil {
		slog.Error("read db url", "error", err.Error())
		return
	}

	db := psql.New(dbUrl, time.Duration(cfg.Postgres.ConnTimeout)*time.Second)

	if err := db.Start(ctx); err != nil {
		slog.Error("connection db", "error", err.Error())
		return
	}
	defer db.Stop()

	app := service.New(ctx, &db)

	webSrv := ginSrvr.New(cfg.Web.ConnectionURL(), &app)

	go func() {
		<-ctx.Done()
		if err = webSrv.Stop(); err != nil {
			slog.Error("closing web server", "error", err.Error())
			return
		}
		slog.Info("web server is closed")
	}()

	err = webSrv.Start()
	if err != nil {
		slog.Error("web server", "error", err.Error())
		return
	}
}
