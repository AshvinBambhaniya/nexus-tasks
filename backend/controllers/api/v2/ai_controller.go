package v2

import (
	"fmt"
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// AiController handles AI related requests.
type AiController struct {
	aiService      services.AiService
	commentService services.CommentService
	taskService    services.TaskService
	logger         *zap.Logger
}

// NewAiController creates a new instance of AiController.
func NewAiController(aiService services.AiService, commentService services.CommentService, taskService services.TaskService, logger *zap.Logger) (*AiController, error) {
	return &AiController{
		aiService:      aiService,
		commentService: commentService,
		taskService:    taskService,
		logger:         logger,
	}, nil
}

// DraftTask handles the POST /api/v2/ai/draft-task endpoint.
/*
swagger:operation POST /ai/draft-task ai draftTask

# AI — Draft a task description

Uses an AI model to generate a detailed task description based on the provided title.
The response contains the generated markdown content as a string.

---
consumes:
- application/json
produces:
- application/json
security:
- cookieAuth: []
- apiKeyAuth: []
parameters:
  - name: body
    in: body
    required: true
    schema:
       "$ref": "#/definitions/ReqDraftTask"

responses:
  "200":
    description: AI-generated task description.
  "400":
    description: Validation error.
  "401":
    description: Not authenticated.
  "500":
    description: AI service error or internal server error.
*/
func (ctrl *AiController) DraftTask(c *fiber.Ctx) error {
	var req structs.ReqDraftTask
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	content, err := ctrl.aiService.DraftTask(req.Title)
	if err != nil {
		ctrl.logger.Error("failed to draft task via AI", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to draft task")
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{
		constants.ResponseKeyContent: content,
	})
}

// SummarizeComments handles the POST /api/v2/workspaces/:wsID/projects/:projectID/tasks/:taskID/ai/summarize-comments endpoint.
//
/*
 swagger:operation POST /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId}/ai/summarize-comments ai summarizeComments

 # AI — Summarize task comments

 Uses an AI model to produce a concise summary of all comments on the specified task.
 Returns an error if the task has no comments.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: AI-generated comment summary.
	400:
	  description: No comments to summarize, or invalid path parameter.
	401:
	  description: Not authenticated.
	500:
	  description: AI service error or internal server error.
*/
func (ctrl *AiController) SummarizeComments(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	comments, err := ctrl.commentService.ListTaskComments(uid, taskID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	if len(comments) == 0 {
		return utils.JSONFail(c, http.StatusBadRequest, "No comments to summarize")
	}

	// Build a transcript
	transcript := ""
	for _, comment := range comments {
		author := comment.AuthorFullName
		if author == "" {
			author = comment.AuthorEmail
		}
		transcript += fmt.Sprintf("%s: %s\n\n", author, comment.Content)
	}

	summary, err := ctrl.aiService.SummarizeComments(transcript)
	if err != nil {
		ctrl.logger.Error("failed to summarize comments via AI", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to summarize comments")
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{
		constants.ResponseKeyContent: summary,
	})
}

// GenerateWeeklyReport handles POST /api/v2/workspaces/:wsID/projects/:projectID/ai/generate-weekly-report.
//
/*
 swagger:operation POST /workspaces/{workspaceId}/projects/{projectId}/ai/generate-weekly-report ai generateWeeklyReport

 # AI — Generate a weekly project report

 Uses an AI model to generate a structured weekly progress report for the project,
 based on tasks completed in the last 7 days and their associated comments.
 Returns an error if no tasks were completed in that window.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: AI-generated weekly report content.
	400:
	  description: No completed tasks found in the last 7 days.
	401:
	  description: Not authenticated.
	500:
	  description: AI service error or internal server error.
*/
func (ctrl *AiController) GenerateWeeklyReport(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
	}

	tasks, commentsByTask, err := ctrl.fetchWeeklyTasksAndComments(uid, projectID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	if len(tasks) == 0 {
		return utils.JSONFail(c, http.StatusBadRequest, "No completed tasks found in the last 7 days.")
	}

	tasksData := ctrl.compilePromptData(tasks, commentsByTask)

	report, err := ctrl.aiService.GenerateWeeklyReport(tasksData)
	if err != nil {
		ctrl.logger.Error("failed to generate weekly report via AI", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to generate weekly report")
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{
		constants.ResponseKeyContent: report,
	})
}

func (ctrl *AiController) fetchWeeklyTasksAndComments(uid, projectID uuid.UUID) ([]models.TaskWithAssignee, map[uuid.UUID][]string, error) {
	tasks, err := ctrl.taskService.ListCompletedTasksInLastDays(uid, projectID, 7)
	if err != nil {
		return nil, nil, err
	}
	if len(tasks) == 0 {
		return tasks, nil, nil
	}

	var taskIDs []uuid.UUID
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	comments, err := ctrl.commentService.ListCommentsForTasks(uid, projectID, taskIDs)
	if err != nil {
		return nil, nil, err
	}

	commentsByTask := make(map[uuid.UUID][]string)
	for _, c := range comments {
		author := c.AuthorFullName
		if author == "" {
			author = c.AuthorEmail
		}
		commentsByTask[c.TaskID] = append(commentsByTask[c.TaskID], fmt.Sprintf("%s: %s", author, c.Content))
	}

	return tasks, commentsByTask, nil
}

func (ctrl *AiController) compilePromptData(tasks []models.TaskWithAssignee, commentsByTask map[uuid.UUID][]string) string {
	var tasksData string
	for _, t := range tasks {
		assignee := "Unassigned"
		if t.AssigneeFullName != nil {
			assignee = *t.AssigneeFullName
		} else if t.AssigneeEmail != nil {
			assignee = *t.AssigneeEmail
		}

		deadlineMissed := "No"
		if t.DueDate != nil && t.CompletedAt != nil && t.CompletedAt.After(*t.DueDate) {
			deadlineMissed = "Yes"
		}

		tasksData += fmt.Sprintf("Task: %s\nAssignee: %s\nDeadline Missed: %s\nDescription: %s\n", t.Title, assignee, deadlineMissed, t.Description)

		taskComments := commentsByTask[t.ID]
		if len(taskComments) > 0 {
			tasksData += "Comments:\n"
			for _, comment := range taskComments {
				tasksData += fmt.Sprintf("- %s\n", comment)
			}
		}
		tasksData += "\n---\n\n"
	}
	return tasksData
}
