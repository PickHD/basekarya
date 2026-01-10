package bootstrap

import (
	"hris-backend/internal/config"
	"hris-backend/internal/infrastructure"
	"hris-backend/internal/middleware"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/auth"
	"hris-backend/internal/modules/health"
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/user"
)

type Container struct {
	Config   *config.Config
	DB       *infrastructure.GormConnectionProvider
	Storage  *infrastructure.MinioStorageProvider
	JWT      *infrastructure.JwtProvider
	Bcrypt   *infrastructure.BcryptHasher
	Location *infrastructure.NominatimFetcher

	HealthCheckHandler *health.Handler
	AuthHandler        *auth.Handler
	UserHandler        *user.Handler
	AttendanceHandler  *attendance.Handler
	MasterHandler      *master.Handler

	AuthMiddleware *middleware.AuthMiddleware
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	db := infrastructure.NewGormConnection(cfg)
	storage := infrastructure.NewMinioStorage(cfg)
	jwt := infrastructure.NewJWTProvider(cfg)
	bcrypt := infrastructure.NewBcryptHasher(12)
	nominatim := infrastructure.NewNominatimFetcher(cfg)

	healthRepo := health.NewRepository(db.GetDB())
	userRepo := user.NewRepository(db.GetDB())
	attendanceRepo := attendance.NewRepository(db.GetDB())
	masterRepo := master.NewRepository(db.GetDB())

	healthSvc := health.NewService(healthRepo)
	authSvc := auth.NewService(userRepo, bcrypt, jwt)
	userSvc := user.NewService(userRepo, bcrypt, storage)
	attendanceSvc := attendance.NewService(attendanceRepo, userRepo, storage)
	masterSvc := master.NewService(masterRepo)

	healthHandler := health.NewHandler(healthSvc)
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler(userSvc)
	attendanceHandler := attendance.NewHandler(attendanceSvc)
	masterHandler := master.NewHandler(masterSvc)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	return &Container{
		Config:   cfg,
		DB:       db,
		Storage:  storage,
		JWT:      jwt,
		Bcrypt:   bcrypt,
		Location: nominatim,

		HealthCheckHandler: healthHandler,
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		AttendanceHandler:  attendanceHandler,
		MasterHandler:      masterHandler,

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
