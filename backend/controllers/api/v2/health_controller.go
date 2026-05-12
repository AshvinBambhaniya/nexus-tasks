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

// Overall check overall health of application as well as dependencies health check
// swagger:route GET /healthz Healthcheck overallHealthCheck
//
//	Overall health check
//
//	Overall health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Overall(ctx *fiber.Ctx) error {
	err := hc.healthService.CheckDatabaseHealth(ctx.Context())
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Self returns the status of the application.
func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Db handles database health check.
// swagger:route GET /healthz/db Healthcheck dbHealthCheck
//
//	Database health check
//
//	Database health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Db(ctx *fiber.Ctx) error {
	err := hc.healthService.CheckDatabaseHealth(ctx.Context())
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}
