# API Design: Dead Simple to Enterprise Grade

## Design Philosophy: Layered Complexity

### Layer 0: Magic Mode (Zero Code)
```bash
# Just run your app with hawk prefix
hawk myapp --port 8080

# Or as a wrapper
hawk -- python manage.py runserver
hawk -- npm start
hawk -- go run main.go
```

Automatically detects and visualizes:
- stdout/stderr as logs
- HTTP requests (if applicable)
- File changes
- Process metrics (CPU, memory)

### Layer 1: One-Line Integration
```python
# Python
import hawk; hawk.auto()

# JavaScript  
require('hawk').auto();

# Go
import _ "github.com/hawk-tui/auto"

# Rust
use hawk::auto;
```

### Layer 2: Simple Customization
```python
import hawk

tui = hawk.TUI("My App")
tui.log("Server started", level="SUCCESS")
tui.metric("requests_per_second", 145)
tui.config("database_url", default="localhost:5432")
```

### Layer 3: Structured Integration
```python
from hawk import TUI, Dashboard, Logger, ConfigPanel

class MyAppTUI:
    def __init__(self):
        self.tui = TUI("my-app")
        self.setup_dashboard()
        self.setup_config()
    
    def setup_dashboard(self):
        dash = self.tui.dashboard("overview")
        dash.add_metric("RPS", self.get_rps, refresh=1.0)
        dash.add_status("Database", self.check_db)
        dash.add_chart("Response Times", self.get_response_times)
    
    def setup_config(self):
        config = self.tui.config_panel("settings")
        config.add_field("port", type="int", default=8080, restart_required=True)
        config.add_field("debug", type="bool", default=False)
        config.add_field("log_level", type="enum", options=["DEBUG", "INFO", "WARN", "ERROR"])
```

### Layer 4: Enterprise Features
```python
from hawk import TUI, Security, RemoteAccess, AuditLog

tui = TUI("production-app")
tui.enable_security(
    auth_provider="ldap",
    permissions=["read", "config_write"],
    audit_log=AuditLog("./hawk-audit.log")
)
tui.enable_remote_access(
    port=9090,
    ssl_cert="./cert.pem",
    allowed_hosts=["admin.company.com"]
)
```

## Core API Components

### 1. Logging Interface
```python
# Simple
hawk.log("Message")
hawk.log("Error occurred", level="ERROR")
hawk.log("Success!", level="SUCCESS", color="green")

# Structured
hawk.log({
    "event": "user_login",
    "user_id": 12345,
    "timestamp": datetime.now(),
    "metadata": {"ip": "192.168.1.1"}
})

# With context
with hawk.context("Database Migration"):
    hawk.log("Starting migration")
    migrate_tables()
    hawk.log("Migration complete")
```

### 2. Metrics Interface
```python
# Counters
hawk.counter("requests_total").inc()
hawk.counter("errors", tags={"type": "404"}).inc()

# Gauges
hawk.gauge("active_connections", 42)
hawk.gauge("memory_usage_mb", psutil.virtual_memory().used / 1024 / 1024)

# Histograms
hawk.histogram("response_time", 0.125)

# Custom widgets
hawk.progress("file_processing", current=75, total=100)
hawk.table("active_users", [
    {"id": 1, "name": "Alice", "status": "online"},
    {"id": 2, "name": "Bob", "status": "offline"}
])
```

### 3. Configuration Interface
```python
# Simple key-value
port = hawk.config("port", default=8080, type=int)
debug = hawk.config("debug", default=False, type=bool)

# Complex forms
config_panel = hawk.config_form("Database Settings")
config_panel.add_field("host", type="string", required=True)
config_panel.add_field("port", type="int", min=1, max=65535)
config_panel.add_field("ssl_mode", type="enum", options=["disable", "require", "verify-full"])
config_panel.on_submit(restart_database_connection)
```

### 4. Interactive Widgets
```python
# File browser
selected_file = hawk.file_browser("Select config file", filter="*.json")

# Command runner
hawk.command_panel("Admin Commands", {
    "Clear Cache": clear_cache,
    "Restart Workers": restart_workers,
    "Export Data": export_data
})

# Real-time charts
chart = hawk.chart("Response Times")
chart.add_series("95th percentile", color="red")
chart.add_series("median", color="blue")
chart.update({"95th": 250, "median": 100})
```

## Language-Specific Idioms

### Python: Decorators & Context Managers
```python
@hawk.monitor
def expensive_function():
    return process_data()

@hawk.timed("api_request")
def handle_request(request):
    return response

with hawk.context("File Processing"):
    process_files()
```

### JavaScript: Middleware & Promises
```javascript
// Express middleware
app.use(hawk.middleware());

// Promise wrapping
const result = await hawk.track('database_query', db.query('SELECT * FROM users'));

// Event-based
hawk.on('config_changed', (key, value) => {
    console.log(`Config ${key} changed to ${value}`);
});
```

### Go: Interfaces & Context
```go
// Interface compliance
func (s *Server) Start() error {
    defer hawk.Track("server_start")()
    
    hawk.Log("Starting server")
    return s.listen()
}

// Context integration
func ProcessRequest(ctx context.Context, req *Request) error {
    ctx = hawk.WithContext(ctx, "request_processing")
    defer hawk.LogFromContext(ctx, "Request completed")
    
    return processRequest(ctx, req)
}
```

### Rust: Traits & Macros
```rust
// Trait implementation
impl HawkMonitor for MyService {
    fn monitor(&self) -> hawk::Metrics {
        hawk::metrics! {
            "active_connections" => self.connections.len(),
            "uptime_seconds" => self.start_time.elapsed().as_secs()
        }
    }
}

// Macro convenience
hawk_log!("Server started on port {}", port);
hawk_metric!("requests_total", 1);
```

## Error Handling & Fallbacks

### Graceful Degradation
```python
import hawk

# If hawk TUI is not available, calls become no-ops
hawk.log("This works whether TUI is running or not")
hawk.metric("counter", 1)  # Silently ignored if TUI unavailable

# Optional features
if hawk.is_available():
    hawk.enable_dashboard()
else:
    print("TUI not available, continuing with normal logging")
```

### Performance Considerations
```python
# Lazy evaluation for expensive operations
hawk.metric("complex_calculation", lambda: expensive_calculation())

# Batching for high-frequency events
with hawk.batch():
    for item in items:
        hawk.log(f"Processing {item}")
        hawk.counter("items_processed").inc()
```

## Configuration & Setup

### Zero Configuration
```python
# Works out of the box
import hawk
hawk.auto()  # Uses sensible defaults
```

### File-Based Configuration
```yaml
# hawk.yml (optional)
app_name: "My Application"
auto_detect:
  logs: true
  metrics: true
  configs: true
dashboard:
  refresh_rate: 1000ms
  theme: "dark"
logging:
  level: "INFO"
  format: "structured"
```

### Environment Variables
```bash
HAWK_APP_NAME="My App"
HAWK_THEME="dark"
HAWK_AUTO_DETECT="logs,metrics"
HAWK_REFRESH_RATE="500ms"
```

### Programmatic Configuration
```python
hawk.configure({
    "app_name": "My Application",
    "theme": "dark",
    "auto_detect": ["logs", "metrics"],
    "dashboard": {
        "refresh_rate": "1s",
        "widgets": ["logs", "metrics", "config"]
    }
})
```

## Success Criteria

### For Small Projects
- ✅ One line of code to get started
- ✅ Zero configuration required
- ✅ Immediate visual feedback
- ✅ No impact on existing code

### For Enterprise
- ✅ Granular permission controls
- ✅ Audit logging and compliance
- ✅ Remote monitoring capabilities
- ✅ Custom branding and themes
- ✅ Integration with existing tools

### For Framework Authors
- ✅ Easy to embed in existing frameworks
- ✅ Minimal dependency footprint
- ✅ Extensible widget system
- ✅ Language-agnostic protocols