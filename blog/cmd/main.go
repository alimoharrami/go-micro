package main

import (
	"context"
	"go-blog/internal/config"
	"go-blog/internal/database"
	"go-blog/internal/handlers"
	"go-blog/internal/helpers"
	"go-blog/internal/repository"
	"go-blog/internal/routes"
	"go-blog/internal/server"
	"go-blog/internal/service"
	"go-blog/migrations"
	"log"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(
			// Config
			config.LoadConfig,

			// Database
			database.NewPostgresConfig,
			database.NewPostgresConnection,

			// GRPC Client
			helpers.InitGRPC,

			// Repositories
			repository.NewPostRepository,

			// Services
			service.NewPostService,

			// Handlers
			handlers.NewPostController,

			// Router
			routes.NewRouter,

			// Server
			server.NewServer,
		),
		fx.Invoke(
			RegisterHooks,
			RunMigrations,
		),
	).Run()
}

func RunMigrations(db *gorm.DB) {
	migrations.AutoMigrate(db)
}

func RegisterHooks(lc fx.Lifecycle, srv *server.Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Start(cfg.Server.Port); err != nil {
					log.Printf("Server failed to start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
