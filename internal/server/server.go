package server

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/middleware"
	"github.com/laksanagusta/identity/pkg/errorhelper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	goccyjson "github.com/goccy/go-json"
	"github.com/segmentio/encoding/json"
)

type Server struct {
	Config config.Config
	Logger *zap.SugaredLogger
	Fiber  *fiber.App
	DB     *sqlx.DB
}

func NewServer(config config.Config, logger *zap.SugaredLogger, db *sqlx.DB) *Server {
	var fiberConfig fiber.Config
	fiberConfig.ErrorHandler = errorhelper.HttpHandleError
	fiberConfig.AppName = config.App.Name
	fiberConfig.DisableStartupMessage = true
	fiberConfig.JSONEncoder = goccyjson.Marshal
	fiberConfig.JSONDecoder = json.Unmarshal

	return &Server{
		Config: config,
		Logger: logger,
		Fiber:  fiber.New(fiberConfig),
		DB:     db,
	}
}

func (s *Server) Run() error {
	// Request Logger Middleware
	if s.Config.App.Env == "local" {
		config := logger.ConfigDefault
		config.Format = "[${time}] ${status} ${method} ${path}\n"
		s.Fiber.Use(logger.New(config))
	}

	s.Fiber.Use(middleware.CORS())

	// Panic Recovery Middleware - mencegah aplikasi berhenti saat terjadi panic
	// Menggunakan custom middleware untuk logging yang lebih baik
	s.Fiber.Use(middleware.RecoveryMiddleware(s.Logger))

	// Backup: Fiber's built-in recover sebagai safety net tambahan
	s.Fiber.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			stackTrace := debug.Stack()
			s.Logger.Errorw("Fiber recover caught panic (backup handler)",
				"error", e,
				"path", c.Path(),
				"stack_trace", string(stackTrace),
			)
		},
	}))

	// Swagger Handler
	// s.Fiber.Get("/swagger/*", swagger.HandlerDefault)

	// Map App Handlers
	err := s.MapHandlers()
	if err != nil {
		return err
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)
	go func() {
		<-quit
		s.Fiber.Shutdown()
	}()

	// Run Fiber
	s.Logger.Infof("App started")
	return s.Fiber.Listen(fmt.Sprintf(":%s", s.Config.App.Port))
}
