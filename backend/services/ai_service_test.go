package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupAiTest(_ *testing.T) AiService {
	// Initialize without an OpenAI key to test the fallback mock logic
	cfg := &config.AppConfig{}
	logger := zap.NewNop()
	return NewAiService(cfg, logger)
}

func TestAiService_DraftTask(t *testing.T) {
	t.Run("fallback mock response", func(t *testing.T) {
		svc := setupAiTest(t)
		title := "Test Bug Fix"

		res, err := svc.DraftTask(title)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.True(t, strings.Contains(res, title))
		assert.True(t, strings.Contains(res, "## Overview"))
		assert.True(t, strings.Contains(res, "## Acceptance Criteria"))
	})
}

func TestAiService_SummarizeComments(t *testing.T) {
	t.Run("fallback mock response", func(t *testing.T) {
		svc := setupAiTest(t)
		comments := "User A: this is a test comment."

		res, err := svc.SummarizeComments(comments)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.True(t, strings.Contains(res, "**Summary**"))
	})
}

func TestAiService_GenerateWeeklyReport(t *testing.T) {
	t.Run("fallback mock response", func(t *testing.T) {
		svc := setupAiTest(t)
		tasksData := "{}"

		res, err := svc.GenerateWeeklyReport(tasksData)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.True(t, strings.Contains(res, "**Weekly Report (Mock)**"))
	})
}

func setupLiveAiTest(_ *testing.T, handler http.HandlerFunc) (AiService, *httptest.Server) {
	ts := httptest.NewServer(handler)
	cfg := &config.AppConfig{
		OpenAIKey:     "fake-key",
		OpenAIBaseURL: ts.URL, // Go-OpenAI appends /chat/completions to BaseURL
	}
	logger := zap.NewNop()
	return NewAiService(cfg, logger), ts
}

func TestAiService_DraftTask_Live(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"choices": [{"message": {"content": "Live Draft output"}}]}`))
		}
		svc, ts := setupLiveAiTest(t, handler)
		defer ts.Close()

		res, err := svc.DraftTask("Test title")
		assert.NoError(t, err)
		assert.Equal(t, "Live Draft output", res)
	})

	t.Run("api error", func(t *testing.T) {
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		svc, ts := setupLiveAiTest(t, handler)
		defer ts.Close()

		res, err := svc.DraftTask("Test title")
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestAiService_SummarizeComments_Live(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"choices": [{"message": {"content": "Live Summary output"}}]}`))
		}
		svc, ts := setupLiveAiTest(t, handler)
		defer ts.Close()

		res, err := svc.SummarizeComments("Comments here")
		assert.NoError(t, err)
		assert.Equal(t, "Live Summary output", res)
	})
}

func TestAiService_GenerateWeeklyReport_Live(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"choices": [{"message": {"content": "Live Weekly Report output"}}]}`))
		}
		svc, ts := setupLiveAiTest(t, handler)
		defer ts.Close()

		res, err := svc.GenerateWeeklyReport("{}")
		assert.NoError(t, err)
		assert.Equal(t, "Live Weekly Report output", res)
	})
}
