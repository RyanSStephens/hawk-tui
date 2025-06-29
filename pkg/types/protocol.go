package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONRPCMessage represents a JSON-RPC 2.0 message
type JSONRPCMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error object
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements the error interface
func (e *RPCError) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// HawkMetadata contains Hawk-specific message metadata
type HawkMetadata struct {
	AppName   string `json:"app_name,omitempty"`
	Component string `json:"component,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Sequence  int64  `json:"sequence,omitempty"`
}

// MessageBatch represents a batch of JSON-RPC messages
type MessageBatch []JSONRPCMessage

// LogLevel represents the severity level of a log message
type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarn    LogLevel = "WARN"
	LogLevelError   LogLevel = "ERROR"
	LogLevelSuccess LogLevel = "SUCCESS"
)

// LogParams represents parameters for hawk.log method
type LogParams struct {
	Message   string                 `json:"message"`
	Level     LogLevel               `json:"level,omitempty"`
	Timestamp *time.Time             `json:"timestamp,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Component string                 `json:"component,omitempty"`
}

// MetricType represents the type of metric being reported
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

// MetricParams represents parameters for hawk.metric method
type MetricParams struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Type      MetricType             `json:"type,omitempty"`
	Timestamp *time.Time             `json:"timestamp,omitempty"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
	Unit      string                 `json:"unit,omitempty"`
}

// ConfigType represents the type of configuration value
type ConfigType string

const (
	ConfigTypeString  ConfigType = "string"
	ConfigTypeInteger ConfigType = "integer"
	ConfigTypeFloat   ConfigType = "float"
	ConfigTypeBoolean ConfigType = "boolean"
	ConfigTypeEnum    ConfigType = "enum"
)

// ConfigParams represents parameters for hawk.config method
type ConfigParams struct {
	Key              string        `json:"key"`
	Value            interface{}   `json:"value,omitempty"`
	Type             ConfigType    `json:"type,omitempty"`
	Description      string        `json:"description,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Min              *float64      `json:"min,omitempty"`
	Max              *float64      `json:"max,omitempty"`
	Options          []interface{} `json:"options,omitempty"`
	RestartRequired  bool          `json:"restart_required,omitempty"`
	Category         string        `json:"category,omitempty"`
}

// ProgressStatus represents the status of a progress operation
type ProgressStatus string

const (
	ProgressStatusPending    ProgressStatus = "pending"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
	ProgressStatusError      ProgressStatus = "error"
)

// ProgressParams represents parameters for hawk.progress method
type ProgressParams struct {
	ID                   string         `json:"id"`
	Label                string         `json:"label"`
	Current              float64        `json:"current"`
	Total                float64        `json:"total"`
	Unit                 string         `json:"unit,omitempty"`
	Status               ProgressStatus `json:"status,omitempty"`
	Details              string         `json:"details,omitempty"`
	EstimatedCompletion  *time.Time     `json:"estimated_completion,omitempty"`
}

// WidgetType represents the type of dashboard widget
type WidgetType string

const (
	WidgetTypeStatusGrid  WidgetType = "status_grid"
	WidgetTypeMetricChart WidgetType = "metric_chart"
	WidgetTypeTable       WidgetType = "table"
	WidgetTypeText        WidgetType = "text"
	WidgetTypeGauge       WidgetType = "gauge"
	WidgetTypeHistogram   WidgetType = "histogram"
)

// WidgetLayout represents the layout configuration for a widget
type WidgetLayout struct {
	Row    int `json:"row"`
	Col    int `json:"col"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DashboardParams represents parameters for hawk.dashboard method
type DashboardParams struct {
	WidgetID string                 `json:"widget_id"`
	Type     WidgetType             `json:"type"`
	Title    string                 `json:"title,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	Layout   *WidgetLayout          `json:"layout,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

// EventSeverity represents the severity level of an event
type EventSeverity string

const (
	EventSeverityInfo     EventSeverity = "info"
	EventSeverityWarning  EventSeverity = "warning"
	EventSeverityError    EventSeverity = "error"
	EventSeverityCritical EventSeverity = "critical"
	EventSeveritySuccess  EventSeverity = "success"
)

// EventParams represents parameters for hawk.event method
type EventParams struct {
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message,omitempty"`
	Severity  EventSeverity          `json:"severity,omitempty"`
	Timestamp *time.Time             `json:"timestamp,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// ConfigUpdateParams represents parameters for hawk.config_update method (TUI → Client)
type ConfigUpdateParams struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ExecuteParams represents parameters for hawk.execute method (TUI → Client)
type ExecuteParams struct {
	Command string                 `json:"command"`
	Args    map[string]interface{} `json:"args,omitempty"`
}

// RequestType represents the type of data being requested
type RequestType string

const (
	RequestTypeMetrics RequestType = "metrics"
	RequestTypeLogs    RequestType = "logs"
	RequestTypeConfig  RequestType = "config"
	RequestTypeStatus  RequestType = "status"
)

// RequestParams represents parameters for hawk.request method (TUI → Client)
type RequestParams struct {
	Type      RequestType            `json:"type"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	TimeRange string                 `json:"timerange,omitempty"`
	Limit     int                    `json:"limit,omitempty"`
}

// StatusGridData represents data for a status grid widget
type StatusGridData map[string]StatusGridItem

// StatusGridItem represents a single item in a status grid
type StatusGridItem struct {
	Status       string      `json:"status"`
	ResponseTime string      `json:"response_time,omitempty"`
	Details      string      `json:"details,omitempty"`
	LastChecked  *time.Time  `json:"last_checked,omitempty"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// TableData represents data for a table widget
type TableData struct {
	Headers []string                   `json:"headers"`
	Rows    [][]interface{}            `json:"rows"`
	Config  map[string]interface{}     `json:"config,omitempty"`
}

// ChartData represents data for chart widgets
type ChartData struct {
	Series []ChartSeries          `json:"series"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// ChartSeries represents a single data series in a chart
type ChartSeries struct {
	Name   string      `json:"name"`
	Data   []ChartPoint `json:"data"`
	Color  string      `json:"color,omitempty"`
	Type   string      `json:"type,omitempty"` // line, bar, area, etc.
}

// ChartPoint represents a single data point in a chart series
type ChartPoint struct {
	X interface{} `json:"x"` // timestamp or category
	Y float64     `json:"y"` // value
}

// GaugeData represents data for gauge widgets
type GaugeData struct {
	Value    float64 `json:"value"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Unit     string  `json:"unit,omitempty"`
	Zones    []GaugeZone `json:"zones,omitempty"`
}

// GaugeZone represents a colored zone in a gauge
type GaugeZone struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Color string  `json:"color"`
	Label string  `json:"label,omitempty"`
}

// HistogramData represents data for histogram widgets
type HistogramData struct {
	Buckets []HistogramBucket      `json:"buckets"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// HistogramBucket represents a single bucket in a histogram
type HistogramBucket struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

// TextData represents data for text widgets
type TextData struct {
	Content string `json:"content"`
	Format  string `json:"format,omitempty"` // plain, markdown, html
}

// Standard JSON-RPC error codes
const (
	ParseErrorCode     = -32700
	InvalidRequestCode = -32600
	MethodNotFoundCode = -32601
	InvalidParamsCode  = -32602
	InternalErrorCode  = -32603

	// Custom Hawk error codes
	HawkInvalidMessageType = -32000
	HawkInvalidData        = -32001
	HawkResourceLimit      = -32002
	HawkAuthenticationError = -32003
)

// Helper functions for creating standard error responses
func NewParseError(data interface{}) *RPCError {
	return &RPCError{
		Code:    ParseErrorCode,
		Message: "Parse error",
		Data:    data,
	}
}

func NewInvalidRequestError(data interface{}) *RPCError {
	return &RPCError{
		Code:    InvalidRequestCode,
		Message: "Invalid Request",
		Data:    data,
	}
}

func NewMethodNotFoundError(method string) *RPCError {
	return &RPCError{
		Code:    MethodNotFoundCode,
		Message: "Method not found",
		Data:    map[string]interface{}{"method": method},
	}
}

func NewInvalidParamsError(data interface{}) *RPCError {
	return &RPCError{
		Code:    InvalidParamsCode,
		Message: "Invalid params",
		Data:    data,
	}
}

func NewInternalError(data interface{}) *RPCError {
	return &RPCError{
		Code:    InternalErrorCode,
		Message: "Internal error",
		Data:    data,
	}
}

// Helper functions for creating JSON-RPC messages
func NewLogMessage(params LogParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.log",
		Params:  params,
		ID:      id,
	}
}

func NewMetricMessage(params MetricParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.metric",
		Params:  params,
		ID:      id,
	}
}

func NewConfigMessage(params ConfigParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.config",
		Params:  params,
		ID:      id,
	}
}

func NewProgressMessage(params ProgressParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.progress",
		Params:  params,
		ID:      id,
	}
}

func NewDashboardMessage(params DashboardParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.dashboard",
		Params:  params,
		ID:      id,
	}
}

func NewEventMessage(params EventParams, id interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		JSONRPC: "2.0",
		Method:  "hawk.event",
		Params:  params,
		ID:      id,
	}
}

// ParseMessage parses a JSON-RPC message from bytes
func ParseMessage(data []byte) (*JSONRPCMessage, error) {
	var msg JSONRPCMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// ParseBatch parses a batch of JSON-RPC messages from bytes
func ParseBatch(data []byte) (MessageBatch, error) {
	var batch MessageBatch
	err := json.Unmarshal(data, &batch)
	return batch, err
}

// IsRequest returns true if the message is a request
func (m *JSONRPCMessage) IsRequest() bool {
	return m.Method != "" && m.Error == nil && m.Result == nil
}

// IsResponse returns true if the message is a response
func (m *JSONRPCMessage) IsResponse() bool {
	return m.Method == "" && (m.Error != nil || m.Result != nil)
}

// IsNotification returns true if the message is a notification (request without ID)
func (m *JSONRPCMessage) IsNotification() bool {
	return m.Method != "" && m.ID == nil && m.Error == nil && m.Result == nil
}

// Validate performs basic validation on the JSON-RPC message
func (m *JSONRPCMessage) Validate() error {
	if m.JSONRPC != "2.0" {
		return NewInvalidRequestError("jsonrpc must be '2.0'")
	}

	if m.IsRequest() && m.Method == "" {
		return NewInvalidRequestError("method is required for requests")
	}

	if m.IsResponse() && m.ID == nil {
		return NewInvalidRequestError("id is required for responses")
	}

	return nil
}

// GetParamsAs unmarshals the params field into the provided struct
func (m *JSONRPCMessage) GetParamsAs(v interface{}) error {
	if m.Params == nil {
		return nil
	}

	// Convert to JSON and back to properly unmarshal into the target type
	data, err := json.Marshal(m.Params)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}