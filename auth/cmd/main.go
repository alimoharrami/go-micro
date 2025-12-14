package main

import (
	"auth/internal/database"
	"auth/internal/handlers"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/migrations"
	"context"
	"log"

	"auth/internal/config"
	"auth/internal/routes"

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

			// Repositories
			repository.NewUserRepository,
			repository.NewRoleRepository,
			repository.NewPermissionRepository,
			repository.NewRolePermissionRepository,

			// Services
			service.NewAuthService,
			service.NewRoleService,
			service.NewPermissionService,
			service.NewRolePermissionService,

			// Handlers
			handlers.NewAuthController,

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
