package bootstrap

import (
	"hris-backend/internal/config"
	"hris-backend/internal/infrastructure"
	"hris-backend/internal/middleware"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/auth"
	"hris-backend/internal/modules/health"
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/payroll"
	"hris-backend/internal/modules/reimbursement"
	"hris-backend/internal/modules/user"
)

type Container struct {
	Config       *config.Config
	DB           *infrastructure.GormConnectionProvider
	Storage      *infrastructure.MinioStorageProvider
	JWT          *infrastructure.JwtProvider
	Bcrypt       *infrastructure.BcryptHasher
	Location     *infrastructure.NominatimFetcher
	GeocodeQueue <-chan attendance.GeocodeJob

	HealthCheckHandler   *health.Handler
	AuthHandler          *auth.Handler
	UserHandler          *user.Handler
	AttendanceHandler    *attendance.Handler
	MasterHandler        *master.Handler
	ReimbursementHandler *reimbursement.Handler
	PayrollHandler       *payroll.Handler

	AuthMiddleware        *middleware.AuthMiddleware
	RateLimiterMiddleware *middleware.RateLimiterMiddleware
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	db := infrastructure.NewGormConnection(cfg)
	storage := infrastructure.NewMinioStorage(cfg)
	jwt := infrastructure.NewJWTProvider(cfg)
	bcrypt := infrastructure.NewBcryptHasher(12)
	nominatim := infrastructure.NewNominatimFetcher(cfg)
	geocodeQueue := make(chan attendance.GeocodeJob, 100)

	healthRepo := health.NewRepository(db.GetDB())
	userRepo := user.NewRepository(db.GetDB())
	attendanceRepo := attendance.NewRepository(db.GetDB())
	masterRepo := master.NewRepository(db.GetDB())
	reimburseRepo := reimbursement.NewRepository(db.GetDB())
	payrollRepo := payroll.NewRepository(db.GetDB())

	healthSvc := health.NewService(healthRepo)
	authSvc := auth.NewService(userRepo, bcrypt, jwt)
	userSvc := user.NewService(userRepo, bcrypt, storage)
	attendanceSvc := attendance.NewService(attendanceRepo, userRepo, storage, geocodeQueue)
	masterSvc := master.NewService(masterRepo)
	reimburseSvc := reimbursement.NewService(reimburseRepo, storage)
	payrollSvc := payroll.NewService(payrollRepo, userRepo, reimburseRepo, attendanceRepo)

	healthHandler := health.NewHandler(healthSvc)
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler(userSvc)
	attendanceHandler := attendance.NewHandler(attendanceSvc)
	masterHandler := master.NewHandler(masterSvc)
	reimburseHandler := reimbursement.NewHandler(reimburseSvc)
	payrollHandler := payroll.NewHandler(payrollSvc)

	authMiddleware := middleware.NewAuthMiddleware(jwt)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware()

	return &Container{
		Config:       cfg,
		DB:           db,
		Storage:      storage,
		JWT:          jwt,
		Bcrypt:       bcrypt,
		Location:     nominatim,
		GeocodeQueue: geocodeQueue,

		HealthCheckHandler:   healthHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		AttendanceHandler:    attendanceHandler,
		MasterHandler:        masterHandler,
		ReimbursementHandler: reimburseHandler,
		PayrollHandler:       payrollHandler,

		AuthMiddleware:        authMiddleware,
		RateLimiterMiddleware: rateLimiterMiddleware,
	}, nil
}

// Close properly closes all resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
