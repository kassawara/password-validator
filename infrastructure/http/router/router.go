package router

import (
	"context"
	"fmt"
	"net/http"
	"password-validator/adapter/controller"
	"password-validator/adapter/presenter"
	"password-validator/adapter/repository"
	"password-validator/core/usecase"
	"password-validator/infrastructure/config"
	"sync"
	"time"

	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/logger"
	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/otel/ginmetric"
	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/otel/gintrace"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	Port int64

	Server interface {
		Listen(ctx context.Context, wg *sync.WaitGroup)
	}

	ginEngine struct {
		router                     *gin.Engine
		port                       int64
		ctxTimeout                 time.Duration
		validatePasswordController controller.ValidatePasswordController
	}
)

func NewGinServer() *ginEngine {
	return &ginEngine{
		router: gin.New(),
	}
}

func (engine *ginEngine) WithPort(p int64) *ginEngine {
	engine.port = p
	return engine
}

func (engine *ginEngine) WithDuration(duration time.Duration) *ginEngine {
	engine.ctxTimeout = duration
	return engine
}

func (engine *ginEngine) WithControllers() *ginEngine {
	passwordRepository := repository.NewPasswordRepository()
	validatePasswordPresenter := presenter.NewValidatePasswordPresenter()
	validatePasswordUseCase := usecase.NewValidatePasswordUseCase(engine.ctxTimeout, passwordRepository, validatePasswordPresenter)
	engine.validatePasswordController = controller.NewValidatePasswordController(validatePasswordUseCase)
	return engine
}

func (engine ginEngine) Listen(ctx context.Context, wg *sync.WaitGroup) {
	gin.Recovery()

	engine.setAppHandlers(engine.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", engine.port),
		Handler:      engine.router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			logger.FromContext(context.Background()).Error("Error starting HTTP server", err)
		}
		if gin.Mode() == gin.DebugMode {
			pprof.Register(engine.router)
			go func() {
				http.ListenAndServe(":6060", nil)
			}()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log := logger.FromContext(ctx)
		log.Info("Shutting down HTTP server...")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("HTTP server forced to shutdown", err)
		}
		log.Info("HTTP server exiting")
	}()
}

func (engine ginEngine) setAppHandlers(router *gin.Engine) {
	ginZapConfig := &ginzap.Config{
		SkipPaths:  []string{"/health", "/metrics", "/favicon.ico"},
		TimeFormat: time.RFC3339Nano,
		UTC:        true,
	}
	router.Use(ginzap.GinzapWithConfig(logger.LogProvider(), ginZapConfig))
	router.Use(ginzap.RecoveryWithZap(logger.LogProvider(), true))
	router.Use(gintrace.Middleware(config.C.AppName, gintrace.SkipUselessRoutesTraceOption()))
	router.Use(ginmetric.Middleware(config.C.AppName, ginmetric.WithShouldRecordFunc(ginmetric.SkipUselessMetric)))
	router.Use(logger.Middleware())

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "UP"}) })
	router.POST("/password/validate", engine.handleValidatePassword())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// Validate Password godoc
//
//	@Summary		Validate password
//	@Description	Validates a password according to security rules
//	@Tags			Password
//	@Accept			json
//	@Produce		json
//	@Param			request	body		input.PasswordInput		true	"Password validation request"
//	@Success		200		{object}	output.PasswordOutput	"Validation result"
//	@Failure		422		{object}	response.Error			"Validation error"
//	@Router			/password/validate [post]
func (engine ginEngine) handleValidatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		engine.validatePasswordController.Execute(ctx.Writer, ctx.Request)
	}
}
