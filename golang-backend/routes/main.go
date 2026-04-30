package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	controller "github.com/AshvinBambhaniya/nexus-tasks/controllers/api/v1"
	"github.com/AshvinBambhaniya/nexus-tasks/middlewares"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, pub *watermill.WatermillPublisher) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger))
	app.Use(middlewares.SentryMiddleware())

	// TODO: Setup swagger docs
	// app.Use(swagger.New(swagger.Config{
	// 	BasePath: "/api/v1/",
	// 	FilePath: "./assets/swagger.json",
	// 	Path:     "docs",
	// 	Title:    "Swagger API Docs",
	// }))

	hub := realtime.NewHub(logger)
	go hub.Run()

	router := app.Group("/api")
	v1 := router.Group("/v1")

	middlewares, err := middlewares.NewMiddleware(goqu, config, logger)
	if err != nil {
		return err
	}

	err = healthCheckController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = setupAuthController(v1, goqu, logger, config, middlewares)
	if err != nil {
		return err
	}

	err = setupWorkspaceController(v1, goqu, logger, config, middlewares, pub)
	if err != nil {
		return err
	}

	err = setupTeamController(v1, goqu, logger, config, middlewares)
	if err != nil {
		return err
	}

	err = setupProjectController(v1, goqu, logger, config, middlewares)
	if err != nil {
		return err
	}

	err = setupTaskController(v1, goqu, logger, config, middlewares, hub)
	if err != nil {
		return err
	}

	err = setupWebsocketController(app, goqu, logger, config, hub)
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

func setupAuthController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, middlewares middlewares.Middleware) error {
	authController, err := controller.NewAuthController(goqu, logger, cfg)
	if err != nil {
		return err
	}

	auth := v1.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Get("/me", middlewares.Authenticated, authController.Me)

	return nil
}

func setupWorkspaceController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, middlewares middlewares.Middleware, publisher *watermill.WatermillPublisher) error {
	wsController, err := controller.NewWorkspaceController(goqu, logger, cfg, publisher)
	if err != nil {
		return err
	}

	ws := v1.Group("/workspaces", middlewares.Authenticated)
	ws.Post("/", wsController.Create)
	ws.Get("/", wsController.List)

	// /api/v1/workspaces/:workspaceId/members
	wsMember := ws.Group(fmt.Sprintf("/:%s/members", constants.ParamWorkspaceID), middlewares.CheckAccess)
	wsMember.Get("/", wsController.ListMembers)
	wsMember.Post("/", wsController.InviteMember)
	wsMember.Delete("/", middlewares.CheckAccess, wsController.RemoveMember)

	return nil
}

func setupTeamController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, middlewares middlewares.Middleware) error {
	teamController, err := controller.NewTeamController(goqu, logger, cfg)
	if err != nil {
		return err
	}

	// /api/v1/workspaces/:workspaceId/teams
	teams := v1.Group(fmt.Sprintf("/workspaces/:%s/teams", constants.ParamWorkspaceID), middlewares.Authenticated, middlewares.CheckAccess)
	teams.Post("/", teamController.Create)
	teams.Get("/", teamController.List)
	teams.Get(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Get)
	teams.Patch(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Update)
	teams.Delete(fmt.Sprintf("/:%s", constants.ParamTeamID), teamController.Delete)

	// Team Members
	teamMember := teams.Group(fmt.Sprintf("/:%s/members", constants.ParamTeamID))
	teamMember.Get("/", teamController.ListMembers)
	teamMember.Post("/", teamController.AddMember)
	teamMember.Delete(fmt.Sprintf("/:%s", constants.ParamUid), teamController.RemoveMember)

	return nil
}

func setupProjectController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, middleware middlewares.Middleware) error {
	projectController, err := controller.NewProjectController(goqu, logger, cfg)
	if err != nil {
		return err
	}

	// 1. Workspace Projects
	// /api/v1/workspaces/:workspaceId/projects
	wsProjects := v1.Group(fmt.Sprintf("/workspaces/:%s/projects", constants.ParamWorkspaceID), middleware.Authenticated, middleware.CheckAccess)
	wsProjects.Post("/", projectController.Create)
	wsProjects.Get("/", projectController.List)
	wsProjects.Get(fmt.Sprintf("/:%s", constants.ParamProjectID), projectController.Get)
	wsProjects.Patch(fmt.Sprintf("/:%s", constants.ParamProjectID), projectController.Update)

	// 2. Project Members
	// /api/v1/projects/:projectId/members
	projectMembers := wsProjects.Group(fmt.Sprintf("/:%s/members", constants.ParamProjectID))
	projectMembers.Get("/", projectController.ListMembers)
	projectMembers.Post("/", projectController.AddMember)
	projectMembers.Delete(fmt.Sprintf("/:%s", constants.ParamUid), projectController.RemoveMember)

	// 3. Project Teams
	// /api/v1/projects/:projectId/teams
	projectTeams := wsProjects.Group(fmt.Sprintf("/:%s/teams", constants.ParamProjectID))
	projectTeams.Get("/", projectController.ListTeams)
	projectTeams.Post("/", projectController.AddTeam)
	projectTeams.Delete(fmt.Sprintf("/:%s", constants.ParamTeamID), projectController.RemoveTeam)

	return nil
}

func setupTaskController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, authMiddleware middlewares.Middleware, hub *realtime.Hub) error {
	taskController, err := controller.NewTaskController(goqu, logger, cfg, hub)
	if err != nil {
		return err
	}

	// 1. Create & List Tasks (Project context)
	// /api/v1/projects/:projectId/tasks
	projectTasks := v1.Group(fmt.Sprintf("/workspaces/:%s/projects/:%s/tasks", constants.ParamWorkspaceID, constants.ParamProjectID), authMiddleware.Authenticated)
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
	v1.Get("/tasks/me", authMiddleware.Authenticated, taskController.ListMyTasks)

	return nil
}

func setupWebsocketController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, hub *realtime.Hub) error {
	wsController, err := controller.NewWebsocketController(hub, goqu, cfg, logger)
	if err != nil {
		return err
	}

	app.Use("/ws", wsController.UpgradeMiddleware)
	app.Get("/ws/:id", websocket.New(wsController.HandleWorkspaceConnection))

	return nil
}
