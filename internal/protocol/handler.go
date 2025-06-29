package protocol

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hawk-tui/hawk-tui/pkg/types"
)

// MessageHandler defines the interface for handling different message types
type MessageHandler interface {
	HandleLog(params types.LogParams, msgID interface{}) error
	HandleMetric(params types.MetricParams, msgID interface{}) error
	HandleConfig(params types.ConfigParams, msgID interface{}) error
	HandleProgress(params types.ProgressParams, msgID interface{}) error
	HandleDashboard(params types.DashboardParams, msgID interface{}) error
	HandleEvent(params types.EventParams, msgID interface{}) error
}

// ResponseSender defines the interface for sending responses back to clients
type ResponseSender interface {
	SendResponse(id interface{}, result interface{}) error
	SendError(id interface{}, err *types.RPCError) error
	SendNotification(method string, params interface{}) error
}

// ProtocolHandler manages the JSON-RPC protocol communication
type ProtocolHandler struct {
	messageHandler  MessageHandler
	responseSender  ResponseSender
	reader          *bufio.Scanner
	writer          io.Writer
	metrics         HandlerMetrics
	config          HandlerConfig
	mu              sync.RWMutex
	running         bool
	stopCh          chan struct{}
	messageSequence int64
}

// HandlerConfig contains configuration options for the protocol handler
type HandlerConfig struct {
	MaxMessageSize    int           // Maximum size of a single message in bytes
	MaxBatchSize      int           // Maximum number of messages in a batch
	RateLimit         int           // Maximum messages per second
	BufferSize        int           // Size of the input buffer
	ValidationTimeout time.Duration // Timeout for message validation
	EnableMetrics     bool          // Whether to collect performance metrics
}

// HandlerMetrics contains performance and diagnostic metrics
type HandlerMetrics struct {
	MessagesReceived    int64
	MessagesProcessed   int64
	MessagesFailed      int64
	BatchesReceived     int64
	AverageProcessTime  time.Duration
	LastMessageTime     time.Time
	ErrorCount          map[string]int64
	MessageTypeCount    map[string]int64
	mu                  sync.RWMutex
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() HandlerConfig {
	return HandlerConfig{
		MaxMessageSize:    1024 * 1024, // 1MB
		MaxBatchSize:      100,
		RateLimit:         1000, // 1000 messages per second
		BufferSize:        64 * 1024, // 64KB
		ValidationTimeout: 100 * time.Millisecond,
		EnableMetrics:     true,
	}
}

// NewProtocolHandler creates a new protocol handler
func NewProtocolHandler(reader io.Reader, writer io.Writer, handler MessageHandler, sender ResponseSender) *ProtocolHandler {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // 64KB initial, 1MB max

	return &ProtocolHandler{
		messageHandler: handler,
		responseSender: sender,
		reader:         scanner,
		writer:         writer,
		config:         DefaultConfig(),
		stopCh:         make(chan struct{}),
		metrics: HandlerMetrics{
			ErrorCount:       make(map[string]int64),
			MessageTypeCount: make(map[string]int64),
		},
	}
}

// SetConfig updates the handler configuration
func (h *ProtocolHandler) SetConfig(config HandlerConfig) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.config = config
}

// Start begins processing messages from the reader
func (h *ProtocolHandler) Start() error {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return fmt.Errorf("handler is already running")
	}
	h.running = true
	h.mu.Unlock()

	go h.processMessages()
	return nil
}

// Stop gracefully stops the message processing
func (h *ProtocolHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return
	}

	h.running = false
	close(h.stopCh)
}

// GetMetrics returns a copy of the current metrics
func (h *ProtocolHandler) GetMetrics() HandlerMetrics {
	h.metrics.mu.RLock()
	defer h.metrics.mu.RUnlock()

	// Create a deep copy
	metrics := h.metrics
	metrics.ErrorCount = make(map[string]int64)
	metrics.MessageTypeCount = make(map[string]int64)

	for k, v := range h.metrics.ErrorCount {
		metrics.ErrorCount[k] = v
	}
	for k, v := range h.metrics.MessageTypeCount {
		metrics.MessageTypeCount[k] = v
	}

	return metrics
}

// processMessages is the main message processing loop
func (h *ProtocolHandler) processMessages() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Protocol handler panic: %v", r)
		}
	}()

	rateLimiter := time.NewTicker(time.Second / time.Duration(h.config.RateLimit))
	defer rateLimiter.Stop()

	for {
		select {
		case <-h.stopCh:
			return
		case <-rateLimiter.C:
			// Process one message per rate limit tick
			if h.reader.Scan() {
				h.processLine(h.reader.Text())
			} else {
				// Check for scanner error
				if err := h.reader.Err(); err != nil {
					log.Printf("Scanner error: %v", err)
					h.incrementErrorCount("scanner_error")
				}
				return
			}
		}
	}
}

// processLine processes a single line of input
func (h *ProtocolHandler) processLine(line string) {
	startTime := time.Now()
	h.incrementMetric("MessagesReceived")

	// Skip empty lines
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check message size limit
	if len(line) > h.config.MaxMessageSize {
		h.handleError(nil, types.NewInvalidRequestError("message too large"))
		return
	}

	// Update last message time
	h.metrics.mu.Lock()
	h.metrics.LastMessageTime = time.Now()
	h.metrics.mu.Unlock()

	// Try to parse as a single message first
	if strings.HasPrefix(line, "{") {
		h.processSingleMessage([]byte(line))
	} else if strings.HasPrefix(line, "[") {
		h.processBatchMessage([]byte(line))
	} else {
		h.handleError(nil, types.NewParseError("invalid JSON format"))
		return
	}

	// Update processing time metrics
	if h.config.EnableMetrics {
		processingTime := time.Since(startTime)
		h.updateAverageProcessTime(processingTime)
	}
}

// processSingleMessage processes a single JSON-RPC message
func (h *ProtocolHandler) processSingleMessage(data []byte) {
	msg, err := types.ParseMessage(data)
	if err != nil {
		h.handleError(nil, types.NewParseError(err.Error()))
		return
	}

	h.handleMessage(msg)
}

// processBatchMessage processes a batch of JSON-RPC messages
func (h *ProtocolHandler) processBatchMessage(data []byte) {
	batch, err := types.ParseBatch(data)
	if err != nil {
		h.handleError(nil, types.NewParseError(err.Error()))
		return
	}

	// Check batch size limit
	if len(batch) > h.config.MaxBatchSize {
		h.handleError(nil, types.NewInvalidRequestError("batch too large"))
		return
	}

	h.incrementMetric("BatchesReceived")

	// Process each message in the batch
	for _, msg := range batch {
		h.handleMessage(&msg)
	}
}

// handleMessage processes a single parsed JSON-RPC message
func (h *ProtocolHandler) handleMessage(msg *types.JSONRPCMessage) {
	// Validate the message
	if err := msg.Validate(); err != nil {
		h.handleError(msg.ID, err.(*types.RPCError))
		return
	}

	// Track message type
	if h.config.EnableMetrics {
		h.incrementMessageTypeCount(msg.Method)
	}

	// Route the message based on method
	switch msg.Method {
	case "hawk.log":
		h.handleLogMessage(msg)
	case "hawk.metric":
		h.handleMetricMessage(msg)
	case "hawk.config":
		h.handleConfigMessage(msg)
	case "hawk.progress":
		h.handleProgressMessage(msg)
	case "hawk.dashboard":
		h.handleDashboardMessage(msg)
	case "hawk.event":
		h.handleEventMessage(msg)
	default:
		h.handleError(msg.ID, types.NewMethodNotFoundError(msg.Method))
	}
}

// handleLogMessage processes a log message
func (h *ProtocolHandler) handleLogMessage(msg *types.JSONRPCMessage) {
	var params types.LogParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Set default values
	if params.Level == "" {
		params.Level = types.LogLevelInfo
	}
	if params.Timestamp == nil {
		now := time.Now()
		params.Timestamp = &now
	}

	// Call the handler
	if err := h.messageHandler.HandleLog(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleMetricMessage processes a metric message
func (h *ProtocolHandler) handleMetricMessage(msg *types.JSONRPCMessage) {
	var params types.MetricParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Set default values
	if params.Type == "" {
		params.Type = types.MetricTypeGauge
	}
	if params.Timestamp == nil {
		now := time.Now()
		params.Timestamp = &now
	}

	// Call the handler
	if err := h.messageHandler.HandleMetric(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleConfigMessage processes a config message
func (h *ProtocolHandler) handleConfigMessage(msg *types.JSONRPCMessage) {
	var params types.ConfigParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Validate required fields
	if params.Key == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("key is required"))
		return
	}

	// Set default values
	if params.Type == "" {
		params.Type = types.ConfigTypeString
	}

	// Call the handler
	if err := h.messageHandler.HandleConfig(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleProgressMessage processes a progress message
func (h *ProtocolHandler) handleProgressMessage(msg *types.JSONRPCMessage) {
	var params types.ProgressParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Validate required fields
	if params.ID == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("id is required"))
		return
	}
	if params.Label == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("label is required"))
		return
	}

	// Set default values
	if params.Status == "" {
		params.Status = types.ProgressStatusInProgress
	}

	// Call the handler
	if err := h.messageHandler.HandleProgress(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleDashboardMessage processes a dashboard message
func (h *ProtocolHandler) handleDashboardMessage(msg *types.JSONRPCMessage) {
	var params types.DashboardParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Validate required fields
	if params.WidgetID == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("widget_id is required"))
		return
	}
	if params.Type == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("type is required"))
		return
	}

	// Call the handler
	if err := h.messageHandler.HandleDashboard(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleEventMessage processes an event message
func (h *ProtocolHandler) handleEventMessage(msg *types.JSONRPCMessage) {
	var params types.EventParams
	if err := msg.GetParamsAs(&params); err != nil {
		h.handleError(msg.ID, types.NewInvalidParamsError(err.Error()))
		return
	}

	// Validate required fields
	if params.Type == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("type is required"))
		return
	}
	if params.Title == "" {
		h.handleError(msg.ID, types.NewInvalidParamsError("title is required"))
		return
	}

	// Set default values
	if params.Severity == "" {
		params.Severity = types.EventSeverityInfo
	}
	if params.Timestamp == nil {
		now := time.Now()
		params.Timestamp = &now
	}

	// Call the handler
	if err := h.messageHandler.HandleEvent(params, msg.ID); err != nil {
		h.handleError(msg.ID, types.NewInternalError(err.Error()))
		return
	}

	h.incrementMetric("MessagesProcessed")
}

// handleError handles errors by sending error responses and updating metrics
func (h *ProtocolHandler) handleError(id interface{}, rpcErr *types.RPCError) {
	h.incrementMetric("MessagesFailed")
	h.incrementErrorCount(rpcErr.Message)

	// Send error response if this is a request (has ID)
	if id != nil && h.responseSender != nil {
		if err := h.responseSender.SendError(id, rpcErr); err != nil {
			log.Printf("Failed to send error response: %v", err)
		}
	}

	// Log the error
	log.Printf("Protocol error: %s (code: %d, data: %v)", rpcErr.Message, rpcErr.Code, rpcErr.Data)
}

// Metrics helper functions
func (h *ProtocolHandler) incrementMetric(metric string) {
	if !h.config.EnableMetrics {
		return
	}

	h.metrics.mu.Lock()
	defer h.metrics.mu.Unlock()

	switch metric {
	case "MessagesReceived":
		h.metrics.MessagesReceived++
	case "MessagesProcessed":
		h.metrics.MessagesProcessed++
	case "MessagesFailed":
		h.metrics.MessagesFailed++
	case "BatchesReceived":
		h.metrics.BatchesReceived++
	}
}

func (h *ProtocolHandler) incrementErrorCount(errorType string) {
	if !h.config.EnableMetrics {
		return
	}

	h.metrics.mu.Lock()
	defer h.metrics.mu.Unlock()
	h.metrics.ErrorCount[errorType]++
}

func (h *ProtocolHandler) incrementMessageTypeCount(messageType string) {
	if !h.config.EnableMetrics {
		return
	}

	h.metrics.mu.Lock()
	defer h.metrics.mu.Unlock()
	h.metrics.MessageTypeCount[messageType]++
}

func (h *ProtocolHandler) updateAverageProcessTime(duration time.Duration) {
	h.metrics.mu.Lock()
	defer h.metrics.mu.Unlock()

	// Simple exponential moving average
	if h.metrics.AverageProcessTime == 0 {
		h.metrics.AverageProcessTime = duration
	} else {
		// EMA with alpha = 0.1
		h.metrics.AverageProcessTime = time.Duration(
			0.9*float64(h.metrics.AverageProcessTime) + 0.1*float64(duration),
		)
	}
}

// SendConfigUpdate sends a configuration update request to the client
func (h *ProtocolHandler) SendConfigUpdate(key string, value interface{}) error {
	if h.responseSender == nil {
		return fmt.Errorf("no response sender configured")
	}

	params := types.ConfigUpdateParams{
		Key:   key,
		Value: value,
	}

	return h.responseSender.SendNotification("hawk.config_update", params)
}

// SendExecuteCommand sends a command execution request to the client
func (h *ProtocolHandler) SendExecuteCommand(command string, args map[string]interface{}) error {
	if h.responseSender == nil {
		return fmt.Errorf("no response sender configured")
	}

	params := types.ExecuteParams{
		Command: command,
		Args:    args,
	}

	h.messageSequence++
	return h.responseSender.SendNotification("hawk.execute", params)
}

// SendDataRequest sends a data request to the client
func (h *ProtocolHandler) SendDataRequest(requestType types.RequestType, filter map[string]interface{}, timeRange string) error {
	if h.responseSender == nil {
		return fmt.Errorf("no response sender configured")
	}

	params := types.RequestParams{
		Type:      requestType,
		Filter:    filter,
		TimeRange: timeRange,
	}

	h.messageSequence++
	return h.responseSender.SendNotification("hawk.request", params)
}

// IsRunning returns whether the handler is currently running
func (h *ProtocolHandler) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.running
}