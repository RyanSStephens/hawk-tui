# API Reference

Complete documentation for the Hawk TUI protocol and client libraries.

## Protocol Overview

Hawk TUI uses a JSON-RPC-like protocol over stdin/stdout for communication between your application and the TUI. The client libraries abstract this protocol, but understanding it helps with debugging and custom integrations.

## Universal Methods

These methods are available in all client libraries with consistent naming.

### Core Functions

#### `hawk.auto()`
Automatically detect and capture logs, metrics, and configuration from your application.

```python
# Python
import hawk
hawk.auto()
```

```javascript
// Node.js
const hawk = require('hawk-tui');
hawk.auto();
```

```go
// Go
import _ "github.com/hawk-tui/hawk-tui/clients/go/auto"
```

#### `hawk.log(message, level?, context?)`
Send a log message to the TUI.

**Parameters:**
- `message` (string): The log message
- `level` (string, optional): Log level (DEBUG, INFO, WARN, ERROR, SUCCESS)
- `context` (object, optional): Additional context data

```python
# Python
hawk.log("Server started")
hawk.log("Database error", level="ERROR")
hawk.log("User login", level="INFO", context={"user_id": 123})
```

```javascript
// Node.js
hawk.log('Server started');
hawk.log('Database error', { level: 'ERROR' });
hawk.log('User login', { level: 'INFO', context: { user_id: 123 } });
```

#### `hawk.metric(name, value, type?)`
Send a metric value to the TUI.

**Parameters:**
- `name` (string): Metric name
- `value` (number): Metric value
- `type` (string, optional): Metric type (counter, gauge, histogram)

```python
# Python
hawk.metric("requests_total", 1)
hawk.metric("response_time", 125.5, type="histogram")
hawk.metric("active_users", 42, type="gauge")
```

```javascript
// Node.js
hawk.metric('requests_total', 1);
hawk.metric('response_time', 125.5, { type: 'histogram' });
hawk.metric('active_users', 42, { type: 'gauge' });
```

### Progress Tracking

#### `hawk.progress(name, total?)`
Create a progress tracker for long-running operations.

**Parameters:**
- `name` (string): Progress tracker name
- `total` (number, optional): Total expected items

**Returns:** Progress object with `update(current)` method

```python
# Python
progress = hawk.progress("Processing files", total=1000)
for i in range(1000):
    process_file(i)
    progress.update(i + 1)
```

```javascript
// Node.js
const progress = hawk.progress('Processing files', { total: 1000 });
for (let i = 0; i < 1000; i++) {
    processFile(i);
    progress.update(i + 1);
}
```

### Configuration

#### `hawk.config(name, default_value, type?)`
Register a configuration parameter that can be changed via the TUI.

**Parameters:**
- `name` (string): Configuration parameter name
- `default_value` (any): Default value
- `type` (string, optional): Type hint (string, int, float, bool)

**Returns:** Current value of the configuration parameter

```python
# Python
port = hawk.config("port", 8080, type="int")
debug = hawk.config("debug", False, type="bool")
log_level = hawk.config("log_level", "INFO", type="string")
```

```javascript
// Node.js
const port = hawk.config('port', 8080, { type: 'int' });
const debug = hawk.config('debug', false, { type: 'bool' });
const logLevel = hawk.config('log_level', 'INFO', { type: 'string' });
```

### Context Management

#### `hawk.context(name)`
Create a context for grouping related operations.

```python
# Python
with hawk.context("Database Migration"):
    migrate_users()
    migrate_products()
```

```javascript
// Node.js
await hawk.context('Database Migration', async () => {
    await migrateUsers();
    await migrateProducts();
});
```

### Status and Health

#### `hawk.status(name, value, status?)`
Update application status information.

**Parameters:**
- `name` (string): Status item name
- `value` (string): Status value
- `status` (string, optional): Status level (ok, warning, error)

```python
# Python
hawk.status("Database", "Connected", status="ok")
hawk.status("Cache", "Disconnected", status="error")
hawk.status("Queue", "3 jobs pending", status="warning")
```

```javascript
// Node.js
hawk.status('Database', 'Connected', { status: 'ok' });
hawk.status('Cache', 'Disconnected', { status: 'error' });
hawk.status('Queue', '3 jobs pending', { status: 'warning' });
```

## Advanced Features

### Custom Dashboards

#### `hawk.dashboard(name)`
Create a custom dashboard with specific widgets.

```python
# Python
dashboard = hawk.dashboard("System Overview")
dashboard.add_metric("CPU Usage", get_cpu_usage, format="{:.1f}%")
dashboard.add_chart("Memory Usage", get_memory_history)
dashboard.add_table("Top Processes", get_top_processes)
```

### Timers and Benchmarks

#### `hawk.timer(name)`
Time operations and automatically send metrics.

```python
# Python
with hawk.timer("database_query"):
    results = db.query("SELECT * FROM users")

# Or as decorator
@hawk.timer("api_request")
def handle_request():
    return process_request()
```

```javascript
// Node.js
const timer = hawk.timer('database_query');
const results = await db.query('SELECT * FROM users');
timer.stop();

// Or with callback
hawk.timer('api_request', () => {
    return processRequest();
});
```

### Events and Alerts

#### `hawk.event(name, data?)`
Send significant events that should be highlighted.

```python
# Python
hawk.event("User Registration", {"user_id": 123, "email": "user@example.com"})
hawk.event("System Restart")
hawk.event("Error Threshold Exceeded", {"error_rate": 15.2})
```

## Language-Specific APIs

### Python Client

```python
import hawk

# Class-based API
class MyApp:
    def __init__(self):
        self.hawk = hawk.TUI("My Application")
        self.metrics = self.hawk.metrics()
        self.logger = self.hawk.logger()
    
    def process(self):
        self.logger.info("Starting process")
        self.metrics.increment("process_count")
```

### Node.js Client

```javascript
const Hawk = require('hawk-tui');

class MyApp {
    constructor() {
        this.hawk = new Hawk('My Application');
        this.metrics = this.hawk.metrics();
        this.logger = this.hawk.logger();
    }
    
    process() {
        this.logger.info('Starting process');
        this.metrics.increment('process_count');
    }
}
```

## Protocol Details

For custom integrations or debugging, here's the underlying protocol:

### Message Format

All messages are JSON objects sent over stdin/stdout:

```json
{
    "type": "log|metric|config|status|progress|event",
    "timestamp": "2023-12-01T10:00:00Z",
    "data": {
        // Type-specific data
    }
}
```

### Log Messages

```json
{
    "type": "log",
    "timestamp": "2023-12-01T10:00:00Z",
    "data": {
        "message": "Server started",
        "level": "INFO",
        "context": {
            "module": "server",
            "pid": 1234
        }
    }
}
```

### Metric Messages

```json
{
    "type": "metric",
    "timestamp": "2023-12-01T10:00:00Z",
    "data": {
        "name": "requests_total",
        "value": 1,
        "metric_type": "counter",
        "labels": {
            "endpoint": "/api/users",
            "method": "GET"
        }
    }
}
```

### Configuration Messages

```json
{
    "type": "config",
    "timestamp": "2023-12-01T10:00:00Z",
    "data": {
        "name": "port",
        "value": 8080,
        "type": "int",
        "description": "Server port number"
    }
}
```

## Error Handling

All client libraries handle errors gracefully:

- Network failures fall back to stdout/stderr
- Invalid data is logged but doesn't crash
- Missing TUI binary continues without visualization

```python
# Python - errors don't affect your application
try:
    hawk.log("This works even if TUI is not available")
    hawk.metric("safe_metric", 42)
except Exception:
    pass  # Hawk errors are caught internally
```

## Performance Considerations

- Messages are batched to reduce overhead
- Metrics are sampled if rate is too high
- Large objects are truncated automatically
- Memory usage is bounded and configurable

## Integration Examples

### Flask/Django

```python
# Flask middleware
from flask import Flask, request
import hawk

app = Flask(__name__)
hawk.auto()

@app.before_request
def before_request():
    hawk.metric("requests_total", 1)
    hawk.log(f"{request.method} {request.path}")

@app.after_request
def after_request(response):
    hawk.metric("response_status", response.status_code)
    return response
```

### Express.js

```javascript
const express = require('express');
const hawk = require('hawk-tui');

const app = express();
hawk.auto();

app.use((req, res, next) => {
    hawk.metric('requests_total', 1);
    hawk.log(`${req.method} ${req.path}`);
    next();
});

app.use((req, res, next) => {
    res.on('finish', () => {
        hawk.metric('response_status', res.statusCode);
    });
    next();
});
```

## Best Practices

1. **Call `hawk.auto()` early** - Place it near the top of your main file
2. **Use descriptive names** - Make metrics and logs searchable
3. **Group related operations** - Use contexts for logical grouping
4. **Don't over-instrument** - Focus on key metrics and events
5. **Handle gracefully** - Ensure your app works without TUI
6. **Use appropriate levels** - INFO for normal operations, ERROR for problems

## Next Steps

- [Configuration Guide](configuration.md) - Customize behavior and appearance
- [Examples](examples.md) - See real-world usage patterns
- [Architecture](architecture.md) - Understand the internals