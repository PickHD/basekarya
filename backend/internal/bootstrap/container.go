package bootstrap

import (
	"basekarya-backend/internal/config"
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/middleware"
	"basekarya-backend/internal/modules/attendance"
	"basekarya-backend/internal/modules/auth"
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/health"
	"basekarya-backend/internal/modules/leave"
	"basekarya-backend/internal/modules/loan"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/notification"
	"basekarya-backend/internal/modules/overtime"
	"basekarya-backend/internal/modules/payroll"
	"basekarya-backend/internal/modules/reimbursement"
	"basekarya-backend/internal/modules/user"
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
	Email        *infrastructure.EmailProvider
	Excel        infrastructure.ExcelProvider

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
	LoanHandler          *loan.Handler
	OvertimeHandler      *overtime.Handler

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
	cronScheduler := infrastructure.NewCronProvider()
	redis := infrastructure.NewRedisClient(&cfg.Redis)
	transactionManager := infrastructure.NewGormTransactionManager(db.GetDB())
	httpClient := infrastructure.NewHttpClientProvider()
	nominatim := infrastructure.NewNominatimFetcher(&cfg.ExternalServiceConfig, httpClient.GetClient())
	email := infrastructure.NewEmailProvider(&cfg.Email)
	excel := infrastructure.NewExcelProvider()

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
	loanRepo := loan.NewRepository(db.GetDB())
	overtimeRepo := overtime.NewRepository(db.GetDB())

	healthSvc := health.NewService(healthRepo)
	notificationSvc := notification.NewService(wsHub, notificationRepo)
	authSvc := auth.NewService(userRepo, bcrypt, jwt)
	attendanceSvc := attendance.NewService(attendanceRepo, userRepo, storage, geocodeWorker, transactionManager, excel)
	masterSvc := master.NewService(masterRepo)
	payrollSvc := payroll.NewService(payrollRepo, userRepo, reimburseRepo, attendanceRepo, companyRepo, notificationSvc, transactionManager, httpClient.GetClient(), email, loanRepo, overtimeRepo)
	leaveSvc := leave.NewService(leaveRepo, storage, notificationSvc, userRepo, transactionManager, excel)
	userSvc := user.NewService(userRepo, bcrypt, storage, leaveSvc, transactionManager)
	reimburseSvc := reimbursement.NewService(reimburseRepo, storage, notificationSvc, userRepo, transactionManager, excel)
	companySvc := company.NewService(companyRepo, storage)
	loanSvc := loan.NewService(loanRepo, notificationSvc, userRepo, transactionManager, excel)
	overtimeSvc := overtime.NewService(overtimeRepo, notificationSvc, userRepo, transactionManager, excel)

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
	loanHandler := loan.NewHandler(loanSvc)
	overtimeHandler := overtime.NewHandler(overtimeSvc)

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
		Email:        email,
		Excel:        excel,

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
		LoanHandler:          loanHandler,
		OvertimeHandler:      overtimeHandler,

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
