# Hawk TUI Python Client Library

**Dead simple TUI integration for Python applications** - Turn any Python app into a beautiful, real-time terminal interface with just one line of code.

## 🚀 Quick Start (5-Minute Rule)

### Minute 1: Installation
```bash
# Option 1: Copy the single file (zero dependencies)
curl -O https://raw.githubusercontent.com/hawk-tui/hawk-tui/main/examples/python/hawk.py

# Option 2: Install with pip (coming soon)
# pip install hawk-tui-client

# Option 3: Clone and install
git clone https://github.com/hawk-tui/hawk-tui.git
cd hawk-tui/examples/python
```

### Minute 2: Magic Mode Integration
```python
import hawk
hawk.auto()  # That's it! 🎉

# Your existing code works unchanged
print("Server starting...")  # Automatically appears in TUI
import logging
logging.info("Database connected")  # Automatically monitored
```

### Minute 3: Add Some Metrics
```python
import hawk

hawk.auto("my-awesome-app")

# Simple logging
hawk.log("Server started successfully")
hawk.error("Database connection failed")
hawk.success("All systems operational")

# Easy metrics
hawk.counter("requests_total")
hawk.gauge("cpu_usage", 65.5)
hawk.histogram("response_time", 0.125)

# Simple configuration
port = hawk.config("port", default=8080, description="Server port")
debug = hawk.config("debug", default=False, type="boolean")
```

### Minute 4: Add Monitoring
```python
import hawk

@hawk.monitor
def process_request():
    """This function is automatically monitored"""
    return "processed"

@hawk.timed("database_query")
def query_database():
    """Execution time automatically tracked"""
    return {"users": 100}

with hawk.context("User Registration"):
    validate_user()
    save_to_database()
    send_welcome_email()
```

### Minute 5: See It in Action
```bash
# Run your app with Hawk TUI
hawk -- python your_app.py

# Or run the included Flask demo
python flask_demo.py
# Then visit http://localhost:5000
```

## 🏗️ Layered Complexity Architecture

Hawk TUI follows a layered complexity approach - start simple and add features as needed:

### Layer 0: Magic Mode (Zero Configuration)
```python
import hawk
hawk.auto()  # Automatically detects and monitors everything
```

**What happens automatically:**
- Python `logging` module integration
- Exception tracking
- Basic metrics collection
- Configuration detection

### Layer 1: Simple Functions
```python
import hawk

# Logging
hawk.log("User logged in", level="INFO")
hawk.debug("Processing request")
hawk.warn("High memory usage detected")
hawk.error("Database connection failed")
hawk.success("Deployment completed")

# Metrics
hawk.counter("api_requests")
hawk.gauge("memory_usage", 75.5, unit="%")
hawk.histogram("response_time", 0.234, unit="seconds")

# Configuration
port = hawk.config("port", default=8080, type="integer")
debug = hawk.config("debug", default=False, type="boolean")
log_level = hawk.config("log_level", default="INFO", 
                       type="enum", options=["DEBUG", "INFO", "WARN", "ERROR"])

# Progress tracking
hawk.progress("file_upload", "Uploading file", 75, 100, unit="%")

# Events
hawk.event("deployment", "New Version Deployed", 
          message="Version 2.1.0 deployed successfully")
```

### Layer 2: Decorators and Context Managers
```python
import hawk
from time import sleep

# Function monitoring decorator
@hawk.monitor
def expensive_function():
    sleep(1)
    return "result"

# Timing decorator
@hawk.timed("api_call_duration")
def make_api_call():
    return requests.get("https://api.example.com")

# Custom monitoring
@hawk.monitor(name="user_registration", 
             log_calls=True, 
             track_time=True,
             log_errors=True)
def register_user(email, password):
    # Registration logic here
    return {"user_id": 123}

# Context managers for grouping
with hawk.context("Database Migration"):
    hawk.log("Starting migration")
    migrate_tables()
    update_schema()
    hawk.success("Migration completed")

# Batch operations for performance
with hawk.batch():
    for i in range(1000):
        hawk.counter("items_processed")
        hawk.log(f"Processed item {i}")
    # All messages sent as a single batch
```

### Layer 3: Advanced Features
```python
from hawk_advanced import Dashboard, ConfigPanel, ProgressTracker

# Real-time dashboard
dashboard = Dashboard("my-app")
dashboard.add_metric("cpu_usage", "CPU Usage %", get_cpu_usage, refresh=1.0)
dashboard.add_status("services", "Service Status", check_services, refresh=5.0)
dashboard.add_chart("response_times", "Response Times", get_response_data)

# Configuration panel with validation
config = ConfigPanel("app-settings")
config.add_field("max_connections", "integer", "Maximum connections",
                 default=100, min_value=1, max_value=1000)
config.add_field("debug_mode", "boolean", "Enable debug mode", 
                 default=False, restart_required=True)

# Advanced progress tracking
progress = ProgressTracker()
with progress.start("data_processing", "Processing data", 1000) as task:
    for i in range(1000):
        process_item(i)
        task.update(i + 1, f"Processed {i + 1}/1000 items")
```

## 🌐 Framework Integration Examples

### Flask Integration
```python
from flask import Flask
import hawk

app = Flask(__name__)
hawk.auto("flask-app")

@app.before_request
def log_request():
    hawk.log(f"{request.method} {request.path}")

@app.after_request  
def log_response(response):
    hawk.metric("http_requests", 1, tags={"status": response.status_code})
    return response

@app.route('/api/users')
@hawk.monitor
def get_users():
    return jsonify(users)

if __name__ == '__main__':
    app.run()
```

### Django Integration
```python
# settings.py
import hawk
hawk.auto("django-app")

# middleware.py
class HawkMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        hawk.counter("django_requests")
        response = self.get_response(request)
        hawk.metric("response_time", response.processing_time)
        return response

# views.py
@hawk.monitor
def user_view(request):
    hawk.log(f"User view accessed by {request.user}")
    return render(request, 'users.html')
```

### FastAPI Integration
```python
from fastapi import FastAPI
import hawk

app = FastAPI()
hawk.auto("fastapi-app")

@app.middleware("http")
async def hawk_middleware(request, call_next):
    hawk.counter("fastapi_requests")
    response = await call_next(request)
    hawk.metric("status_codes", 1, tags={"code": response.status_code})
    return response

@app.get("/users")
@hawk.monitor
async def get_users():
    return {"users": []}
```

## 📊 Real-World Examples

### Web Server Monitoring
```python
import hawk
from flask import Flask, request
import time
import psutil

app = Flask(__name__)
hawk.auto("web-server")

# System metrics
def collect_system_metrics():
    while True:
        hawk.gauge("cpu_usage", psutil.cpu_percent())
        hawk.gauge("memory_usage", psutil.virtual_memory().percent)
        hawk.gauge("disk_usage", psutil.disk_usage('/').percent)
        time.sleep(5)

# Request monitoring
@app.before_request
def monitor_request():
    request.start_time = time.time()
    hawk.counter("http_requests", tags={"method": request.method})

@app.after_request
def monitor_response(response):
    duration = time.time() - request.start_time
    hawk.histogram("request_duration", duration)
    hawk.counter("http_responses", tags={"status": response.status_code})
    
    if response.status_code >= 500:
        hawk.error(f"Server error: {response.status_code}")
    elif response.status_code >= 400:
        hawk.warn(f"Client error: {response.status_code}")
    
    return response

# Background metrics collection
import threading
metrics_thread = threading.Thread(target=collect_system_metrics, daemon=True)
metrics_thread.start()
```

### Database Operations
```python
import hawk
import sqlite3
from contextlib import contextmanager

@contextmanager
def database_connection():
    hawk.counter("db_connections")
    conn = sqlite3.connect('app.db')
    try:
        yield conn
    except Exception as e:
        hawk.error(f"Database error: {e}")
        hawk.counter("db_errors")
        raise
    finally:
        conn.close()

@hawk.timed("db_query_duration")
def execute_query(query, params=None):
    with database_connection() as conn:
        cursor = conn.cursor()
        
        hawk.log(f"Executing query: {query[:50]}...", level="DEBUG")
        
        try:
            if params:
                cursor.execute(query, params)
            else:
                cursor.execute(query)
            
            results = cursor.fetchall()
            hawk.counter("db_queries_success")
            hawk.metric("db_rows_returned", len(results))
            
            return results
            
        except Exception as e:
            hawk.error(f"Query failed: {e}")
            hawk.counter("db_queries_failed")
            raise
```

### Background Task Processing
```python
import hawk
from celery import Celery

app = Celery('tasks')
hawk.auto("task-worker")

@app.task
@hawk.monitor
def process_image(image_path):
    with hawk.context("Image Processing"):
        hawk.log(f"Processing image: {image_path}")
        
        # Progress tracking
        with hawk.progress_tracker.start("image_process", "Processing image", 4) as progress:
            progress.update(1, "Loading image")
            image = load_image(image_path)
            
            progress.update(2, "Applying filters")
            filtered = apply_filters(image)
            
            progress.update(3, "Compressing")
            compressed = compress_image(filtered)
            
            progress.update(4, "Saving result")
            save_image(compressed, output_path)
        
        hawk.success(f"Image processed: {image_path}")
        hawk.counter("images_processed")
        
        return output_path

@hawk.monitor
def process_queue():
    """Background queue processor"""
    while True:
        try:
            task = get_next_task()
            if task:
                process_image.delay(task.image_path)
                hawk.counter("tasks_queued")
            else:
                time.sleep(1)
        except Exception as e:
            hawk.error(f"Queue processing error: {e}")
            time.sleep(5)
```

## 🔧 Configuration and Customization

### Environment Variables
```bash
# Application identification
export HAWK_APP_NAME="my-app"
export HAWK_SESSION_ID="prod-session-1"

# Behavior configuration
export HAWK_AUTO_DETECT="true"
export HAWK_THREAD_SAFE="true"
export HAWK_GRACEFUL_FALLBACK="true"

# Performance tuning
export HAWK_BUFFER_SIZE="100"
export HAWK_FLUSH_INTERVAL="0.1"

# Theme and display
export HAWK_THEME="dark"
export HAWK_REFRESH_RATE="1000"
```

### Programmatic Configuration
```python
import hawk

# Configure before use
hawk.configure(
    app_name="my-application",
    buffer_size=200,
    flush_interval=0.05,
    auto_detect=True,
    graceful_fallback=True,
    thread_safe=True,
    debug=False
)

# Or use custom client
from hawk import HawkClient, HawkConfig

config = HawkConfig(
    app_name="custom-app",
    buffer_size=500,
    flush_interval=0.02,
    auto_detect=False
)

client = HawkClient(config)
client.log("Using custom client")
```

### File-Based Configuration
```yaml
# hawk.yml
app_name: "my-application"
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
performance:
  buffer_size: 100
  flush_interval: 100ms
  thread_safe: true
```

## 🚀 Performance and Best Practices

### High-Performance Usage
```python
import hawk
from hawk_advanced import BatchOperations

# Use batch operations for high-frequency events
batch = BatchOperations(batch_size=1000, flush_interval=1.0)

for i in range(10000):
    batch.add_metric("requests", 1)
    batch.add_log(f"Processing request {i}")

# Explicit flush when needed
batch.flush()

# Use lazy evaluation for expensive operations
hawk.metric("expensive_calculation", lambda: expensive_function())

# Batch context for multiple operations
with hawk.batch():
    for item in large_dataset:
        hawk.log(f"Processing {item}")
        hawk.counter("items")
```

### Memory Management
```python
import hawk

# Limit context stack depth
with hawk.context("Level 1"):
    with hawk.context("Level 2"):
        # Avoid too many nested contexts
        pass

# Clean shutdown
import atexit
atexit.register(hawk.shutdown)

# Manual cleanup when needed
hawk.shutdown()
```

### Error Handling
```python
import hawk

# Hawk never breaks your application
try:
    risky_operation()
except Exception as e:
    hawk.error(f"Operation failed: {e}")
    # Application continues normally

# Graceful degradation
if hawk.is_available():
    hawk.log("TUI is running")
else:
    print("TUI not available, using fallback logging")

# Disable Hawk in production if needed
if os.getenv('ENVIRONMENT') == 'production':
    hawk.configure(graceful_fallback=True, debug=False)
```

## 🧪 Testing and Development

### Running Examples
```bash
# Basic example
python -c "import hawk; hawk.auto(); hawk.log('Hello!')"

# Flask demo with full features
python flask_demo.py

# Advanced features demo
python hawk_advanced.py

# Run with actual TUI (requires hawk binary)
hawk -- python flask_demo.py
```

### Testing Your Integration
```python
import hawk
import pytest

def test_hawk_integration():
    """Test that Hawk integration doesn't break application"""
    hawk.auto("test-app")
    
    # These should never raise exceptions
    hawk.log("Test message")
    hawk.metric("test_metric", 42)
    hawk.config("test_config", default="value")
    
    # Test decorators
    @hawk.monitor
    def test_function():
        return "success"
    
    result = test_function()
    assert result == "success"
    
    # Test context managers
    with hawk.context("Test Context"):
        pass
    
    # Clean up
    hawk.shutdown()

def test_graceful_fallback():
    """Test that Hawk works even when TUI is not available"""
    # Hawk should work silently without TUI
    hawk.log("This should not fail")
    assert True  # If we get here, Hawk didn't crash the app
```

### Development Setup
```bash
# Clone repository
git clone https://github.com/hawk-tui/hawk-tui.git
cd hawk-tui/examples/python

# Install development dependencies
pip install -r requirements.txt

# Install in development mode
pip install -e .

# Run tests
python -m pytest tests/

# Format code
black *.py

# Type checking
mypy hawk.py
```

## 🔒 Security and Enterprise Features

### Audit Logging
```python
from hawk_advanced import AuditLogger

audit = AuditLogger("./logs/security_audit.log")

# Log security events
audit.log_access("admin", "/sensitive-data", "READ", "success")
audit.log_config_change("admin", "max_connections", 100, 200)
audit.log_command("admin", "restart_service", ["api-server"], "success")

# Verify log integrity
integrity = audit.verify_integrity()
print(f"Audit log integrity: {integrity['integrity']:.2%}")
```

### Configuration Security
```python
from hawk_advanced import ConfigPanel

config = ConfigPanel("secure-config")

# Sensitive configuration with validation
config.add_field("api_key", "string", "API Key", 
                required=True,
                validation_func=lambda x: len(x) >= 32)

config.add_field("max_connections", "integer", "Max connections",
                default=100, min_value=1, max_value=1000)

# Configuration change callbacks
def on_security_change(key, old_value, new_value):
    audit.log_config_change("system", key, old_value, new_value)
    if key == "api_key":
        hawk.warn("API key changed - verify this was authorized")

config.on_change("api_key", on_security_change)
```

## 🤝 Integration Patterns

### Microservices Architecture
```python
import hawk
import requests

class ServiceClient:
    def __init__(self, service_name, base_url):
        self.service_name = service_name
        self.base_url = base_url
        hawk.auto(f"client-{service_name}")
    
    @hawk.monitor
    def call_service(self, endpoint, **kwargs):
        url = f"{self.base_url}/{endpoint}"
        
        with hawk.context(f"Service Call: {self.service_name}"):
            hawk.log(f"Calling {url}")
            
            try:
                response = requests.get(url, **kwargs)
                hawk.metric("service_response_time", response.elapsed.total_seconds())
                hawk.counter("service_calls", tags={"service": self.service_name, "status": "success"})
                return response.json()
            
            except Exception as e:
                hawk.error(f"Service call failed: {e}")
                hawk.counter("service_calls", tags={"service": self.service_name, "status": "error"})
                raise

# Usage
user_service = ServiceClient("user-service", "http://users.internal")
order_service = ServiceClient("order-service", "http://orders.internal")

users = user_service.call_service("users")
orders = order_service.call_service("orders")
```

### CI/CD Pipeline Integration
```python
import hawk
import subprocess
import sys

def deploy_application():
    hawk.auto("deployment-pipeline")
    
    steps = [
        ("Running tests", "python -m pytest"),
        ("Building application", "docker build -t myapp ."),
        ("Pushing to registry", "docker push myapp:latest"),
        ("Deploying to production", "kubectl apply -f deployment.yml")
    ]
    
    with hawk.context("Application Deployment"):
        for i, (description, command) in enumerate(steps):
            hawk.progress("deployment", description, i + 1, len(steps))
            
            try:
                result = subprocess.run(command.split(), capture_output=True, text=True)
                if result.returncode == 0:
                    hawk.success(f"✓ {description}")
                else:
                    hawk.error(f"✗ {description}: {result.stderr}")
                    sys.exit(1)
            except Exception as e:
                hawk.error(f"✗ {description}: {e}")
                sys.exit(1)
        
        hawk.event("deployment_complete", "Deployment Successful", 
                  "Application deployed to production successfully")

if __name__ == "__main__":
    deploy_application()
```

## 📝 API Reference

### Core Functions

#### `hawk.auto(app_name=None, **kwargs)`
Enable magic mode with automatic detection.

#### `hawk.log(message, level="INFO", **kwargs)`
Send a log message.

#### `hawk.metric(name, value, **kwargs)`
Send a metric value.

#### `hawk.config(key, value=None, **kwargs)`
Define or get configuration.

#### `hawk.progress(progress_id, label, current, total, **kwargs)`
Update progress tracking.

#### `hawk.event(event_type, title, **kwargs)`
Send application event.

### Decorators

#### `@hawk.monitor(func=None, *, name=None, log_calls=True, log_errors=True, track_time=True)`
Monitor function execution.

#### `@hawk.timed(name=None, unit="seconds")`
Time function execution.

### Context Managers

#### `hawk.context(name, log_entry=True, log_exit=True)`
Group related operations.

#### `hawk.batch()`
Batch multiple operations.

### Advanced Classes

#### `Dashboard(name, client=None)`
Create real-time dashboards.

#### `ConfigPanel(name, client=None)`
Create configuration panels.

#### `ProgressTracker(client=None)`
Advanced progress tracking.

#### `AuditLogger(audit_file=None, client=None)`
Security audit logging.

## 🐛 Troubleshooting

### Common Issues

**Q: Hawk messages aren't appearing in the TUI**
```python
# Check if TUI is running
if hawk.is_available():
    print("TUI is available")
else:
    print("TUI not available - messages go to stdout")

# Force flush messages
hawk.shutdown()
```

**Q: High memory usage with lots of messages**
```python
# Use batch operations for high-frequency events
from hawk_advanced import BatchOperations

batch = BatchOperations(batch_size=1000)
for i in range(10000):
    batch.add_log(f"Message {i}")
batch.flush()
```

**Q: Application performance degraded**
```python
# Disable Hawk in performance-critical sections
hawk.configure(graceful_fallback=True)

# Or use lazy evaluation
hawk.metric("expensive", lambda: expensive_calculation())
```

**Q: Threading issues**
```python
# Ensure thread safety is enabled
hawk.configure(thread_safe=True)

# Or use separate clients per thread
import threading

thread_local = threading.local()

def get_client():
    if not hasattr(thread_local, 'client'):
        thread_local.client = HawkClient()
    return thread_local.client
```

### Debug Mode
```python
# Enable debug output
hawk.configure(debug=True)

# Check configuration
print(hawk._get_global_client().config)

# Manual message testing
hawk.log("Test message")
hawk.shutdown()  # Force flush
```

## 📄 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.

## 🔗 Links

- [Main Repository](https://github.com/hawk-tui/hawk-tui)
- [Documentation](https://docs.hawk-tui.dev)
- [Examples](https://github.com/hawk-tui/hawk-tui/tree/main/examples)
- [Issue Tracker](https://github.com/hawk-tui/hawk-tui/issues)

---

**Made with ❤️ by the Hawk TUI team**