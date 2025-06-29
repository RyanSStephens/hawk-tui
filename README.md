# 🦅 Hawk TUI

**Transform any command-line application into a beautiful, interactive TUI in minutes.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

![Hawk TUI Demo](docs/demo.gif)

## What is Hawk TUI?

Hawk TUI is a **universal TUI framework** that can transform any application in any programming language into a rich, interactive terminal interface. Whether you're monitoring a web server, running database migrations, or training ML models, Hawk TUI provides instant visualization with **zero architectural changes** to your existing code.

## Quick Start (< 5 minutes)

### 1. Install Hawk TUI
```bash
# macOS/Linux
curl -sSL https://hawk-tui.dev/install.sh | sh

# Or build from source
git clone https://github.com/hawk-tui/hawk.git
cd hawk
go build ./cmd/hawk
```

### 2. Add to Your Python App
```python
import hawk
hawk.auto()  # One line - that's it!

# Your existing code works unchanged
print("Server starting...")  # Appears in TUI
logger.info("Database connected")  # Automatically captured
metrics.counter("requests").inc()  # Automatically graphed
```

### 3. Run with TUI
```bash
python your_app.py | hawk
```

**That's it!** Your application now has a beautiful TUI interface.

## Core Philosophy

### Drop-in Integration
- **One line of code** to get started
- **Zero dependencies** for host applications  
- **Works with existing logging/metrics**
- **Never breaks your application**

### Universal Language Support
```python
# Python
import hawk; hawk.auto()
```
```javascript
// Node.js  
require('hawk').auto();
```
```go
// Go
import _ "github.com/hawk-tui/go-client/auto"
```
```rust
// Rust
use hawk::auto;
```

### Enterprise Ready
- **Security & Compliance**: Local-only communication, audit trails
- **Scale & Performance**: Multi-instance support, resource limits
- **Remote Monitoring**: Secure tunneling for production environments

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Your App      │    │   Hawk TUI      │    │   Language      │
│   (Any Lang)    │◄──►│   (Go Binary)   │◄──►│   Client        │
│                 │    │                 │    │   (Thin Library)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
     JSON over             Bubble Tea           Python/Node/Go
     stdin/stdout          Renderer            Rust/Java/etc
```

- **Universal Protocol**: JSON-RPC over stdin/stdout works with any language
- **High Performance**: 60 FPS updates, 100K+ messages/second throughput
- **Zero Dependencies**: Host applications don't install TUI libraries
- **Extensible**: Plugin system for custom widgets and integrations

## Features

### Real-time Monitoring
- **Live Metrics**: Counters, gauges, histograms with auto-scaling charts
- **Intelligent Logging**: Filtering, search, context-aware display
- **Interactive Dashboards**: Customizable widgets and layouts
- **Configuration Management**: Live parameter editing with validation

### Developer Experience
- **Vim-like Navigation**: Intuitive keyboard shortcuts
- **Beautiful Themes**: Professional dark/light themes
- **Smart Search**: Filter logs, metrics, and configuration
- **Responsive Design**: Adapts to any terminal size

### Production Ready
- **Security**: Input validation, resource limits, audit logging
- **Performance**: Memory efficient, graceful degradation
- **Monitoring**: Built-in performance metrics and health checks
- **Remote Access**: Secure tunneling for production monitoring

## Examples

### Web Server Monitoring
```python
from flask import Flask
import hawk

app = Flask(__name__)
tui = hawk.TUI("E-commerce API")

@app.route('/api/products')
def get_products():
    with hawk.context("Database Query"):
        products = db.query("SELECT * FROM products")
        hawk.metric("db_query_time", timer.elapsed())
        hawk.log(f"Retrieved {len(products)} products")
        return products

if __name__ == '__main__':
    app.run()
```

### Database Migration
```python
import hawk

tui = hawk.TUI("Database Migration")

def migrate_table(table_name, total_records):
    progress = hawk.progress(f"Migrating {table_name}", total=total_records)
    
    for i, record in enumerate(get_records(table_name)):
        migrate_record(record)
        progress.update(i + 1)
        
        if i % 100 == 0:
            hawk.log(f"Migrated {i+1}/{total_records} records")
    
    hawk.log(f"Completed {table_name}", level="SUCCESS")
```

### Machine Learning Training
```python
import hawk

monitor = hawk.TUI("ResNet Training")
dashboard = monitor.dashboard("Training Metrics")

for epoch in range(100):
    for batch_idx, (data, targets) in enumerate(dataloader):
        loss = train_step(data, targets)
        
        hawk.metric("batch_loss", loss.item())
        hawk.metric("learning_rate", optimizer.lr)
        
        if batch_idx % 10 == 0:
            hawk.log(f"Epoch {epoch}, Batch {batch_idx}: Loss = {loss:.4f}")
```

## Use Cases

### Web Development
- API endpoint monitoring and performance tracking
- Database query analysis and optimization
- Request/response logging with filtering
- Configuration hot-reloading and validation

### DevOps & Infrastructure  
- Container orchestration dashboards
- Log aggregation and real-time analysis
- Metrics collection and alerting interfaces
- CI/CD pipeline visualization and control

### Data Processing
- ETL pipeline monitoring and control
- Data quality dashboards and validation
- Processing job status and progress tracking
- Performance profiling and optimization

### Game Development
- Asset loading progress and performance
- Debug console with live parameter tweaking
- Performance metrics and frame rate analysis
- Live configuration and gameplay tuning

## Language Support

| Language | Status | Installation | Example |
|----------|--------|--------------|---------|
| **Python** | Ready | `pip install hawk-tui` | [Flask Demo](examples/python/flask_demo.py) |
| **Node.js** | In Progress | `npm install hawk-tui` | [Express Demo](examples/nodejs/express_demo.js) |
| **Go** | In Progress | `go get hawk-tui/client` | [Gin Demo](examples/go/gin_demo.go) |
| **Rust** | Planned | `cargo add hawk-tui` | Coming Soon |
| **Java** | Planned | Maven/Gradle | Coming Soon |

## Getting Started

### Installation Options

#### Option 1: Install Script (Recommended)
```bash
curl -sSL https://hawk-tui.dev/install.sh | sh
```

#### Option 2: Download Binary
```bash
# Download latest release
wget https://github.com/hawk-tui/hawk/releases/latest/download/hawk-linux-amd64
chmod +x hawk-linux-amd64
sudo mv hawk-linux-amd64 /usr/local/bin/hawk
```

#### Option 3: Build from Source
```bash
git clone https://github.com/hawk-tui/hawk.git
cd hawk
go build ./cmd/hawk
```

### Basic Usage

#### Magic Mode (Zero Configuration)
```python
import hawk
hawk.auto()  # Automatically detects logs, metrics, configs

# Your existing application code
logger.info("Server started")
metrics.counter("requests").inc()
config_value = os.getenv("PORT", 8080)
```

#### Structured Integration
```python
import hawk

# Setup application monitoring
tui = hawk.TUI("My Application")
dashboard = tui.dashboard("Overview")

# Add metrics and monitoring
dashboard.add_metric("Requests/sec", get_request_rate)
dashboard.add_chart("Response Times", get_response_times)
dashboard.add_table("Active Users", get_active_users)

# Configuration management
config = tui.config_panel("Settings")
config.add_field("port", type="int", default=8080)
config.add_field("debug", type="bool", default=False)
```

#### Running Your Application
```bash
# Basic usage
python your_app.py | hawk

# With configuration
python your_app.py | hawk --theme dark --refresh-rate 500ms

# Remote monitoring
python your_app.py | hawk --remote --port 9090
```

## Advanced Features

### Dashboard Creation
```python
dashboard = hawk.dashboard("System Overview")

# Add various widget types
dashboard.add_metric("CPU Usage", get_cpu_usage, format="{:.1f}%")
dashboard.add_gauge("Memory", get_memory_usage, max_value=100)
dashboard.add_chart("Network I/O", get_network_stats, chart_type="line")
dashboard.add_table("Processes", get_top_processes, columns=["PID", "Name", "CPU"])
dashboard.add_status("Services", get_service_status)
```

### Configuration Management
```python
config = hawk.config_panel("Database Settings")
config.add_field("host", type="string", required=True, default="localhost")
config.add_field("port", type="int", min=1, max=65535, default=5432)
config.add_field("ssl_mode", type="enum", options=["disable", "require", "verify-full"])

# React to configuration changes
@config.on_change("host")
def reconnect_database(new_host):
    db.reconnect(host=new_host)
```

### Performance Monitoring
```python
# Decorators for automatic monitoring
@hawk.monitor
@hawk.timed("api_request")
def handle_request(request):
    return process_request(request)

# Context managers for sections
with hawk.context("Database Migration"):
    with hawk.context("Table: users"):
        migrate_users_table()
```

## Configuration

### Environment Variables
```bash
export HAWK_APP_NAME="My Application"
export HAWK_THEME="dark"
export HAWK_REFRESH_RATE="1000ms"
export HAWK_AUTO_DETECT="logs,metrics,config"
```

### Configuration File (`hawk.yml`)
```yaml
app_name: "My Application"
theme: "dark"
refresh_rate: 1000ms

auto_detect:
  logs: true
  metrics: true
  configs: true

dashboard:
  widgets: ["logs", "metrics", "status"]
  layout: "vertical"

security:
  enable_audit_log: true
  max_message_rate: 1000
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup
```bash
git clone https://github.com/hawk-tui/hawk.git
cd hawk
go mod download
go run cmd/hawk/main.go
```

### Running Tests
```bash
go test ./...
cd examples/python && python -m pytest test_hawk.py
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=hawk-tui/hawk&type=Date)](https://star-history.com/#hawk-tui/hawk&Date)

## Links

- **Website**: https://hawk-tui.dev
- **Documentation**: https://docs.hawk-tui.dev  
- **Examples**: [examples/](examples/)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)
- **Roadmap**: [ROADMAP.md](ROADMAP.md)

---

**Made with care by developers, for developers.**

*Transform your CLI tools into beautiful, interactive experiences.*