package v2

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HealthController handles health check requests.
type HealthController struct {
	healthService services.HealthService
	logger        *zap.Logger
}

// NewHealthController creates a new instance of HealthController.
func NewHealthController(healthService services.HealthService, logger *zap.Logger) (*HealthController, error) {
	return &HealthController{
		healthService: healthService,
		logger:        logger,
	}, nil
}

// Overall check overall health of application as well as dependencies health check.
//
/*
 swagger:operation GET /healthz health overallHealthCheck

 # Overall system health check

 Returns 200 if the API server and its critical dependencies (database) are healthy.

 ---
 produces:
 - application/json
 responses:

	200:
	  description: All systems healthy.
	500:
	  description: One or more dependencies are unhealthy.
*/
func (hc *HealthController) Overall(ctx *fiber.Ctx) error {
	err := hc.healthService.CheckDatabaseHealth(ctx.Context())
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Self returns the status of the application process.
//
/*
 swagger:operation GET /healthz/self health selfHealthCheck

 # Self / liveness check

 Returns 200 as long as the HTTP server process is running. Does not check external dependencies.

 ---
 produces:
 - application/json
 responses:

	200:
	  description: Server process is alive.
*/
func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Db handles database health check.
//
/*
 swagger:operation GET /healthz/db health dbHealthCheck

 # Database connectivity check

 Performs a lightweight ping against the primary database and returns 200 if reachable.

 ---
 produces:
 - application/json
 responses:

	200:
	  description: Database is reachable.
	500:
	  description: Database connectivity check failed.
*/
func (hc *HealthController) Db(ctx *fiber.Ctx) error {
	err := hc.healthService.CheckDatabaseHealth(ctx.Context())
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}
