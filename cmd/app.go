package cmd

import (
	"context"
	"fmt"
	"github.com/elnoro/foxylock/m/v2/admin"
	"github.com/elnoro/foxylock/m/v2/config"
	"github.com/elnoro/foxylock/m/v2/db"
	coredns_integration "github.com/elnoro/foxylock/m/v2/dns/coredns-integration"
	"go.uber.org/zap"
)

type App struct {
	config *config.Config
	log    *zap.SugaredLogger
}

func NewProduction() (*App, error) {
	c, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config, %w", err)
	}

	p, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("getting logger, %w", err)
	}

	return newApp(c, p.Sugar()), nil
}

func newApp(config *config.Config, log *zap.SugaredLogger) *App {
	return &App{config: config, log: log}
}

func (a App) Start(ctx context.Context) error {
	a.log.Infof("Loading file: %s", a.config.HostsFile)
	inMemoryDb, err := db.NewInMemoryFromFile(a.config.HostsFile)
	if err != nil {
		return fmt.Errorf("db init, %w", inMemoryDb)
	}

	if a.config.RedisPass != "" {
		a.log.Infof("Starting Redis at %s", a.config.RedisAddr)
		rs := admin.NewRedisLikeServer(inMemoryDb, a.config.RedisAddr, a.config.RedisPass)
		a.startServer(rs, ctx)
	}

	if a.config.HttpAddr != "" {
		a.log.Infof("Starting Http at %s", a.config.HttpAddr)
		gs := admin.NewHttpServer(inMemoryDb, a.config.HttpAddr)
		a.startServer(gs, ctx)
	}

	err = coredns_integration.NewCoreDns(inMemoryDb).Start()
	if err != nil {
		return fmt.Errorf("starting CoreDNS, %w", err)
	}

	a.log.Info("Running...")

	return nil
}

func (a App) startServer(s admin.DbServer, ctx context.Context) {
	go func(s admin.DbServer, l *zap.SugaredLogger) {
		err := s.Run(ctx)
		l.Errorf("failed to stop a server, %v", err)
	}(s, a.log)
}
