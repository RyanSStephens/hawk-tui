# Python Client

The Hawk TUI Python client provides seamless integration with Python applications.

## Installation

```bash
pip install hawk-tui
```

## Quick Start

```python
import hawk

# Zero-configuration setup
hawk.auto()

# Your existing code works unchanged
print("Server starting...")
import logging
logging.info("Application initialized")
```

## API Reference

### Core Functions

#### `hawk.auto()`
Automatically captures logs, metrics, and configuration from your application.

```python
import hawk
hawk.auto()

# Now all print() and logging calls are captured
print("This appears in the TUI")
logging.info("So does this")
```

#### `hawk.log(message, level=None, **context)`
Send a log message to the TUI.

```python
hawk.log("Server started")
hawk.log("Database error", level="ERROR")
hawk.log("User action", level="INFO", user_id=123, action="login")
```

**Levels:** DEBUG, INFO, WARN, ERROR, SUCCESS

#### `hawk.metric(name, value, type=None, **labels)`
Send a metric to the TUI.

```python
hawk.metric("requests_total", 1)
hawk.metric("response_time", 125.5, type="histogram")
hawk.metric("active_users", 42, type="gauge", region="us-east")
```

**Types:** counter, gauge, histogram

### Progress Tracking

#### `hawk.progress(name, total=None)`
Create a progress tracker for long-running operations.

```python
# With known total
progress = hawk.progress("Processing files", total=1000)
for i in range(1000):
    process_file(i)
    progress.update(i + 1)

# With unknown total
progress = hawk.progress("Loading data")
for item in data_stream:
    process_item(item)
    progress.increment()

# As context manager
with hawk.progress("Migration", total=len(tables)) as p:
    for i, table in enumerate(tables):
        migrate_table(table)
        p.update(i + 1)
```

### Configuration Management

#### `hawk.config(name, default=None, type=None, description=None)`
Register a configuration parameter.

```python
# Basic configuration
port = hawk.config("port", default=8080, type=int)
debug = hawk.config("debug", default=False, type=bool)

# With description
log_level = hawk.config(
    "log_level", 
    default="INFO", 
    type=str,
    description="Application log level"
)

# Use in your application
if debug:
    logging.getLogger().setLevel(logging.DEBUG)
```

### Context Management

#### `hawk.context(name)`
Group related operations under a named context.

```python
# As context manager
with hawk.context("Database Migration"):
    with hawk.context("Users Table"):
        migrate_users()
    with hawk.context("Products Table"):
        migrate_products()

# Manual context control
ctx = hawk.context("API Request")
ctx.enter()
process_request()
ctx.exit()
```

### Status Updates

#### `hawk.status(name, value, status=None)`
Update application status information.

```python
hawk.status("Database", "Connected", status="ok")
hawk.status("Cache", "Redis unavailable", status="error")
hawk.status("Queue", f"{queue.size()} jobs pending", status="warning")
```

**Status levels:** ok, warning, error

### Events and Alerts

#### `hawk.event(name, data=None, level="info")`
Send significant events.

```python
hawk.event("User Registration", {"user_id": 123, "email": "user@example.com"})
hawk.event("System Restart", level="warning")
hawk.event("Critical Error", {"error": str(e)}, level="error")
```

### Decorators

#### `@hawk.monitor`
Automatically monitor function calls.

```python
@hawk.monitor
def process_request(request):
    # Function timing and error handling is automatic
    return handle_request(request)

@hawk.monitor(name="custom_operation")
def complex_operation():
    # Custom monitoring name
    return do_work()
```

#### `@hawk.timed(name)`
Time function execution.

```python
@hawk.timed("database_query")
def fetch_users():
    return db.query("SELECT * FROM users")

# Results in automatic metrics:
# - database_query_duration_seconds
# - database_query_calls_total
```

#### `@hawk.retry(max_attempts=3)`
Monitor retry operations.

```python
@hawk.retry(max_attempts=3)
def unreliable_operation():
    if random.random() < 0.7:
        raise Exception("Simulated failure")
    return "success"
```

## Class-Based API

For more complex applications, use the class-based API:

```python
from hawk import TUI, MetricsCollector, Logger

class MyApplication:
    def __init__(self):
        self.tui = TUI("My Application")
        self.metrics = MetricsCollector(self.tui)
        self.logger = Logger(self.tui)
        
        # Create custom dashboard
        self.dashboard = self.tui.dashboard("Overview")
        self.dashboard.add_metric("Requests/sec", self.get_request_rate)
        self.dashboard.add_chart("Response Times", self.get_response_times)
    
    def get_request_rate(self):
        return self.metrics.get_rate("requests_total")
    
    def get_response_times(self):
        return self.metrics.get_histogram("response_time")
```

## Framework Integration

### Flask

```python
from flask import Flask, request, g
import hawk
import time

app = Flask(__name__)
hawk.auto()

@app.before_request
def before_request():
    g.start_time = time.time()
    hawk.metric("http_requests_total", 1, method=request.method, endpoint=request.endpoint)

@app.after_request
def after_request(response):
    duration = time.time() - g.start_time
    hawk.metric("http_request_duration_seconds", duration, 
                method=request.method, status=response.status_code)
    return response

@app.route('/api/users')
def get_users():
    with hawk.context("Database Query"):
        users = User.query.all()
        hawk.log(f"Retrieved {len(users)} users")
        return {"users": [u.to_dict() for u in users]}
```

### Django

```python
# middleware.py
import hawk
import time

class HawkMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response
        hawk.auto()

    def __call__(self, request):
        start_time = time.time()
        
        hawk.metric("django_requests_total", 1, 
                   method=request.method, view=request.resolver_match.view_name)
        
        response = self.get_response(request)
        
        duration = time.time() - start_time
        hawk.metric("django_request_duration_seconds", duration,
                   method=request.method, status=response.status_code)
        
        return response

# settings.py
MIDDLEWARE = [
    'myapp.middleware.HawkMiddleware',
    # ... other middleware
]
```

### FastAPI

```python
from fastapi import FastAPI, Request
import hawk
import time

app = FastAPI()
hawk.auto()

@app.middleware("http")
async def hawk_middleware(request: Request, call_next):
    start_time = time.time()
    
    hawk.metric("fastapi_requests_total", 1, 
               method=request.method, path=request.url.path)
    
    response = await call_next(request)
    
    duration = time.time() - start_time
    hawk.metric("fastapi_request_duration_seconds", duration,
               method=request.method, status=response.status_code)
    
    return response

@app.get("/api/users")
async def get_users():
    with hawk.context("Database Query"):
        users = await fetch_users_from_db()
        hawk.log(f"Retrieved {len(users)} users")
        return {"users": users}
```

### Celery

```python
from celery import Celery
import hawk

app = Celery('myapp')
hawk.auto()

@app.task
@hawk.monitor
def process_data(data_id):
    with hawk.context(f"Processing data {data_id}"):
        hawk.log(f"Starting processing for {data_id}")
        
        # Simulate work with progress
        progress = hawk.progress("Processing steps", total=100)
        for i in range(100):
            process_step(i)
            progress.update(i + 1)
            time.sleep(0.1)
        
        hawk.log(f"Completed processing for {data_id}", level="SUCCESS")
        return {"status": "completed", "data_id": data_id}
```

## Data Science Integration

### Pandas

```python
import pandas as pd
import hawk

hawk.auto()

def analyze_dataset(df):
    hawk.log(f"Analyzing dataset with {len(df)} rows, {len(df.columns)} columns")
    hawk.metric("dataset_rows", len(df))
    hawk.metric("dataset_memory_mb", df.memory_usage(deep=True).sum() / 1024**2)
    
    with hawk.context("Data Cleaning"):
        # Track missing values
        missing_data = df.isnull().sum()
        hawk.metric("missing_values_total", missing_data.sum())
        
        # Clean data with progress
        progress = hawk.progress("Cleaning columns", total=len(df.columns))
        for i, col in enumerate(df.columns):
            clean_column(df, col)
            progress.update(i + 1)
    
    hawk.status("Data Quality", "Good", status="ok")
    return df
```

### Machine Learning

```python
import hawk
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier

hawk.auto()

def train_model(X, y):
    with hawk.context("Model Training"):
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)
        
        hawk.log(f"Training set: {len(X_train)} samples")
        hawk.log(f"Test set: {len(X_test)} samples")
        
        # Track training progress
        model = RandomForestClassifier(n_estimators=100)
        
        with hawk.timer("model_training"):
            model.fit(X_train, y_train)
        
        # Evaluate model
        train_score = model.score(X_train, y_train)
        test_score = model.score(X_test, y_test)
        
        hawk.metric("model_train_accuracy", train_score)
        hawk.metric("model_test_accuracy", test_score)
        hawk.metric("model_features", len(X.columns))
        
        if test_score > 0.9:
            hawk.status("Model Quality", "Excellent", status="ok")
        elif test_score > 0.8:
            hawk.status("Model Quality", "Good", status="warning")
        else:
            hawk.status("Model Quality", "Poor", status="error")
    
    return model
```

## Configuration

### Environment Variables

```python
import os
import hawk

# Configure via environment
hawk.config_from_env({
    'HAWK_APP_NAME': 'My Python App',
    'HAWK_LOG_LEVEL': 'INFO',
    'HAWK_METRICS_ENABLED': 'true'
})

hawk.auto()
```

### Configuration File

```python
# config.yaml
app_name: "My Python Application"
log_level: "INFO"
auto_detect:
  logs: true
  metrics: true
  config: true
```

```python
import hawk
import yaml

with open('config.yaml') as f:
    config = yaml.safe_load(f)

hawk.configure(config)
hawk.auto()
```

## Error Handling

The Python client handles errors gracefully and never crashes your application:

```python
import hawk

# This is safe even if Hawk TUI is not available
hawk.auto()
hawk.log("This works regardless")
hawk.metric("safe_metric", 42)

# Hawk errors are logged but don't propagate
try:
    risky_operation()
except Exception as e:
    hawk.log(f"Operation failed: {e}", level="ERROR")
    # Your error handling continues normally
```

## Performance Tips

1. **Batch operations** when possible:
```python
# Better: batch multiple metrics
with hawk.batch():
    hawk.metric("cpu_usage", get_cpu())
    hawk.metric("memory_usage", get_memory())
    hawk.metric("disk_usage", get_disk())
```

2. **Use sampling** for high-frequency metrics:
```python
# Sample 10% of requests
if random.random() < 0.1:
    hawk.metric("request_details", detailed_info)
```

3. **Avoid blocking operations** in metric collection:
```python
# Good: async metrics
async def collect_metrics():
    cpu = await get_cpu_async()
    hawk.metric("cpu_usage", cpu)

# Bad: blocking metrics in main thread
def slow_metric():
    time.sleep(1)  # Don't do this
    return expensive_calculation()
```

## Debugging

Enable debug logging to troubleshoot issues:

```python
import logging
import hawk

logging.basicConfig(level=logging.DEBUG)
hawk.set_debug(True)
hawk.auto()

# Now you'll see internal Hawk TUI messages
```

## Examples

See the [examples/python/](../../examples/python/) directory for complete working examples:

- [Basic Usage](../../examples/python/demo.py)
- [Web Server](../../examples/python/web_server.py)
- [Data Processing](../../examples/python/data_processing.py)
- [Machine Learning](../../examples/python/ml_training.py)