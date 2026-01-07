package bootstrap

import (
	"hris-backend/internal/config"
	"hris-backend/internal/infrastructure"
	"hris-backend/internal/middleware"
	"hris-backend/internal/modules/auth"
	"hris-backend/internal/modules/health"
	"hris-backend/internal/modules/user"

	"golang.org/x/crypto/bcrypt"
)

type Container struct {
	Config  *config.Config
	DB      *infrastructure.GormConnectionProvider
	Storage *infrastructure.MinioStorageProvider
	JWT     *infrastructure.JwtProvider
	Bcrypt  *infrastructure.BcryptHasher

	HealthCheckHandler *health.Handler
	AuthHandler        *auth.Handler
	UserHandler        *user.Handler

	AuthMiddleware *middleware.AuthMiddleware
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	db := infrastructure.NewGormConnection(cfg)
	storage := infrastructure.NewMinioStorage(cfg)
	jwt := infrastructure.NewJWTProvider(cfg)
	bcrypt := infrastructure.NewBcryptHasher(bcrypt.DefaultCost)

	healthRepo := health.NewRepository(db.GetDB())
	userRepo := user.NewRepository(db.GetDB())

	healthSvc := health.NewService(healthRepo)
	authSvc := auth.NewService(userRepo, bcrypt, jwt)
	userSvc := user.NewService(userRepo, bcrypt, storage)

	healthHandler := health.NewHandler(healthSvc)
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler(userSvc)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	return &Container{
		Config:  cfg,
		DB:      db,
		Storage: storage,
		JWT:     jwt,
		Bcrypt:  bcrypt,

		HealthCheckHandler: healthHandler,
		AuthHandler:        authHandler,
		UserHandler:        userHandler,

		AuthMiddleware: authMiddleware,
	}, nil
}

// Close properly closes all resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
