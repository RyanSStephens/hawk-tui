package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hawk-tui/hawk-tui/pkg/types"
)

func main() {
	fmt.Fprintln(os.Stderr, "Starting Hawk TUI Demo Application...")
	fmt.Fprintln(os.Stderr, "This application will send sample data to stdout in Hawk TUI format.")
	fmt.Fprintln(os.Stderr, "To view in the TUI, run: go run examples/go/simple_demo.go | hawk")
	fmt.Fprintln(os.Stderr, "")

	// Send initial configuration
	sendConfig("app.name", "Demo Application", "Application name")
	sendConfig("app.version", "1.0.0", "Application version")
	sendConfig("log.level", "INFO", "Logging level")
	sendConfig("server.port", 8080, "Server port number")
	sendConfig("server.timeout", 30.0, "Request timeout in seconds")
	
	// Send initial log messages
	sendLog("Application starting up...", "INFO", map[string]interface{}{
		"version": "1.0.0",
		"pid":     os.Getpid(),
	})
	
	sendLog("Configuration loaded", "SUCCESS", nil)
	sendLog("Database connection established", "INFO", map[string]interface{}{
		"host": "localhost",
		"port": 5432,
	})

	// Create some dashboard widgets
	sendDashboard("server_status", "status_grid", "Server Status", map[string]interface{}{
		"API Server":    map[string]interface{}{"status": "healthy", "response_time": "12ms"},
		"Database":      map[string]interface{}{"status": "healthy", "response_time": "3ms"},
		"Cache":         map[string]interface{}{"status": "healthy", "response_time": "1ms"},
		"Load Balancer": map[string]interface{}{"status": "degraded", "response_time": "45ms"},
	})

	sendDashboard("system_info", "text", "System Information", map[string]interface{}{
		"content": "OS: Linux\nArch: x86_64\nMemory: 16GB\nCPU: 8 cores\nUptime: 2d 5h 30m",
		"format":  "plain",
	})

	// Send progress for startup tasks
	progressID := "startup_tasks"
	for i := 0; i <= 100; i += 10 {
		sendProgress(progressID, "Initializing application", float64(i), 100, "%", "in_progress", fmt.Sprintf("Step %d/10", i/10+1))
		time.Sleep(200 * time.Millisecond)
	}
	sendProgress(progressID, "Initialization complete", 100, 100, "%", "completed", "All systems ready")

	// Send an event
	sendEvent("deployment_completed", "Application Deployed", "Version 1.0.0 deployed successfully", "success", map[string]interface{}{
		"version":  "1.0.0",
		"duration": "2m 34s",
		"services": []string{"api", "worker", "scheduler"},
	})

	// Simulate real-time metrics
	fmt.Fprintln(os.Stderr, "Sending real-time metrics... (Press Ctrl+C to stop)")
	
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(startTime)
			
			// System metrics
			sendMetric("system.cpu.usage", 15+rand.Float64()*70, "gauge", "%")
			sendMetric("system.memory.usage", 30+rand.Float64()*40, "gauge", "%")
			sendMetric("system.disk.usage", 25+rand.Float64()*20, "gauge", "%")
			sendMetric("system.network.rx", rand.Float64()*1000, "gauge", "KB/s")
			sendMetric("system.network.tx", rand.Float64()*500, "gauge", "KB/s")
			
			// Application metrics
			sendMetric("app.requests.total", float64(time.Now().Unix()%10000), "counter", "")
			sendMetric("app.requests.per_second", 50+rand.Float64()*200, "gauge", "req/s")
			sendMetric("app.response.time.avg", 10+rand.Float64()*100, "gauge", "ms")
			sendMetric("app.errors.rate", rand.Float64()*5, "gauge", "%")
			sendMetric("app.active.connections", float64(50+rand.Intn(200)), "gauge", "")
			
			// Database metrics
			sendMetric("db.connections.active", float64(10+rand.Intn(90)), "gauge", "")
			sendMetric("db.queries.per_second", 20+rand.Float64()*80, "gauge", "q/s")
			sendMetric("db.query.time.avg", 5+rand.Float64()*50, "gauge", "ms")
			
			// Occasionally send log messages
			if elapsed.Seconds() > 0 && int(elapsed.Seconds())%5 == 0 {
				switch rand.Intn(4) {
				case 0:
					sendLog("Request processed successfully", "INFO", map[string]interface{}{
						"method":      "GET",
						"path":        "/api/users",
						"status":      200,
						"duration_ms": rand.Intn(100),
					})
				case 1:
					sendLog("Cache miss for key", "WARN", map[string]interface{}{
						"key":    fmt.Sprintf("user_%d", rand.Intn(1000)),
						"reason": "expired",
					})
				case 2:
					sendLog("Database query executed", "DEBUG", map[string]interface{}{
						"query":      "SELECT * FROM users WHERE active = true",
						"duration":   fmt.Sprintf("%dms", rand.Intn(50)),
						"rows":       rand.Intn(100),
					})
				case 3:
					if rand.Float64() < 0.1 { // 10% chance of error
						sendLog("Request failed with error", "ERROR", map[string]interface{}{
							"error":  "Connection timeout",
							"method": "POST",
							"path":   "/api/upload",
							"status": 500,
						})
					}
				}
			}
			
			// Occasionally send events
			if elapsed.Seconds() > 0 && int(elapsed.Seconds())%30 == 0 {
				events := []struct {
					eventType string
					title     string
					message   string
					severity  string
				}{
					{"health_check", "Health Check", "All services are healthy", "info"},
					{"backup_completed", "Backup Complete", "Daily backup completed successfully", "success"},
					{"high_cpu_usage", "High CPU Usage", "CPU usage is above 80%", "warning"},
				}
				
				event := events[rand.Intn(len(events))]
				sendEvent(event.eventType, event.title, event.message, event.severity, map[string]interface{}{
					"timestamp": time.Now().Format(time.RFC3339),
					"host":      "app-server-01",
				})
			}
		}
	}
}

func sendLog(message, level string, context map[string]interface{}) {
	params := types.LogParams{
		Message: message,
		Level:   types.LogLevel(level),
		Context: context,
	}
	
	now := time.Now()
	params.Timestamp = &now
	
	msg := types.NewLogMessage(params, nil)
	sendMessage(msg)
}

func sendMetric(name string, value float64, metricType, unit string) {
	params := types.MetricParams{
		Name:  name,
		Value: value,
		Type:  types.MetricType(metricType),
		Unit:  unit,
	}
	
	now := time.Now()
	params.Timestamp = &now
	
	msg := types.NewMetricMessage(params, nil)
	sendMessage(msg)
}

func sendConfig(key string, value interface{}, description string) {
	params := types.ConfigParams{
		Key:         key,
		Value:       value,
		Description: description,
	}
	
	// Set type based on value
	switch value.(type) {
	case string:
		params.Type = types.ConfigTypeString
	case int, int64:
		params.Type = types.ConfigTypeInteger
	case float64:
		params.Type = types.ConfigTypeFloat
	case bool:
		params.Type = types.ConfigTypeBoolean
	}
	
	msg := types.NewConfigMessage(params, nil)
	sendMessage(msg)
}

func sendProgress(id, label string, current, total float64, unit, status, details string) {
	params := types.ProgressParams{
		ID:      id,
		Label:   label,
		Current: current,
		Total:   total,
		Unit:    unit,
		Status:  types.ProgressStatus(status),
		Details: details,
	}
	
	msg := types.NewProgressMessage(params, nil)
	sendMessage(msg)
}

func sendDashboard(widgetID, widgetType, title string, data interface{}) {
	params := types.DashboardParams{
		WidgetID: widgetID,
		Type:     types.WidgetType(widgetType),
		Title:    title,
		Data:     data,
	}
	
	msg := types.NewDashboardMessage(params, nil)
	sendMessage(msg)
}

func sendEvent(eventType, title, message, severity string, data map[string]interface{}) {
	params := types.EventParams{
		Type:     eventType,
		Title:    title,
		Message:  message,
		Severity: types.EventSeverity(severity),
		Data:     data,
	}
	
	now := time.Now()
	params.Timestamp = &now
	
	msg := types.NewEventMessage(params, nil)
	sendMessage(msg)
}

func sendMessage(msg *types.JSONRPCMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling message: %v\n", err)
		return
	}
	
	fmt.Println(string(data))
}