package bootstrap

import (
	"hris-backend/internal/config"
	"hris-backend/internal/infrastructure"
	"hris-backend/internal/middleware"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/auth"
	"hris-backend/internal/modules/company"
	"hris-backend/internal/modules/health"
	"hris-backend/internal/modules/leave"
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/notification"
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
	WebsocketHub *infrastructure.Hub
	Redis        *infrastructure.RedisClientProvider

	HealthCheckHandler   *health.Handler
	AuthHandler          *auth.Handler
	UserHandler          *user.Handler
	AttendanceHandler    *attendance.Handler
	MasterHandler        *master.Handler
	ReimbursementHandler *reimbursement.Handler
	PayrollHandler       *payroll.Handler
	LeaveHandler         *leave.Handler
	NotificationHandler  *notification.Handler
	CompanyHandler       *company.Handler

	AuthMiddleware        *middleware.AuthMiddleware
	RateLimiterMiddleware *middleware.RateLimiterMiddleware

	GeocodeWorker         attendance.GeocodeWorker
	LeaveScheduler        leave.Scheduler
	NotificationScheduler notification.Scheduler
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	db := infrastructure.NewGormConnection(&cfg.Database)
	storage := infrastructure.NewMinioStorage(&cfg.Minio)
	jwt := infrastructure.NewJWTProvider(&cfg.JWT)
	bcrypt := infrastructure.NewBcryptHasher(12)
	nominatim := infrastructure.NewNominatimFetcher(&cfg.ExternalServiceConfig)
	cronScheduler := infrastructure.NewCronProvider()
	redis := infrastructure.NewRedisClient(&cfg.Redis)
	transactionManager := infrastructure.NewGormTransactionManager(db.GetDB())

	wsHub := infrastructure.NewHub(redis.GetClient())
	geocodeWorker := attendance.NewGeocodeWorker(db.GetDB(), nominatim, 100)

	healthRepo := health.NewRepository(db.GetDB())
	userRepo := user.NewRepository(db.GetDB())
	attendanceRepo := attendance.NewRepository(db.GetDB())
	masterRepo := master.NewRepository(db.GetDB())
	reimburseRepo := reimbursement.NewRepository(db.GetDB())
	payrollRepo := payroll.NewRepository(db.GetDB())
	leaveRepo := leave.NewRepository(db.GetDB())
	notificationRepo := notification.NewRepository(db.GetDB())
	companyRepo := company.NewRepository(db.GetDB())

	healthSvc := health.NewService(healthRepo)
	notificationSvc := notification.NewService(wsHub, notificationRepo)
	authSvc := auth.NewService(userRepo, bcrypt, jwt)
	attendanceSvc := attendance.NewService(attendanceRepo, userRepo, storage, geocodeWorker)
	masterSvc := master.NewService(masterRepo)
	payrollSvc := payroll.NewService(payrollRepo, userRepo, reimburseRepo, attendanceRepo, companyRepo, notificationSvc, transactionManager)
	leaveSvc := leave.NewService(leaveRepo, storage, notificationSvc, userRepo)
	userSvc := user.NewService(userRepo, bcrypt, storage, leaveSvc)
	reimburseSvc := reimbursement.NewService(reimburseRepo, storage, notificationSvc, userRepo)
	companySvc := company.NewService(companyRepo, storage)

	healthHandler := health.NewHandler(healthSvc)
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler(userSvc)
	attendanceHandler := attendance.NewHandler(attendanceSvc)
	masterHandler := master.NewHandler(masterSvc)
	reimburseHandler := reimbursement.NewHandler(reimburseSvc)
	payrollHandler := payroll.NewHandler(payrollSvc)
	leaveHandler := leave.NewHandler(leaveSvc)
	notificationHandler := notification.NewHandler(wsHub, notificationSvc)
	companyHandler := company.NewHandler(companySvc)

	authMiddleware := middleware.NewAuthMiddleware(jwt)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware()

	leaveScheduler := leave.NewScheduler(cronScheduler, leaveSvc)
	notificationScheduler := notification.NewScheduler(cronScheduler, notificationSvc)

	return &Container{
		Config:       cfg,
		DB:           db,
		Storage:      storage,
		JWT:          jwt,
		Bcrypt:       bcrypt,
		Location:     nominatim,
		WebsocketHub: wsHub,
		Redis:        redis,

		HealthCheckHandler:   healthHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		AttendanceHandler:    attendanceHandler,
		MasterHandler:        masterHandler,
		ReimbursementHandler: reimburseHandler,
		PayrollHandler:       payrollHandler,
		LeaveHandler:         leaveHandler,
		NotificationHandler:  notificationHandler,
		CompanyHandler:       companyHandler,

		AuthMiddleware:        authMiddleware,
		RateLimiterMiddleware: rateLimiterMiddleware,

		GeocodeWorker:         geocodeWorker,
		LeaveScheduler:        leaveScheduler,
		NotificationScheduler: notificationScheduler,
	}, nil
}

// Close properly closes all resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	if c.GeocodeWorker != nil {
		c.GeocodeWorker.Stop()
	}

	if c.LeaveScheduler != nil {
		c.LeaveScheduler.Stop()
	}

	if c.NotificationScheduler != nil {
		c.NotificationScheduler.Stop()
	}

	if c.Redis != nil {
		c.Redis.Close()
	}

	return nil
}
