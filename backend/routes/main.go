// Package routes handles the registration of all API routes.
package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	controller "github.com/AshvinBambhaniya/nexus-tasks/v2/controllers/api/v2"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/middlewares"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/watermill"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config *config.AppConfig, pub *watermill.Publisher) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger))
	app.Use(middlewares.SentryMiddleware())

	// Serve Swagger UI at /api/docs — reads the pre-generated ./assets/swagger.json.
	// Run `make swagger-gen` to regenerate the spec after changing annotations.
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v2/",
		FilePath: "./assets/swagger.json",
		Path:     "docs",
		Title:    "Nexus Tasks API Docs",
		CacheAge: 3600,
	}))

	hub := realtime.NewHub(logger)
	go hub.Run()

	// Initialize the Storage (Unit of Work)
	storage := models.NewStorage(goqu)

	// Initialize the Services
	apiKeyService := services.NewAPIKeyService(storage, logger)
	userService := services.NewUserService(storage, logger, config)
	workspaceService := services.NewWorkspaceService(storage, logger, pub)
	teamService := services.NewTeamService(storage, logger)
	projectService := services.NewProjectService(storage, logger)
	taskService := services.NewTaskService(storage, logger, hub)
	websocketService := services.NewWebsocketService(workspaceService, projectService, logger)
	commentService := services.NewCommentService(storage, projectService, logger)
	healthService := services.NewHealthService(storage, logger)
	aiService := services.NewAiService(config, logger)
	notificationService := services.NewNotificationService(storage, logger)
	timeTrackingService := services.NewTimeTrackingService(storage, logger, hub)

	router := app.Group("/api")
	v2 := router.Group("/v2")

	middlewares := middlewares.NewMiddleware(goqu, config, logger, apiKeyService)

	err := healthCheckController(app, healthService, logger)
	if err != nil {
		return err
	}

	err = setupAuthController(v2, userService, logger, config, middlewares)
	if err != nil {
		return err
	}

	err = setupWorkspaceController(v2, workspaceService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupTeamController(v2, teamService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupProjectController(v2, projectService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupTaskController(v2, taskService, commentService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupAiController(v2, aiService, commentService, taskService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupWebsocketController(app, logger, config, hub, websocketService)
	if err != nil {
		return err
	}

	err = setupNotificationController(v2, notificationService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupTimeTrackingController(v2, timeTrackingService, logger, middlewares)
	if err != nil {
		return err
	}

	err = setupAPIKeyController(v2, apiKeyService, logger, middlewares)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func healthCheckController(app *fiber.App, healthService services.HealthService, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(healthService, logger)
	if err != nil {
		return err
	}

	healthz := app.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/self", healthController.Self)
	healthz.Get("/db", healthController.Db)
	return nil
}

func setupAuthController(v2 fiber.Router, userSvc services.UserService, logger *zap.Logger, cfg *config.AppConfig, middlewares middlewares.Middleware) error {
	authController, err := controller.NewAuthController(userSvc, logger, cfg)
	if err != nil {
		return err
	}

	auth := v2.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Get("/me", middlewares.Authenticated, authController.Me)

	return nil
}

func setupWorkspaceController(v2 fiber.Router, workspaceService services.WorkspaceService, logger *zap.Logger, middlewares middlewares.Middleware) error {
	wsController, err := controller.NewWorkspaceController(workspaceService, logger)
	if err != nil {
		return err
	}

	ws := v2.Group("/workspaces", middlewares.Authenticated)
	ws.Post("/", wsController.Create)
	ws.Get("/", wsController.List)

	// /api/v2/workspaces/:workspaceId/members
	wsMember := ws.Group(fmt.Sprintf("/:%s/members", constants.ParamWorkspaceID), middlewares.CheckAccess)
	wsMember.Get("/", wsController.ListMembers)
	wsMember.Post("/", wsController.InviteMember)
	wsMember.Delete("/", middlewares.CheckAccess, wsController.RemoveMember)

	return nil
}

func setupTeamController(v2 fiber.Router, teamService services.TeamService, logger *zap.Logger, middlewares middlewares.Middleware) error {
	teamController, err := controller.NewTeamController(teamService, logger)
	if err != nil {
		return err
	}

	// /api/v2/workspaces/:workspaceId/teams
	teams := v2.Group(fmt.Sprintf("/workspaces/:%s/teams", constants.ParamWorkspaceID), middlewares.Authenticated, middlewares.CheckAccess)
	teams.Post("/", teamController.Create)
	teams.Get("/", teamController.List)
	teams.Get(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Get)
	teams.Patch(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Update)
	teams.Delete(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Delete)

	// Team Members
	teamMember := teams.Group(fmt.Sprintf("/:%s/members", constants.ParamTeamID))
	teamMember.Get("/", teamController.ListMembers)
	teamMember.Post("/", teamController.AddMember)
	teamMember.Delete(fmt.Sprintf("/:%s", constants.ParamUID), teamController.RemoveMember)

	return nil
}

func setupProjectController(v2 fiber.Router, projectService services.ProjectService, logger *zap.Logger, middleware middlewares.Middleware) error {
	projectController, err := controller.NewProjectController(projectService, logger)
	if err != nil {
		return err
	}

	// 1. Workspace Projects
	// /api/v2/workspaces/:workspaceId/projects
	wsProjects := v2.Group(fmt.Sprintf("/workspaces/:%s/projects", constants.ParamWorkspaceID), middleware.Authenticated, middleware.CheckAccess)
	wsProjects.Post("/", projectController.Create)
	wsProjects.Get("/", projectController.List)
	wsProjects.Get(fmt.Sprintf("/:%s", constants.ParamProjectID), projectController.Get)
	wsProjects.Patch(fmt.Sprintf("/:%s", constants.ParamProjectID), projectController.Update)

	// 2. Project Members
	// /api/v2/projects/:projectId/members
	projectMembers := wsProjects.Group(fmt.Sprintf("/:%s/members", constants.ParamProjectID))
	projectMembers.Get("/", projectController.ListMembers)
	projectMembers.Post("/", projectController.AddMember)
	projectMembers.Delete(fmt.Sprintf("/:%s", constants.ParamUID), projectController.RemoveMember)

	// 3. Project Teams
	// /api/v2/projects/:projectId/teams
	projectTeams := wsProjects.Group(fmt.Sprintf("/:%s/teams", constants.ParamProjectID))
	projectTeams.Get("/", projectController.ListTeams)
	projectTeams.Post("/", projectController.AddTeam)
	projectTeams.Delete(fmt.Sprintf("/:%s", constants.ParamTeamID), projectController.RemoveTeam)

	return nil
}

func setupTaskController(v2 fiber.Router, taskService services.TaskService, commentService services.CommentService, logger *zap.Logger, authMiddleware middlewares.Middleware) error {
	taskController, err := controller.NewTaskController(taskService, commentService, logger)
	if err != nil {
		return err
	}

	// 1. Create & List Tasks (Project context)
	// /api/v2/projects/:projectId/tasks
	projectTasks := v2.Group(fmt.Sprintf("/workspaces/:%s/projects/:%s/tasks", constants.ParamWorkspaceID, constants.ParamProjectID), authMiddleware.Authenticated)
	projectTasks.Post("/", taskController.CreateTask)
	projectTasks.Get("/", taskController.ListProjectTasks)

	// 2. Task Management
	tasks := projectTasks.Group(fmt.Sprintf("/:%s", constants.ParamTaskID))
	tasks.Get("/", taskController.GetTask)
	tasks.Patch("/", taskController.UpdateTask)
	tasks.Delete("/", taskController.DeleteTask)

	// Comments
	tasks.Get("/comments", taskController.ListTaskComments)
	tasks.Post("/comments", taskController.CreateComment)
	tasks.Delete(fmt.Sprintf("/comments/:%s", constants.ParamCommentID), taskController.DeleteComment)

	// 3. Global Task Routes
	v2.Get("/tasks/me", authMiddleware.Authenticated, taskController.ListMyTasks)
	v2.Get(fmt.Sprintf("/tasks/:%s", constants.ParamTaskID), authMiddleware.Authenticated, taskController.GetTask)
	v2.Get(fmt.Sprintf("/tasks/:%s/comments", constants.ParamTaskID), authMiddleware.Authenticated, taskController.ListTaskComments)
	v2.Post(fmt.Sprintf("/tasks/:%s/comments", constants.ParamTaskID), authMiddleware.Authenticated, taskController.CreateComment)

	return nil
}

func setupWebsocketController(app *fiber.App, logger *zap.Logger, cfg *config.AppConfig, hub realtime.IHub, wsSvc services.WebsocketService) error {
	wsController, err := controller.NewWebsocketController(hub, wsSvc, cfg, logger)
	if err != nil {
		return err
	}

	app.Use("/ws", wsController.UpgradeMiddleware)
	app.Get("/ws/:id", websocket.New(wsController.HandleWorkspaceConnection))

	return nil
}

func setupAiController(v2 fiber.Router, aiService services.AiService, commentService services.CommentService, taskService services.TaskService, logger *zap.Logger, middleware middlewares.Middleware) error {
	aiController, err := controller.NewAiController(aiService, commentService, taskService, logger)
	if err != nil {
		return err
	}

	ai := v2.Group("/ai", middleware.Authenticated)
	ai.Post("/draft-task", aiController.DraftTask)

	taskAi := v2.Group(fmt.Sprintf("/workspaces/:%s/projects/:%s/tasks/:%s/ai", constants.ParamWorkspaceID, constants.ParamProjectID, constants.ParamTaskID), middleware.Authenticated)
	taskAi.Post("/summarize-comments", aiController.SummarizeComments)

	projectAi := v2.Group(fmt.Sprintf("/workspaces/:%s/projects/:%s/ai", constants.ParamWorkspaceID, constants.ParamProjectID), middleware.Authenticated)
	projectAi.Post("/generate-weekly-report", aiController.GenerateWeeklyReport)

	return nil
}

func setupNotificationController(v2 fiber.Router, notificationService services.NotificationService, logger *zap.Logger, middlewares middlewares.Middleware) error {
	notificationController := controller.NewNotificationController(notificationService, logger)

	inbox := v2.Group("/inbox", middlewares.Authenticated)
	inbox.Get("/", notificationController.GetInbox)
	inbox.Patch("/clear-all", notificationController.ClearAll)
	inbox.Patch("/:notificationId/read", notificationController.MarkAsRead)
	inbox.Patch("/:notificationId/clear", notificationController.MarkAsCleared)

	return nil
}

func setupAPIKeyController(v2 fiber.Router, apiKeyService services.APIKeyService, logger *zap.Logger, middlewares middlewares.Middleware) error {
	apiKeyController := controller.NewAPIKeyController(apiKeyService, logger)

	apiKeys := v2.Group("/auth/api-keys", middlewares.Authenticated)
	apiKeys.Post("/", apiKeyController.Create)
	apiKeys.Get("/", apiKeyController.List)
	apiKeys.Delete(fmt.Sprintf("/:%s", constants.ParamKeyID), apiKeyController.Revoke)

	return nil
}

func setupTimeTrackingController(v2 fiber.Router, service services.TimeTrackingService, logger *zap.Logger, middleware middlewares.Middleware) error {
	ctrl, err := controller.NewTimeTrackingController(service, logger)
	if err != nil {
		return err
	}

	// Active timer
	v2.Get("/timer/active", middleware.Authenticated, ctrl.GetActiveTimer)

	// Task-scoped timer operations
	tasks := v2.Group(fmt.Sprintf("/tasks/:%s", constants.ParamTaskID), middleware.Authenticated)
	tasks.Post("/timer/start", ctrl.StartTimer)
	tasks.Post("/timer/stop", ctrl.StopTimer)
	tasks.Post("/timer/discard", ctrl.DiscardTimer)
	tasks.Post("/time-entries", ctrl.LogManualTime)
	tasks.Get("/time-entries", ctrl.ListTaskTimeEntries)

	// Time entry management
	v2.Delete(fmt.Sprintf("/time-entries/:%s", constants.ParamEntryID), middleware.Authenticated, ctrl.DeleteTimeEntry)

	// Project analytics (nested under workspace/project)
	projectAnalytics := v2.Group(fmt.Sprintf("/workspaces/:%s/projects/:%s", constants.ParamWorkspaceID, constants.ParamProjectID), middleware.Authenticated)
	projectAnalytics.Get("/time-analytics", ctrl.GetProjectAnalytics)
	projectAnalytics.Get("/time-entries", ctrl.ListProjectTimeEntries)

	return nil
}
