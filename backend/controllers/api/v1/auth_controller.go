package v1

import (
	"net/http"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type AuthController struct {
	userService services.UserService
	logger      *zap.Logger
	config      *config.AppConfig
}

func NewAuthController(userSvc services.UserService, logger *zap.Logger, cfg *config.AppConfig) (*AuthController, error) {
	return &AuthController{
		userService: userSvc,
		logger:      logger,
		config:      cfg,
	}, nil
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req structs.ReqRegisterUser
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, token, err := ctrl.userService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		ctrl.logger.Error("failed to register user", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to register user")
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieUser,
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   ctrl.config.JwtExpirationHours * 60 * 60,
	})

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResUser{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req structs.ReqLoginUser

	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	token, err := ctrl.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		return utils.JSONFail(c, http.StatusUnauthorized, "Incorrect email or password")
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieUser,
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   ctrl.config.JwtExpirationHours * 60 * 60,
	})

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Login successful"})
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieUser,
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(-24 * time.Hour),
		MaxAge:   -1,
	})
	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Logout successful"})
}

func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	// Got from middleware
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error("invalid user id in context", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "invalid user id in context")
	}

	user, err := ctrl.userService.GetByID(uid)
	if err != nil {
		return utils.JSONFail(c, http.StatusNotFound, "User not found")
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResUser{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	})
}
