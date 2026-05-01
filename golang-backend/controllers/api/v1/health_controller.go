package v1

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthController struct {
	healthService services.HealthService
	logger        *zap.Logger
}

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

func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Database health check
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
