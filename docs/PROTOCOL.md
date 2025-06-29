# Hawk TUI Communication Protocol Specification

## Overview

The Hawk TUI protocol enables language-agnostic communication between client applications and the Hawk TUI framework using JSON-RPC 2.0 over stdin/stdout. This design ensures universal compatibility while maintaining high performance for real-time updates.

## Design Goals

- **Language Agnostic**: Works with any language that can write JSON to stdout
- **High Performance**: Efficient for high-frequency metrics and log updates
- **Extensible**: Easy to add new message types and features
- **Error Resistant**: Malformed messages don't crash the TUI
- **Real-time**: Support for streaming updates and subscriptions

## Transport Layer

### Primary: JSON-RPC over stdin/stdout
- **Bidirectional**: Client → TUI (stdout) and TUI → Client (stdin)
- **Line-delimited**: Each message is a single JSON line
- **UTF-8 encoding**: All text content must be valid UTF-8
- **Newline terminated**: Each message ends with `\n`

### Alternative Transports (Future)
- Named pipes (high-performance local IPC)
- Unix domain sockets
- HTTP endpoints (web-based applications)

## Message Format

All messages follow JSON-RPC 2.0 specification with Hawk-specific extensions:

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.log",
  "params": {
    "message": "Server started successfully",
    "level": "INFO",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "id": "msg_001"
}
```

### Common Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `jsonrpc` | string | Yes | Always "2.0" |
| `method` | string | Yes | Hawk method name (e.g., "hawk.log") |
| `params` | object | No | Method-specific parameters |
| `id` | string/number | No | Request identifier for responses |

### Hawk Extensions

#### Metadata
Every message can include optional metadata:

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.log",
  "params": { /* ... */ },
  "hawk_meta": {
    "app_name": "my-api-server",
    "component": "database",
    "session_id": "sess_123",
    "sequence": 42
  }
}
```

#### Batching
Multiple messages can be sent as a JSON array for efficiency:

```json
[
  {"jsonrpc": "2.0", "method": "hawk.log", "params": {"message": "Processing batch"}},
  {"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "batch_size", "value": 100}},
  {"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "processing_time", "value": 0.125}}
]
```

## Core Message Types

### 1. Logging (`hawk.log`)

Send log messages to be displayed in the TUI.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.log",
  "params": {
    "message": "User authentication successful",
    "level": "INFO",
    "timestamp": "2024-01-15T10:30:00Z",
    "context": {
      "user_id": 12345,
      "ip_address": "192.168.1.100",
      "session_id": "sess_abc123"
    },
    "tags": ["auth", "security"],
    "component": "auth-service"
  }
}
```

#### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `message` | string | Yes | Log message content |
| `level` | string | No | Log level (DEBUG, INFO, WARN, ERROR, SUCCESS) |
| `timestamp` | string | No | ISO 8601 timestamp (auto-generated if missing) |
| `context` | object | No | Structured context data |
| `tags` | array[string] | No | Searchable tags |
| `component` | string | No | Component/module name |

#### Log Levels
- `DEBUG`: Detailed debugging information
- `INFO`: General information (default)
- `WARN`: Warning conditions
- `ERROR`: Error conditions
- `SUCCESS`: Success/completion messages

### 2. Metrics (`hawk.metric`)

Report numerical metrics for dashboards and monitoring.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.metric",
  "params": {
    "name": "api_requests_per_second",
    "value": 145.7,
    "type": "gauge",
    "timestamp": "2024-01-15T10:30:00Z",
    "tags": {
      "endpoint": "/api/users",
      "method": "GET",
      "status": "200"
    },
    "unit": "req/s"
  }
}
```

#### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Metric name (dot notation recommended) |
| `value` | number | Yes | Metric value |
| `type` | string | No | Metric type (counter, gauge, histogram) |
| `timestamp` | string | No | ISO 8601 timestamp |
| `tags` | object | No | Key-value tags for filtering |
| `unit` | string | No | Unit of measurement |

#### Metric Types
- `counter`: Monotonically increasing value (e.g., total requests)
- `gauge`: Point-in-time value (e.g., CPU usage, active connections)
- `histogram`: Distribution of values (e.g., response times)

### 3. Configuration (`hawk.config`)

Define and manage configuration parameters.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.config",
  "params": {
    "key": "database.connection_pool_size",
    "value": 20,
    "type": "integer",
    "description": "Maximum number of database connections",
    "default": 10,
    "min": 1,
    "max": 100,
    "restart_required": true,
    "category": "Database"
  }
}
```

#### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Configuration key (dot notation) |
| `value` | any | No | Current value |
| `type` | string | No | Value type (string, integer, float, boolean, enum) |
| `description` | string | No | Human-readable description |
| `default` | any | No | Default value |
| `min` | number | No | Minimum value (for numbers) |
| `max` | number | No | Maximum value (for numbers) |
| `options` | array | No | Valid options (for enum type) |
| `restart_required` | boolean | No | Whether change requires restart |
| `category` | string | No | Configuration category |

### 4. Progress Tracking (`hawk.progress`)

Display progress bars and status updates.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.progress",
  "params": {
    "id": "file_upload_001",
    "label": "Uploading user_data.csv",
    "current": 75,
    "total": 100,
    "unit": "%",
    "status": "in_progress",
    "details": "Uploading to S3 bucket",
    "estimated_completion": "2024-01-15T10:35:00Z"
  }
}
```

#### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique progress identifier |
| `label` | string | Yes | Progress bar label |
| `current` | number | Yes | Current progress value |
| `total` | number | Yes | Total/maximum value |
| `unit` | string | No | Unit of measurement |
| `status` | string | No | Status (pending, in_progress, completed, error) |
| `details` | string | No | Additional status details |
| `estimated_completion` | string | No | ISO 8601 estimated completion time |

### 5. Dashboard Widgets (`hawk.dashboard`)

Create and update dashboard widgets.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.dashboard",
  "params": {
    "widget_id": "server_status",
    "type": "status_grid",
    "title": "Service Status",
    "data": {
      "Database": {"status": "healthy", "response_time": "12ms"},
      "Redis": {"status": "healthy", "response_time": "2ms"},
      "API Gateway": {"status": "degraded", "response_time": "850ms"}
    },
    "layout": {
      "row": 0,
      "col": 0,
      "width": 6,
      "height": 3
    }
  }
}
```

#### Widget Types
- `status_grid`: Service status overview
- `metric_chart`: Time-series charts
- `table`: Tabular data display
- `text`: Rich text content
- `gauge`: Single metric gauges
- `histogram`: Value distribution charts

### 6. Events (`hawk.event`)

Send application events and notifications.

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.event",
  "params": {
    "type": "deployment_completed",
    "title": "Production Deployment Complete",
    "message": "Version 2.1.4 deployed successfully",
    "severity": "info",
    "timestamp": "2024-01-15T10:30:00Z",
    "data": {
      "version": "2.1.4",
      "duration": "2m 34s",
      "affected_services": ["api", "worker", "scheduler"]
    }
  }
}
```

## Bidirectional Communication

### Client → TUI (stdout)
All message types above are sent from client to TUI.

### TUI → Client (stdin)
TUI can send requests back to client applications:

#### Configuration Updates (`hawk.config_update`)
```json
{
  "jsonrpc": "2.0",
  "method": "hawk.config_update",
  "params": {
    "key": "log_level",
    "value": "DEBUG"
  },
  "id": "config_001"
}
```

#### Command Execution (`hawk.execute`)
```json
{
  "jsonrpc": "2.0",
  "method": "hawk.execute",
  "params": {
    "command": "restart_workers",
    "args": {"force": true}
  },
  "id": "cmd_001"
}
```

#### Data Requests (`hawk.request`)
```json
{
  "jsonrpc": "2.0",
  "method": "hawk.request",
  "params": {
    "type": "metrics",
    "filter": {"component": "database"},
    "timerange": "last_5_minutes"
  },
  "id": "req_001"
}
```

## Error Handling

### Error Response Format
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32600,
    "message": "Invalid Request",
    "data": {
      "invalid_field": "timestamp",
      "reason": "Invalid ISO 8601 format"
    }
  },
  "id": "msg_001"
}
```

### Standard Error Codes
- `-32700`: Parse error (invalid JSON)
- `-32600`: Invalid request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error
- `-32000 to -32099`: Custom Hawk errors

### Graceful Degradation
- Invalid messages are ignored, not fatal
- Missing optional fields use sensible defaults
- Unknown methods log warnings but don't crash
- Malformed JSON is skipped with error logging

## Performance Optimizations

### Batching
Send multiple messages in arrays to reduce overhead:
```json
[
  {"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "cpu", "value": 45.2}},
  {"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "memory", "value": 67.8}},
  {"jsonrpc": "2.0", "method": "hawk.metric", "params": {"name": "disk", "value": 23.1}}
]
```

### Compression
For high-frequency updates, consider:
- Abbreviated field names in production mode
- Delta compression for time-series data
- Message deduplication

### Buffering
- Client libraries should buffer messages for batch sending
- Configurable flush intervals (default: 100ms or 10 messages)
- Emergency flush on critical errors

## Security Considerations

### Input Validation
- All string fields are sanitized for display
- Numeric fields are validated for reasonable ranges
- Timestamp fields must be valid ISO 8601
- JSON depth is limited to prevent DoS

### Resource Limits
- Maximum message size: 1MB
- Maximum batch size: 100 messages
- Rate limiting: 1000 messages/second per client
- Memory usage capped per client session

## Implementation Guidelines

### Client Libraries
1. **Buffering**: Implement message batching for performance
2. **Error Handling**: Graceful fallbacks when TUI is unavailable
3. **Threading**: Non-blocking message sending
4. **Validation**: Client-side validation before sending

### TUI Implementation
1. **Parser**: Robust JSON parsing with error recovery
2. **Validation**: Server-side validation and sanitization
3. **Routing**: Efficient message routing to UI components
4. **Storage**: In-memory state with optional persistence

## Future Extensions

### Planned Features
- Binary protocol option for extreme performance
- Message compression for high-volume scenarios
- Authentication and encryption for remote access
- Plugin system for custom message types
- Schema validation for structured data

### Backward Compatibility
- New fields are always optional
- Deprecated fields are supported for 2 major versions
- Version negotiation for protocol upgrades
- Graceful fallback for unsupported features

## Examples

See the `examples/` directory for complete implementations in:
- Python (`examples/python/`)
- Node.js (`examples/nodejs/`)
- Go (`examples/go/`)

Each example demonstrates:
- Basic logging and metrics
- Configuration management
- Dashboard widgets
- Error handling
- Performance optimization techniques