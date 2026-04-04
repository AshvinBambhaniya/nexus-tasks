package routes

import (
	"sync"

	"go.uber.org/zap"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	controller "github.com/AshvinBambhaniya/nexus-tasks/controllers/api/v1"
	"github.com/AshvinBambhaniya/nexus-tasks/middlewares"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger))
	app.Use(middlewares.SentryMiddleware())

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./assets/swagger.json",
		Path:     "docs",
		Title:    "Swagger API Docs",
	}))

	router := app.Group("/api")
	_ = router.Group("/v1")

	_ = middlewares.NewMiddleware(config, logger)

	err := healthCheckController(app, goqu, logger)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func healthCheckController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := app.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}
