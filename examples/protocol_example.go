package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hawk-tui/hawk-tui/internal/protocol"
	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// ExampleMessageHandler implements the MessageHandler interface for demonstration
type ExampleMessageHandler struct{}

func (h *ExampleMessageHandler) HandleLog(params types.LogParams, msgID interface{}) error {
	fmt.Printf("[%s] %s: %s\n",
		params.Timestamp.Format("15:04:05"),
		params.Level,
		params.Message)
	if params.Context != nil {
		fmt.Printf("  Context: %+v\n", params.Context)
	}
	return nil
}

func (h *ExampleMessageHandler) HandleMetric(params types.MetricParams, msgID interface{}) error {
	fmt.Printf("METRIC [%s] %s = %.2f %s (%s)\n",
		params.Timestamp.Format("15:04:05"),
		params.Name,
		params.Value,
		params.Unit,
		params.Type)
	if params.Tags != nil {
		fmt.Printf("  Tags: %+v\n", params.Tags)
	}
	return nil
}

func (h *ExampleMessageHandler) HandleConfig(params types.ConfigParams, msgID interface{}) error {
	fmt.Printf("CONFIG %s = %v (%s)\n", params.Key, params.Value, params.Type)
	if params.Description != "" {
		fmt.Printf("  Description: %s\n", params.Description)
	}
	return nil
}

func (h *ExampleMessageHandler) HandleProgress(params types.ProgressParams, msgID interface{}) error {
	percentage := (params.Current / params.Total) * 100
	fmt.Printf("PROGRESS [%s] %s: %.1f%% (%.1f/%.1f %s)\n",
		params.ID,
		params.Label,
		percentage,
		params.Current,
		params.Total,
		params.Unit)
	return nil
}

func (h *ExampleMessageHandler) HandleDashboard(params types.DashboardParams, msgID interface{}) error {
	fmt.Printf("DASHBOARD [%s] %s (%s)\n", params.WidgetID, params.Title, params.Type)
	if params.Data != nil {
		dataJSON, _ := json.MarshalIndent(params.Data, "  ", "  ")
		fmt.Printf("  Data: %s\n", string(dataJSON))
	}
	return nil
}

func (h *ExampleMessageHandler) HandleEvent(params types.EventParams, msgID interface{}) error {
	fmt.Printf("EVENT [%s] %s: %s (%s)\n",
		params.Timestamp.Format("15:04:05"),
		params.Type,
		params.Title,
		params.Severity)
	if params.Message != "" {
		fmt.Printf("  Message: %s\n", params.Message)
	}
	return nil
}

// ExampleResponseSender implements the ResponseSender interface for demonstration
type ExampleResponseSender struct{}

func (s *ExampleResponseSender) SendResponse(id interface{}, result interface{}) error {
	response := types.JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(response)
	fmt.Printf("RESPONSE: %s\n", string(data))
	return nil
}

func (s *ExampleResponseSender) SendError(id interface{}, err *types.RPCError) error {
	response := types.JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}
	data, _ := json.Marshal(response)
	fmt.Printf("ERROR RESPONSE: %s\n", string(data))
	return nil
}

func (s *ExampleResponseSender) SendNotification(method string, params interface{}) error {
	notification := types.JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	data, _ := json.Marshal(notification)
	fmt.Printf("NOTIFICATION: %s\n", string(data))
	return nil
}

func main() {
	fmt.Println("Hawk TUI Protocol Example")
	fmt.Println("=========================")

	// Create example input messages
	exampleMessages := []string{
		// Log message
		`{"jsonrpc": "2.0", "method": "hawk.log", "params": {"message": "Server started successfully", "level": "SUCCESS", "component": "server"}, "id": 1}`,

		// Metric message
		`{"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "requests_per_second", "value": 145.7, "type": "gauge", "unit": "req/s", "tags": {"endpoint": "/api/users"}}}`,

		// Config message
		`{"jsonrpc": "2.0", "method": "hawk.config", "params": {"key": "server.port", "value": 8080, "type": "integer", "description": "HTTP server port", "min": 1, "max": 65535}}`,

		// Progress message
		`{"jsonrpc": "2.0", "method": "hawk.progress", "params": {"id": "upload_001", "label": "Uploading file", "current": 75, "total": 100, "unit": "%", "status": "in_progress"}}`,

		// Dashboard message with status grid
		`{"jsonrpc": "2.0", "method": "hawk.dashboard", "params": {"widget_id": "services", "type": "status_grid", "title": "Service Status", "data": {"Database": {"status": "healthy", "response_time": "12ms"}, "Redis": {"status": "healthy", "response_time": "2ms"}}}}`,

		// Event message
		`{"jsonrpc": "2.0", "method": "hawk.event", "params": {"type": "deployment", "title": "Deployment Completed", "message": "Version 2.1.0 deployed successfully", "severity": "success"}}`,

		// Batch message
		`[
			{"jsonrpc": "2.0", "method": "hawk.log", "params": {"message": "Processing batch"}},
			{"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "batch_size", "value": 100}},
			{"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "processing_time", "value": 0.125, "unit": "seconds"}}
		]`,

		// Invalid message (will generate error)
		`{"jsonrpc": "2.0", "method": "hawk.invalid", "params": {}}`,
	}

	// Create input reader from example messages
	input := strings.NewReader(strings.Join(exampleMessages, "\n"))

	// Create output buffer (normally this would be stdout)
	output := &bytes.Buffer{}

	// Create handler and sender
	messageHandler := &ExampleMessageHandler{}
	responseSender := &ExampleResponseSender{}

	// Create protocol handler
	handler := protocol.NewProtocolHandler(input, output, messageHandler, responseSender)

	// Configure the handler
	config := protocol.DefaultConfig()
	config.EnableMetrics = true
	config.RateLimit = 100 // Lower rate limit for demo
	handler.SetConfig(config)

	// Start processing
	fmt.Println("\nProcessing messages:")
	fmt.Println("--------------------")

	if err := handler.Start(); err != nil {
		log.Fatalf("Failed to start handler: %v", err)
	}

	// Wait a bit for processing to complete
	time.Sleep(100 * time.Millisecond)

	// Stop the handler
	handler.Stop()

	// Show metrics
	fmt.Println("\nHandler Metrics:")
	fmt.Println("----------------")
	metrics := handler.GetMetrics()
	fmt.Printf("Messages Received: %d\n", metrics.MessagesReceived)
	fmt.Printf("Messages Processed: %d\n", metrics.MessagesProcessed)
	fmt.Printf("Messages Failed: %d\n", metrics.MessagesFailed)
	fmt.Printf("Batches Received: %d\n", metrics.BatchesReceived)
	fmt.Printf("Average Process Time: %v\n", metrics.AverageProcessTime)
	fmt.Printf("Last Message Time: %v\n", metrics.LastMessageTime.Format("15:04:05"))

	if len(metrics.ErrorCount) > 0 {
		fmt.Println("\nError Counts:")
		for errorType, count := range metrics.ErrorCount {
			fmt.Printf("  %s: %d\n", errorType, count)
		}
	}

	if len(metrics.MessageTypeCount) > 0 {
		fmt.Println("\nMessage Type Counts:")
		for msgType, count := range metrics.MessageTypeCount {
			fmt.Printf("  %s: %d\n", msgType, count)
		}
	}

	// Demonstrate sending messages from TUI back to client
	fmt.Println("\nTUI → Client Communication:")
	fmt.Println("---------------------------")

	// Send configuration update
	handler.SendConfigUpdate("log_level", "DEBUG")

	// Send command execution request
	handler.SendExecuteCommand("restart_workers", map[string]interface{}{
		"force":   true,
		"timeout": 30,
	})

	// Send data request
	handler.SendDataRequest(types.RequestTypeMetrics, map[string]interface{}{
		"component": "database",
	}, "last_5_minutes")

	fmt.Println("\nExample completed!")
}

// Helper function to create example message structs
func createExampleMessages() {
	fmt.Println("\nExample message creation using helper functions:")
	fmt.Println("------------------------------------------------")

	// Create log message
	logMsg := types.NewLogMessage(types.LogParams{
		Message: "Application started",
		Level:   types.LogLevelInfo,
		Context: map[string]interface{}{
			"version": "1.0.0",
			"port":    8080,
		},
		Tags: []string{"startup", "server"},
	}, "log_001")

	logJSON, _ := json.MarshalIndent(logMsg, "", "  ")
	fmt.Printf("Log Message:\n%s\n\n", string(logJSON))

	// Create metric message
	metricMsg := types.NewMetricMessage(types.MetricParams{
		Name:  "cpu_usage",
		Value: 65.4,
		Type:  types.MetricTypeGauge,
		Unit:  "%",
		Tags: map[string]interface{}{
			"host": "server-01",
			"core": "all",
		},
	}, "metric_001")

	metricJSON, _ := json.MarshalIndent(metricMsg, "", "  ")
	fmt.Printf("Metric Message:\n%s\n\n", string(metricJSON))

	// Create dashboard message with chart data
	chartData := types.ChartData{
		Series: []types.ChartSeries{
			{
				Name:  "Response Time",
				Color: "#007acc",
				Data: []types.ChartPoint{
					{X: time.Now().Add(-4 * time.Minute), Y: 120.5},
					{X: time.Now().Add(-3 * time.Minute), Y: 115.2},
					{X: time.Now().Add(-2 * time.Minute), Y: 130.8},
					{X: time.Now().Add(-1 * time.Minute), Y: 125.1},
					{X: time.Now(), Y: 118.9},
				},
			},
		},
	}

	dashboardMsg := types.NewDashboardMessage(types.DashboardParams{
		WidgetID: "response_times",
		Type:     types.WidgetTypeMetricChart,
		Title:    "API Response Times",
		Data:     chartData,
		Layout: &types.WidgetLayout{
			Row:    0,
			Col:    0,
			Width:  8,
			Height: 4,
		},
	}, "dashboard_001")

	dashboardJSON, _ := json.MarshalIndent(dashboardMsg, "", "  ")
	fmt.Printf("Dashboard Message:\n%s\n", string(dashboardJSON))
}

func init() {
	// Set up logging
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
