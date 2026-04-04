package services

import (
	"database/sql"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type UserService struct {
	userModel      *models.UserModel
	workspaceModel *models.WorkspaceModel
	db             *goqu.Database
	logger         *zap.Logger
}

func NewUserService(db *goqu.Database, logger *zap.Logger, userModel *models.UserModel, workspaceModel *models.WorkspaceModel) *UserService {
	return &UserService{
		userModel:      userModel,
		workspaceModel: workspaceModel,
		db:             db,
		logger:         logger,
	}
}

func (s *UserService) Register(email, password, fullName string) (models.User, error) {

	// 1. Hash Password
	hashedPwd, err := utils.PasswordHash(password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Email:          email,
		FullName:       fullName,
		HashedPassword: hashedPwd,
		IsActive:       true,
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return user, err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				s.logger.Error("error during commit in register user", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				s.logger.Error("error during rollback in register user", zap.Error(err))
			}
		}
	}()

	// A. Create User
	createdUser, err := s.userModel.CreateUser(transaction, user)
	if err != nil {
		return user, err
	}

	// B. Create Personal Workspace
	personalWs := models.Workspace{
		Name:    constants.WorkspaceTypePersonal,
		Type:    models.WorkspaceTypePersonal,
		OwnerID: createdUser.ID,
	}
	createdWs, err := s.workspaceModel.CreateWorkspace(transaction, personalWs)
	if err != nil {
		return user, err
	}

	// C. Add Member (Admin)
	err = s.workspaceModel.AddMember(transaction, models.WorkspaceMember{
		WorkspaceID: createdWs.ID,
		UserID:      createdUser.ID,
		Role:        models.WorkspaceRoleAdmin,
	})
	if err != nil {
		return user, err
	}

	isOk = true
	return createdUser, nil
}

func (s *UserService) GetUser(id int) (models.User, error) {
	return s.userModel.GetByID(id)
}

func (s *UserService) Authenticate(email, password string) (models.User, error) {
	user, err := s.userModel.GetByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	if !utils.CheckPasswordHash(password, user.HashedPassword) {
		return models.User{}, sql.ErrNoRows // Treat as not found/invalid
	}

	return user, nil
}
