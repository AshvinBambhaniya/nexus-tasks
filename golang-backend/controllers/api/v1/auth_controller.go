package v1

import (
	"net/http"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/jwt"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type AuthController struct {
	userService *services.UserService
	userModel   *models.UserModel
	logger      *zap.Logger
	config      config.AppConfig
}

func NewAuthController(goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		return nil, err
	}
	workspaceModel, err := models.InitWorkspaceModel(goqu)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(goqu, logger, &userModel, &workspaceModel)

	return &AuthController{
		userService: userSvc,
		userModel:   &userModel,
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

	// Check if user exists
	_, err := ctrl.userModel.GetByEmail(req.Email)
	if err == nil {
		return utils.JSONFail(c, http.StatusBadRequest, "Email already registered")
	}

	user, err := ctrl.userService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		ctrl.logger.Error("failed to register user", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to register user")
	}

	// Generate Token
	token, err := jwt.CreateToken(ctrl.config, user.ID.String(), time.Now().Add(time.Duration(ctrl.config.JwtExpirationHours)*time.Hour))
	if err != nil {
		ctrl.logger.Error("failed to generate token", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to generate token")
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

	user, err := ctrl.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		return utils.JSONFail(c, http.StatusUnauthorized, "Incorrect email or password")
	}

	// Generate Token
	// Using ID (UUID) as Subject.
	token, err := jwt.CreateToken(ctrl.config, user.ID.String(), time.Now().Add(time.Duration(ctrl.config.JwtExpirationHours)*time.Hour))
	if err != nil {
		ctrl.logger.Error("failed to generate token", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to generate token")
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
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error("invalid user id in context", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "invalid user id in context")
	}

	user, err := ctrl.userModel.GetByID(uid)
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
