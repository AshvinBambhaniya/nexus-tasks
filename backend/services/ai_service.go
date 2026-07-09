package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// AiService defines the AI service interface
type AiService interface {
	DraftTask(title string) (string, error)
	SummarizeComments(commentsText string) (string, error)
	GenerateWeeklyReport(tasksData string) (string, error)
}

// aiService implementation
type aiService struct {
	config *config.AppConfig
	logger *zap.Logger
	client *openai.Client
}

// NewAiService creates a new AiService
func NewAiService(cfg *config.AppConfig, logger *zap.Logger) AiService {
	var client *openai.Client
	if cfg.OpenAIKey != "" {
		openaiConfig := openai.DefaultConfig(cfg.OpenAIKey)
		if cfg.OpenAIBaseURL != "" {
			openaiConfig.BaseURL = cfg.OpenAIBaseURL
		}
		client = openai.NewClientWithConfig(openaiConfig)
	}

	return &aiService{
		config: cfg,
		logger: logger,
		client: client,
	}
}

// DraftTask generates a task description using OpenAI
func (s *aiService) DraftTask(title string) (string, error) {
	if s.client == nil {
		s.logger.Warn("OpenAI client is missing. Using mock response for AI Draft.")
		return s.generateMockDraft(title), nil
	}

	systemPrompt := `You are a Senior Technical Product Manager. The user will provide a brief, vague title for a software task or feature.
Your job is to break it down into a comprehensive, professional Markdown description.
The markdown MUST include:
1. **Overview**: A 2-3 sentence summary of the feature.
2. **Acceptance Criteria**: A bulleted list of what constitutes "Done".
3. **Implementation Steps**: A markdown checklist (using '- [ ] ') of 3-5 subtasks required to build this feature.

Return ONLY the markdown content, nothing else.`

	model := s.config.OpenAIModel
	if model == "" {
		model = openai.GPT4oMini
	}

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: title,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		s.logger.Error("Failed to call OpenAI API", zap.Error(err))
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices found in AI response")
	}

	return resp.Choices[0].Message.Content, nil
}

// SummarizeComments summarizes a thread of comments
func (s *aiService) SummarizeComments(commentsText string) (string, error) {
	if s.client == nil {
		s.logger.Warn("OpenAI API client is missing. Using mock response for AI Summarize.")
		return "**Summary**\n- Discussed integration details.\n- Blocked by backend deployment.\n- Assigned to Ashvin.", nil
	}

	systemPrompt := `You are an expert Project Manager. You will be provided with a transcript of comments on a software task.
Your job is to provide a concise, 3-bullet point summary of the conversation. Focus on:
1. Key decisions made.
2. Current blockers or open questions.
3. Next steps.

Output strictly as a Markdown unordered list. Do not include any intro or outro text.`

	model := s.config.OpenAIModel
	if model == "" {
		model = openai.GPT4oMini
	}

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: commentsText,
				},
			},
			Temperature: 0.5,
		},
	)

	if err != nil {
		s.logger.Error("Failed to call OpenAI API for summarization", zap.Error(err))
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices found in AI response")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateWeeklyReport generates a professional markdown changelog from recent tasks
func (s *aiService) GenerateWeeklyReport(tasksData string) (string, error) {
	if s.client == nil {
		s.logger.Warn("OpenAI API client is missing. Using mock response for AI Weekly Report.")
		return "**Weekly Report (Mock)**\n- Shipped X feature.\n- 1 deadline missed.", nil
	}

	systemPrompt := `You are an expert Technical Writer and Agile Project Manager.
You will be provided with a JSON-like text dump of all tasks completed in the last 7 days, including their comments and whether they missed their deadlines.
Your job is to synthesize this into a beautiful, professional Markdown weekly report / release notes document.

The document MUST contain:
1. **Highlights & Achievements**: A short paragraph summarizing the biggest wins (infer from titles and comments).
2. **Completed Tasks**: A categorized bulleted list of what was done. Add a brief 1-sentence context for each (derived from comments if available).
3. **Missed Deadlines / Bottlenecks**: A section highlighting any tasks that missed their deadlines, and any blockers mentioned in comments. If none, say "All tasks completed on time."

Use engaging and professional language. Use emojis tastefully.`

	model := s.config.OpenAIModel
	if model == "" {
		model = openai.GPT4oMini
	}

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: tasksData,
				},
			},
			Temperature: 0.6,
		},
	)

	if err != nil {
		s.logger.Error("Failed to call OpenAI API for weekly report", zap.Error(err))
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices found in AI response")
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *aiService) generateMockDraft(title string) string {
	return fmt.Sprintf(`## Overview
This task aims to implement **%s**. It will enhance our platform's capabilities and provide a better experience for our users.

## Acceptance Criteria
* [ ] The feature is fully functional and tested.
* [ ] UI matches the high-fidelity designs.
* [ ] Edge cases and error states are handled.

## Implementation Steps
- [ ] Investigate the current architecture and identify integration points.
- [ ] Build the backend API endpoints or background services.
- [ ] Implement the frontend components and state management.
- [ ] Write unit tests and conduct QA validation.`, title)
}
