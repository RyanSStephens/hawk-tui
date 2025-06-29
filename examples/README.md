# Hawk TUI Protocol Examples

This directory contains examples demonstrating the Hawk TUI communication protocol implementation.

## Files Overview

### Protocol Implementation
- **`protocol_example.go`** - Go example showing how to use the protocol handler
- **`python/hawk_client_example.py`** - Python client library demonstration
- **`nodejs/hawk_client_example.js`** - Node.js client library demonstration

## Running the Examples

### Go Protocol Handler Example
```bash
# From project root
go run examples/protocol_example.go
```

This example demonstrates:
- JSON-RPC message parsing and validation
- Error handling and recovery
- Message batching
- Bidirectional communication (TUI ↔ Client)
- Performance metrics collection

### Python Client Example
```bash
# From project root
python3 examples/python/hawk_client_example.py
```

This example simulates:
- Web server monitoring with real-time metrics
- Configuration management
- Progress tracking (file uploads, shutdowns)
- Error scenarios and recovery
- Dashboard widgets (status grids, charts, tables)

### Node.js Client Example
```bash
# From project root
node examples/nodejs/hawk_client_example.js
```

This example demonstrates:
- Express.js API server simulation
- Real-time WebSocket metrics
- Service health monitoring
- Database reconnection scenarios
- Performance optimization techniques

## Protocol Features Demonstrated

### Message Types
- **Logging**: Structured logs with levels (DEBUG, INFO, WARN, ERROR, SUCCESS)
- **Metrics**: Counters, gauges, and histograms with tags
- **Configuration**: Type-safe config parameters with validation
- **Progress**: Long-running operation tracking
- **Dashboard**: Interactive widgets (charts, tables, status grids)
- **Events**: Application lifecycle and notification events

### Advanced Features
- **Batching**: Multiple messages in single JSON arrays
- **Error Handling**: Graceful degradation and recovery
- **Real-time Updates**: High-frequency metric streaming
- **Bidirectional Communication**: TUI can send commands back to clients
- **Metadata**: Session tracking and message sequencing

## Integration Patterns

### Simple Integration (One-liner)
```python
# Python
import hawk; hawk.auto()

# JavaScript
require('hawk').auto();

# Go
import _ "github.com/hawk-tui/auto"
```

### Structured Integration
```python
from hawk import TUI, Logger, Metrics, Dashboard

tui = TUI("my-app")
logger = Logger()
metrics = Metrics()
dashboard = Dashboard()

# Use throughout your application
logger.info("Application started")
metrics.gauge("cpu_usage", 45.2)
dashboard.status_grid("services", {...})
```

### Enterprise Integration
```python
from hawk import TUI, Security, RemoteAccess

tui = TUI("production-app")
tui.enable_security(auth_provider="ldap")
tui.enable_remote_access(port=9090, ssl_cert="cert.pem")
```

## Message Format

All messages follow JSON-RPC 2.0 specification:

```json
{
  "jsonrpc": "2.0",
  "method": "hawk.log",
  "params": {
    "message": "Server started",
    "level": "INFO",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "id": "msg_001"
}
```

With optional Hawk metadata:
```json
{
  "jsonrpc": "2.0",
  "method": "hawk.metric",
  "params": {...},
  "hawk_meta": {
    "app_name": "my-app",
    "session_id": "sess_123",
    "sequence": 42
  }
}
```

## Performance Considerations

### Client-Side Optimizations
- **Message Batching**: Group multiple messages for efficiency
- **Async Sending**: Non-blocking message transmission
- **Buffer Management**: Configurable flush intervals
- **Error Recovery**: Graceful fallbacks when TUI unavailable

### TUI-Side Optimizations
- **Rate Limiting**: Configurable message processing limits
- **Memory Management**: Bounded memory usage per client
- **Efficient Parsing**: Stream-based JSON processing
- **Message Validation**: Fast validation with early rejection

## Error Handling

The protocol is designed to be error-resistant:

- **Invalid JSON**: Logged and skipped, doesn't crash TUI
- **Unknown Methods**: Generate warnings, don't fail
- **Missing Fields**: Use sensible defaults
- **Malformed Data**: Validate and sanitize input
- **Resource Limits**: Configurable message size and rate limits

## Next Steps

1. **See Protocol Documentation**: `docs/PROTOCOL.md` for complete specification
2. **Review Go Types**: `pkg/types/protocol.go` for message structures
3. **Examine Handler**: `internal/protocol/handler.go` for processing logic
4. **Build Client Libraries**: Use examples as reference for other languages

## Contributing

When adding new message types or features:

1. Update the protocol specification in `docs/PROTOCOL.md`
2. Add Go types to `pkg/types/protocol.go`
3. Implement handlers in `internal/protocol/handler.go`
4. Add examples for each supported language
5. Update tests and documentation