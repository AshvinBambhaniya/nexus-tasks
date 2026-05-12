// Package v2 provides the version 1 of the API controllers.
package v2

import (
	"net/http"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// AuthController handles authentication related requests.
type AuthController struct {
	userService services.UserService
	logger      *zap.Logger
	config      *config.AppConfig
}

const cookieSameSiteLax = "Lax"

// NewAuthController creates a new instance of AuthController.
func NewAuthController(userSvc services.UserService, logger *zap.Logger, cfg *config.AppConfig) (*AuthController, error) {
	return &AuthController{
		userService: userSvc,
		logger:      logger,
		config:      cfg,
	}, nil
}

// Register handles user registration.
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
		SameSite: cookieSameSiteLax,
		MaxAge:   ctrl.config.JwtExpirationHours * 60 * 60,
	})

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResUser{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	})
}

// Login handles user login and sets the authentication cookie.
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
		SameSite: cookieSameSiteLax,
		MaxAge:   ctrl.config.JwtExpirationHours * 60 * 60,
	})

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: "Login successful"})
}

// Logout handles user logout by clearing the authentication cookie.
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieUser,
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: cookieSameSiteLax,
		Expires:  time.Now().Add(-24 * time.Hour),
		MaxAge:   -1,
	})
	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: "Logout successful"})
}

// Me returns the current authenticated user's information.
func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	// Got from middleware
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error(constants.ErrInvalidUserID, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
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
