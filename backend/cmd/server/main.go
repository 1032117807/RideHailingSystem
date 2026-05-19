package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ridehailing/backend/internal/config"
	"ridehailing/backend/internal/handler"
	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/pkg/jwtutil"
	"ridehailing/backend/internal/pkg/mailer"
	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/repository"
	"ridehailing/backend/internal/router"
	"ridehailing/backend/internal/service"
)

func main() {
	cfg := config.Load()

	gormLogger := logger.New(
		log.New(os.Stdout, "[gorm] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Trip{},
		&model.TripStop{},
		&model.Order{},
		&model.Payment{},
		&model.Notification{},
		&model.TokenUsage{},
		&model.RiskEvent{},
		&model.ElectronicTicket{},
		&model.TicketVerification{},
		&model.RefundAuditLog{},
		&model.AuditLog{},
		&model.Passenger{},
		&model.PriceAlert{},
		&model.DriverProfile{},
		&model.Vehicle{},
		&model.DriverSettlement{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	userRepo := repository.NewGormUserRepository(db)
	tripRepo := repository.NewGormTripRepository(db)
	orderRepo := repository.NewGormOrderRepository(db)
	paymentRepo := repository.NewGormPaymentRepository(db)
	notificationRepo := repository.NewGormNotificationRepository(db)
	electronicTicketRepo := repository.NewGormElectronicTicketRepository(db)
	auditRepo := repository.NewGormAuditRepository(db)
	passengerRepo := repository.NewGormPassengerRepository(db)
	priceAlertRepo := repository.NewGormPriceAlertRepository(db)
	driverProfileRepo := repository.NewGormDriverProfileRepository(db)
	codeRepo := repository.NewRedisCodeRepository(redisClient, cfg.RedisKeyPrefix)
	aiLimiter := repository.NewRedisAIRateLimiter(redisClient, cfg.AIRateLimitPrefix)
	tokenUsageRepo := repository.NewTokenUsageRepository(db)
	riskEventRepo := repository.NewRiskEventRepository(db)

	jwtManager := jwtutil.NewManager(cfg.JWTSecret, cfg.JWTExpireTime)
	devMailer := mailer.NewLogMailer()

	authService := service.NewAuthService(userRepo, codeRepo, jwtManager, cfg.CodeTTL, devMailer)
	userService := service.NewUserService(userRepo)
	ticketService := service.NewTicketService(tripRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	electronicTicketService := service.NewElectronicTicketService(electronicTicketRepo, orderRepo, cfg.JWTSecret)
	orderService := service.NewOrderService(orderRepo, tripRepo, userRepo, paymentRepo, notificationRepo, auditRepo, electronicTicketService)
	tripService := service.NewTripService(tripRepo, orderRepo)
	auditService := service.NewAuditService(auditRepo)
	passengerService := service.NewPassengerService(passengerRepo)
	priceAlertService := service.NewPriceAlertService(priceAlertRepo)
	driverProfileService := service.NewDriverProfileService(driverProfileRepo)
	tokenUsageService := service.NewTokenUsageService(tokenUsageRepo)
	riskService := service.NewRiskService(riskEventRepo, tokenUsageRepo)
	knowledgeService := service.NewKnowledgeService(
		cfg.KnowledgeMemoryDir,
		cfg.EmbeddingAPIKey,
		cfg.EmbeddingBaseURL,
		cfg.EmbeddingModel,
		cfg.RerankAPIKey,
		cfg.RerankBaseURL,
		cfg.RerankModel,
		tokenUsageService,
		cfg.KnowledgeChunkSize,
		cfg.KnowledgeChunkOverlap,
		cfg.KnowledgeRecallLimit,
		cfg.KnowledgeTopK,
	)
	aiService := service.NewAIService(
		cfg.OpenAIAPIKey,
		cfg.OpenAIBaseURL,
		cfg.OpenAIModel,
		cfg.OpenAITimeout,
		ticketService,
		orderService,
		userService,
		knowledgeService,
		tokenUsageService,
	)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo, electronicTicketService)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	driverTripHandler := handler.NewDriverTripHandler(tripService)
	tokenUsageHandler := handler.NewTokenUsageHandler(tokenUsageService)
	aiHandler := handler.NewAIHandler(
		aiService,
		aiLimiter,
		cfg.AIPassengerChatLimit,
		cfg.AIDriverDraftLimit,
		cfg.AIRateLimitWindow,
		riskService,
	)
	ticketHandler := handler.NewTicketHandler(ticketService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	electronicTicketHandler := handler.NewElectronicTicketHandler(electronicTicketService)
	auditHandler := handler.NewAuditHandler(auditService)
	passengerHandler := handler.NewPassengerHandler(passengerService)
	priceAlertHandler := handler.NewPriceAlertHandler(priceAlertService)
	driverProfileHandler := handler.NewDriverProfileHandler(driverProfileService)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeService)
	riskHandler := handler.NewRiskHandler(riskService)

	appRouter := router.NewRouter(router.Dependencies{
		AuthHandler:             authHandler,
		UserHandler:             userHandler,
		DriverTripHandler:       driverTripHandler,
		AIHandler:               aiHandler,
		KnowledgeHandler:        knowledgeHandler,
		RiskHandler:             riskHandler,
		TicketHandler:           ticketHandler,
		OrderHandler:            orderHandler,
		PaymentHandler:          paymentHandler,
		NotificationHandler:     notificationHandler,
		ElectronicTicketHandler: electronicTicketHandler,
		AuditHandler:            auditHandler,
		PassengerHandler:        passengerHandler,
		PriceAlertHandler:       priceAlertHandler,
		DriverProfileHandler:    driverProfileHandler,
		AuthMiddleware:          middleware.AuthMiddleware(jwtManager),
		EnableMockPayment:       cfg.EnableMockPayment,
		TokenUsageHandler:       tokenUsageHandler,
	})

	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      appRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			closedCount, err := orderService.CloseExpiredPendingOrders(context.Background())
			if err != nil {
				log.Printf("close expired pending orders failed: %v", err)
				continue
			}
			if closedCount > 0 {
				log.Printf("closed %d expired pending orders", closedCount)
			}
		}
	}()

	log.Printf("server listening on %s", cfg.ServerAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server start failed: %v", err)
	}
}
