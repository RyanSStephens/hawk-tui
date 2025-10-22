package tui

import (
	"testing"
	"time"

	"github.com/hawk-tui/hawk-tui/pkg/types"
)

func TestTUIInitialization(t *testing.T) {
	config := Config{
		AppName:  "test-app",
		LogLevel: "info",
		Debug:    false,
	}
	
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Expected no error creating model, got %v", err)
	}

	if model == nil {
		t.Fatal("Expected Model instance, got nil")
	}

	if model.logs == nil {
		t.Fatal("Expected logs slice to be initialized")
	}

	if model.metrics == nil {
		t.Fatal("Expected metrics map to be initialized")
	}

	if model.configs == nil {
		t.Fatal("Expected configs map to be initialized")
	}

	if model.progress == nil {
		t.Fatal("Expected progress map to be initialized")
	}

	if model.events == nil {
		t.Fatal("Expected events slice to be initialized")
	}
}

func TestHandleLogMessage(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test log message handling
	now := time.Now()
	logParams := types.LogParams{
		Message:   "Test log message",
		Level:     types.LogLevelInfo,
		Timestamp: &now,
		Component: "test-app",
	}

	err = model.HandleLog(logParams, nil)
	if err != nil {
		t.Errorf("Expected no error handling log, got %v", err)
	}

	model.mu.RLock()
	defer model.mu.RUnlock()

	if len(model.logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(model.logs))
	}

	log := model.logs[0]
	if log.Message != "Test log message" {
		t.Errorf("Expected message 'Test log message', got '%s'", log.Message)
	}

	if log.Level != types.LogLevelInfo {
		t.Errorf("Expected level 'INFO', got '%s'", log.Level)
	}

	if log.Component != "test-app" {
		t.Errorf("Expected component 'test-app', got '%s'", log.Component)
	}
}

func TestHandleMetricMessage(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test metric message handling
	now := time.Now()
	metricParams := types.MetricParams{
		Name:      "cpu_usage",
		Value:     75.5,
		Type:      types.MetricTypeGauge,
		Unit:      "%",
		Timestamp: &now,
	}

	err = model.HandleMetric(metricParams, nil)
	if err != nil {
		t.Errorf("Expected no error handling metric, got %v", err)
	}

	model.mu.RLock()
	defer model.mu.RUnlock()

	if len(model.metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(model.metrics))
	}

	metric, exists := model.metrics["cpu_usage"]
	if !exists {
		t.Fatal("Expected metric 'cpu_usage' to exist")
	}

	if metric.Value != 75.5 {
		t.Errorf("Expected value 75.5, got %f", metric.Value)
	}

	if metric.Type != types.MetricTypeGauge {
		t.Errorf("Expected type 'gauge', got '%s'", metric.Type)
	}

	if metric.Unit != "%" {
		t.Errorf("Expected unit '%%', got '%s'", metric.Unit)
	}
}

func TestHandleConfigMessage(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test config message handling
	configParams := types.ConfigParams{
		Key:         "server_port",
		Value:       3000,
		Type:        types.ConfigTypeInteger,
		Description: "HTTP server port",
	}

	err = model.HandleConfig(configParams, nil)
	if err != nil {
		t.Errorf("Expected no error handling config, got %v", err)
	}

	model.mu.RLock()
	defer model.mu.RUnlock()

	if len(model.configs) != 1 {
		t.Fatalf("Expected 1 config item, got %d", len(model.configs))
	}

	configEntry, exists := model.configs["server_port"]
	if !exists {
		t.Fatal("Expected config 'server_port' to exist")
	}

	if configEntry.Value != 3000 {
		t.Errorf("Expected value 3000, got %v", configEntry.Value)
	}

	if configEntry.Type != types.ConfigTypeInteger {
		t.Errorf("Expected type 'integer', got '%s'", configEntry.Type)
	}

	if configEntry.Description != "HTTP server port" {
		t.Errorf("Expected description 'HTTP server port', got '%s'", configEntry.Description)
	}
}

func TestHandleProgressMessage(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test progress message handling
	progressParams := types.ProgressParams{
		ID:      "file_upload",
		Label:   "File Upload Progress",
		Current: 75,
		Total:   100,
		Unit:    "files",
		Status:  types.ProgressStatusInProgress,
		Details: "Uploading files...",
	}

	err = model.HandleProgress(progressParams, nil)
	if err != nil {
		t.Errorf("Expected no error handling progress, got %v", err)
	}

	model.mu.RLock()
	defer model.mu.RUnlock()

	if len(model.progress) != 1 {
		t.Fatalf("Expected 1 progress item, got %d", len(model.progress))
	}

	progress, exists := model.progress["file_upload"]
	if !exists {
		t.Fatal("Expected progress 'file_upload' to exist")
	}

	if progress.Current != 75 {
		t.Errorf("Expected current 75, got %f", progress.Current)
	}

	if progress.Total != 100 {
		t.Errorf("Expected total 100, got %f", progress.Total)
	}

	if progress.Unit != "files" {
		t.Errorf("Expected unit 'files', got '%s'", progress.Unit)
	}

	if progress.Details != "Uploading files..." {
		t.Errorf("Expected details 'Uploading files...', got '%s'", progress.Details)
	}

	if progress.Status != types.ProgressStatusInProgress {
		t.Errorf("Expected status 'in_progress', got '%s'", progress.Status)
	}
}

func TestHandleEventMessage(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test event message handling
	now := time.Now()
	eventParams := types.EventParams{
		Type:      "user_login",
		Title:     "User Login Event",
		Message:   "User authentication successful",
		Severity:  types.EventSeveritySuccess,
		Timestamp: &now,
		Data:      map[string]interface{}{"user_id": 123, "username": "john_doe"},
	}

	err = model.HandleEvent(eventParams, nil)
	if err != nil {
		t.Errorf("Expected no error handling event, got %v", err)
	}

	model.mu.RLock()
	defer model.mu.RUnlock()

	if len(model.events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(model.events))
	}

	event := model.events[0]
	if event.Type != "user_login" {
		t.Errorf("Expected type 'user_login', got '%s'", event.Type)
	}

	if event.Title != "User Login Event" {
		t.Errorf("Expected title 'User Login Event', got '%s'", event.Title)
	}

	if event.Severity != types.EventSeveritySuccess {
		t.Errorf("Expected severity 'success', got '%s'", event.Severity)
	}

	if event.Data == nil {
		t.Fatal("Expected event data to be set")
	}
}

func TestConcurrentAccess(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test concurrent message processing
	done := make(chan bool, 2)
	
	// Goroutine 1: Add log messages
	go func() {
		for i := 0; i < 100; i++ {
			now := time.Now()
			logParams := types.LogParams{
				Message:   "Test message",
				Level:     types.LogLevelInfo,
				Timestamp: &now,
				Component: "test-app",
			}
			_ = model.HandleLog(logParams, nil)
		}
		done <- true
	}()
	
	// Goroutine 2: Add metric messages
	go func() {
		for i := 0; i < 100; i++ {
			now := time.Now()
			metricParams := types.MetricParams{
				Name:      "test_metric",
				Value:     float64(i),
				Type:      types.MetricTypeGauge,
				Timestamp: &now,
			}
			_ = model.HandleMetric(metricParams, nil)
		}
		done <- true
	}()
	
	// Wait for both goroutines to complete
	<-done
	<-done
	
	// Verify state
	model.mu.RLock()
	defer model.mu.RUnlock()
	
	if len(model.logs) != 100 {
		t.Errorf("Expected 100 log entries, got %d", len(model.logs))
	}
	
	if len(model.metrics) != 1 {
		t.Errorf("Expected 1 metric (overwritten), got %d", len(model.metrics))
	}
	
	// The last metric value should be 99
	if metric, exists := model.metrics["test_metric"]; exists {
		if metric.Value != 99.0 {
			t.Errorf("Expected final metric value 99, got %f", metric.Value)
		}
	} else {
		t.Error("Expected metric 'test_metric' to exist")
	}
}

func TestViewModeString(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test different view modes
	testCases := []struct {
		mode     ViewMode
		expected string
	}{
		{ViewModeLogs, "Logs"},
		{ViewModeMetrics, "Metrics"},
		{ViewModeDashboard, "Dashboard"},
		{ViewModeConfig, "Config"},
		{ViewModeHelp, "Help"},
	}
	
	for _, tc := range testCases {
		model.viewMode = tc.mode
		result := model.getViewModeString()
		if result != tc.expected {
			t.Errorf("Expected view mode string '%s', got '%s'", tc.expected, result)
		}
	}
}

func TestStatusUpdate(t *testing.T) {
	config := Config{AppName: "test-app", Debug: false}
	model, err := NewModel(config)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	
	// Test that handling messages updates counts
	initialCount := model.messageCount
	
	now := time.Now()
	logParams := types.LogParams{
		Message:   "Test message",
		Level:     types.LogLevelInfo,
		Timestamp: &now,
		Component: "test-app",
	}
	
	_ = model.HandleLog(logParams, nil)
	
	// Note: messageCount is updated in handleDataUpdate, not in Handle* methods
	// So we can't test it directly here without triggering the full message flow
	
	// Instead, test that the data was stored
	model.mu.RLock()
	defer model.mu.RUnlock()
	
	if len(model.logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(model.logs))
	}
	
	// Test that last update time is reasonable
	if time.Since(model.lastUpdate) > time.Minute {
		t.Error("Last update time seems too old")
	}
	
	_ = initialCount // Suppress unused variable warning
}