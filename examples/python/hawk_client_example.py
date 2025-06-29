#!/usr/bin/env python3
"""
Hawk TUI Python Client Example

This example demonstrates how to send messages to the Hawk TUI using the JSON-RPC protocol.
In a real implementation, this would be part of a client library that abstracts the protocol details.
"""

import json
import sys
import time
import threading
from datetime import datetime, timezone
from typing import Any, Dict, List, Optional, Union


class HawkClient:
    """
    Simple Hawk TUI client that sends JSON-RPC messages to stdout.
    
    In a production implementation, this would handle:
    - Message batching for performance
    - Buffering and async sending
    - Error handling and retries
    - Configuration management
    """
    
    def __init__(self, app_name: str = "python-app"):
        self.app_name = app_name
        self.session_id = f"sess_{int(time.time())}"
        self.sequence = 0
        self._buffer = []
        self._buffer_lock = threading.Lock()
        
    def _send_message(self, method: str, params: Dict[str, Any], message_id: Optional[Union[str, int]] = None):
        """Send a JSON-RPC message to the TUI via stdout."""
        self.sequence += 1
        
        message = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params
        }
        
        if message_id is not None:
            message["id"] = message_id
            
        # Add Hawk metadata
        message["hawk_meta"] = {
            "app_name": self.app_name,
            "session_id": self.session_id,
            "sequence": self.sequence
        }
        
        # Send to stdout (this would go to the TUI process)
        json_str = json.dumps(message, default=self._json_serializer)
        print(json_str, flush=True)
        
    def _send_batch(self, messages: List[Dict[str, Any]]):
        """Send a batch of messages for efficiency."""
        batch_json = json.dumps(messages, default=self._json_serializer)
        print(batch_json, flush=True)
        
    def _json_serializer(self, obj):
        """Custom JSON serializer for datetime objects."""
        if isinstance(obj, datetime):
            return obj.isoformat()
        raise TypeError(f"Object of type {type(obj)} is not JSON serializable")
    
    # Logging methods
    def log(self, message: str, level: str = "INFO", 
            context: Optional[Dict[str, Any]] = None, 
            tags: Optional[List[str]] = None,
            component: Optional[str] = None):
        """Send a log message."""
        params = {
            "message": message,
            "level": level,
            "timestamp": datetime.now(timezone.utc)
        }
        
        if context:
            params["context"] = context
        if tags:
            params["tags"] = tags
        if component:
            params["component"] = component
            
        self._send_message("hawk.log", params)
    
    def debug(self, message: str, **kwargs):
        """Send a debug log message."""
        self.log(message, level="DEBUG", **kwargs)
        
    def info(self, message: str, **kwargs):
        """Send an info log message."""
        self.log(message, level="INFO", **kwargs)
        
    def warn(self, message: str, **kwargs):
        """Send a warning log message."""
        self.log(message, level="WARN", **kwargs)
        
    def error(self, message: str, **kwargs):
        """Send an error log message."""
        self.log(message, level="ERROR", **kwargs)
        
    def success(self, message: str, **kwargs):
        """Send a success log message."""
        self.log(message, level="SUCCESS", **kwargs)
    
    # Metrics methods
    def metric(self, name: str, value: float, 
               metric_type: str = "gauge",
               unit: Optional[str] = None,
               tags: Optional[Dict[str, Any]] = None):
        """Send a metric value."""
        params = {
            "name": name,
            "value": value,
            "type": metric_type,
            "timestamp": datetime.now(timezone.utc)
        }
        
        if unit:
            params["unit"] = unit
        if tags:
            params["tags"] = tags
            
        self._send_message("hawk.metric", params)
    
    def counter(self, name: str, value: float = 1, **kwargs):
        """Send a counter metric."""
        self.metric(name, value, metric_type="counter", **kwargs)
        
    def gauge(self, name: str, value: float, **kwargs):
        """Send a gauge metric."""
        self.metric(name, value, metric_type="gauge", **kwargs)
        
    def histogram(self, name: str, value: float, **kwargs):
        """Send a histogram metric."""
        self.metric(name, value, metric_type="histogram", **kwargs)
    
    # Configuration methods
    def config(self, key: str, value: Any = None, 
               config_type: str = "string",
               description: Optional[str] = None,
               default: Any = None,
               **kwargs):
        """Define or update a configuration parameter."""
        params = {
            "key": key,
            "type": config_type
        }
        
        if value is not None:
            params["value"] = value
        if description:
            params["description"] = description
        if default is not None:
            params["default"] = default
            
        # Add any additional parameters (min, max, options, etc.)
        params.update(kwargs)
        
        self._send_message("hawk.config", params)
    
    # Progress tracking methods
    def progress(self, progress_id: str, label: str, 
                current: float, total: float,
                unit: str = "",
                status: str = "in_progress",
                details: Optional[str] = None):
        """Update progress for a long-running operation."""
        params = {
            "id": progress_id,
            "label": label,
            "current": current,
            "total": total,
            "status": status
        }
        
        if unit:
            params["unit"] = unit
        if details:
            params["details"] = details
            
        self._send_message("hawk.progress", params)
    
    def progress_complete(self, progress_id: str, label: str, total: float, unit: str = ""):
        """Mark progress as completed."""
        self.progress(progress_id, label, total, total, unit, status="completed")
    
    def progress_error(self, progress_id: str, label: str, current: float, total: float, 
                      error_message: str, unit: str = ""):
        """Mark progress as failed."""
        self.progress(progress_id, label, current, total, unit, 
                     status="error", details=error_message)
    
    # Dashboard methods
    def dashboard_widget(self, widget_id: str, widget_type: str, 
                        title: str, data: Any,
                        layout: Optional[Dict[str, int]] = None):
        """Create or update a dashboard widget."""
        params = {
            "widget_id": widget_id,
            "type": widget_type,
            "title": title,
            "data": data
        }
        
        if layout:
            params["layout"] = layout
            
        self._send_message("hawk.dashboard", params)
    
    def status_grid(self, widget_id: str, title: str, services: Dict[str, Dict[str, Any]], **kwargs):
        """Create a status grid widget."""
        self.dashboard_widget(widget_id, "status_grid", title, services, **kwargs)
    
    def metric_chart(self, widget_id: str, title: str, chart_data: Dict[str, Any], **kwargs):
        """Create a metric chart widget."""
        self.dashboard_widget(widget_id, "metric_chart", title, chart_data, **kwargs)
    
    def table(self, widget_id: str, title: str, headers: List[str], rows: List[List[Any]], **kwargs):
        """Create a table widget."""
        table_data = {
            "headers": headers,
            "rows": rows
        }
        self.dashboard_widget(widget_id, "table", title, table_data, **kwargs)
    
    # Event methods
    def event(self, event_type: str, title: str, 
              message: Optional[str] = None,
              severity: str = "info",
              data: Optional[Dict[str, Any]] = None):
        """Send an application event."""
        params = {
            "type": event_type,
            "title": title,
            "severity": severity,
            "timestamp": datetime.now(timezone.utc)
        }
        
        if message:
            params["message"] = message
        if data:
            params["data"] = data
            
        self._send_message("hawk.event", params)


def simulate_web_server():
    """Simulate a web server with various TUI updates."""
    
    hawk = HawkClient("demo-web-server")
    
    # Initial setup
    hawk.info("Starting web server simulation", component="server")
    
    # Configuration
    hawk.config("server.port", 8080, "integer", "HTTP server port", min=1, max=65535)
    hawk.config("server.workers", 4, "integer", "Number of worker processes", default=4)
    hawk.config("debug_mode", False, "boolean", "Enable debug logging")
    hawk.config("log_level", "INFO", "enum", "Logging level", 
               options=["DEBUG", "INFO", "WARN", "ERROR"])
    
    # Initial metrics
    hawk.gauge("server.uptime", 0, unit="seconds")
    hawk.gauge("server.active_connections", 0)
    hawk.counter("server.requests_total", 0)
    
    # Status dashboard
    hawk.status_grid("service_status", "Service Status", {
        "HTTP Server": {"status": "starting", "response_time": "0ms"},
        "Database": {"status": "healthy", "response_time": "12ms"},
        "Redis Cache": {"status": "healthy", "response_time": "2ms"},
        "Load Balancer": {"status": "healthy", "response_time": "5ms"}
    })
    
    hawk.success("Server started successfully on port 8080", component="server")
    hawk.event("server_started", "Server Started", "Web server is now accepting connections")
    
    # Simulate server activity
    uptime = 0
    requests_total = 0
    
    for i in range(20):
        time.sleep(0.5)  # Simulate real-time updates
        uptime += 0.5
        
        # Simulate varying load
        active_connections = max(0, 50 + int(30 * (0.5 + 0.5 * (i % 10) / 10)))
        requests_per_second = max(0, 100 + int(50 * (0.5 + 0.5 * (i % 8) / 8)))
        response_time = 50 + (i % 5) * 10  # Varying response time
        
        requests_total += requests_per_second * 0.5
        
        # Update metrics
        hawk.gauge("server.uptime", uptime, unit="seconds")
        hawk.gauge("server.active_connections", active_connections)
        hawk.gauge("server.requests_per_second", requests_per_second, unit="req/s")
        hawk.gauge("server.avg_response_time", response_time, unit="ms")
        hawk.counter("server.requests_total", requests_total)
        
        # Update status with current response time
        hawk.status_grid("service_status", "Service Status", {
            "HTTP Server": {"status": "healthy", "response_time": f"{response_time}ms"},
            "Database": {"status": "healthy", "response_time": "12ms"},
            "Redis Cache": {"status": "healthy", "response_time": "2ms"},
            "Load Balancer": {"status": "healthy", "response_time": "5ms"}
        })
        
        # Occasional log messages
        if i % 5 == 0:
            hawk.info(f"Processed {int(requests_total)} total requests", 
                     context={"active_connections": active_connections})
        
        # Simulate occasional warnings
        if response_time > 80:
            hawk.warn(f"High response time detected: {response_time}ms", 
                     component="performance")
        
        # Progress simulation (file upload)
        if i == 10:
            hawk.info("Starting file upload", component="upload")
            
        if 10 <= i <= 15:
            progress = ((i - 10) / 5) * 100
            hawk.progress("file_upload", "Uploading user_data.csv", 
                         progress, 100, unit="%", 
                         details=f"Uploading to S3 bucket ({progress:.0f}% complete)")
        
        if i == 15:
            hawk.progress_complete("file_upload", "Upload completed", 100, unit="%")
            hawk.success("File upload completed successfully", component="upload")
    
    # Simulate graceful shutdown
    hawk.info("Initiating graceful shutdown", component="server")
    
    for i in range(5):
        remaining_connections = max(0, active_connections - (i + 1) * 10)
        hawk.gauge("server.active_connections", remaining_connections)
        hawk.progress("shutdown", "Graceful shutdown", i + 1, 5, 
                     details=f"Waiting for {remaining_connections} connections to close")
        time.sleep(0.3)
    
    hawk.progress_complete("shutdown", "Shutdown complete", 5)
    hawk.event("server_stopped", "Server Stopped", "Web server shutdown completed", severity="info")
    hawk.success("Server shutdown completed", component="server")


def demonstrate_error_handling():
    """Demonstrate error scenarios and recovery."""
    
    hawk = HawkClient("error-demo")
    
    hawk.info("Starting error handling demonstration")
    
    # Simulate various error conditions
    hawk.error("Database connection failed", 
              context={"error": "connection refused", "host": "db.example.com", "port": 5432},
              component="database")
    
    hawk.warn("Redis cache miss rate high", 
             context={"miss_rate": 85.5, "threshold": 80.0},
             component="cache")
    
    # Simulate error recovery
    hawk.info("Attempting database reconnection", component="database")
    
    for attempt in range(3):
        hawk.progress("db_reconnect", "Reconnecting to database", 
                     attempt + 1, 3, details=f"Attempt {attempt + 1}/3")
        time.sleep(0.5)
        
        if attempt == 2:  # Success on third attempt
            hawk.progress_complete("db_reconnect", "Reconnection successful", 3)
            hawk.success("Database connection restored", component="database")
            break
        else:
            hawk.warn(f"Reconnection attempt {attempt + 1} failed", component="database")


def demonstrate_dashboard_widgets():
    """Demonstrate various dashboard widget types."""
    
    hawk = HawkClient("dashboard-demo")
    
    hawk.info("Creating dashboard widgets")
    
    # Table widget
    hawk.table("user_sessions", "Active User Sessions",
              headers=["User ID", "Username", "Last Activity", "Status"],
              rows=[
                  [1001, "alice", "2 minutes ago", "active"],
                  [1002, "bob", "5 minutes ago", "idle"],
                  [1003, "charlie", "1 minute ago", "active"],
                  [1004, "diana", "10 minutes ago", "idle"]
              ])
    
    # Chart widget with time series data
    chart_data = {
        "series": [
            {
                "name": "CPU Usage",
                "color": "#ff6b6b",
                "data": [
                    {"x": datetime.now().timestamp() - 300, "y": 45.2},
                    {"x": datetime.now().timestamp() - 240, "y": 52.1},
                    {"x": datetime.now().timestamp() - 180, "y": 48.7},
                    {"x": datetime.now().timestamp() - 120, "y": 61.3},
                    {"x": datetime.now().timestamp() - 60, "y": 55.8},
                    {"x": datetime.now().timestamp(), "y": 49.2}
                ]
            },
            {
                "name": "Memory Usage",
                "color": "#4ecdc4",
                "data": [
                    {"x": datetime.now().timestamp() - 300, "y": 68.5},
                    {"x": datetime.now().timestamp() - 240, "y": 70.2},
                    {"x": datetime.now().timestamp() - 180, "y": 69.8},
                    {"x": datetime.now().timestamp() - 120, "y": 72.1},
                    {"x": datetime.now().timestamp() - 60, "y": 71.5},
                    {"x": datetime.now().timestamp(), "y": 70.9}
                ]
            }
        ]
    }
    
    hawk.metric_chart("system_metrics", "System Resource Usage", chart_data,
                     layout={"row": 0, "col": 0, "width": 8, "height": 4})


if __name__ == "__main__":
    print("Hawk TUI Python Client Example", file=sys.stderr)
    print("==============================", file=sys.stderr)
    print("Note: Messages are sent to stdout in JSON-RPC format", file=sys.stderr)
    print("In a real scenario, these would be processed by the Hawk TUI", file=sys.stderr)
    print("", file=sys.stderr)
    
    try:
        # Run the web server simulation
        simulate_web_server()
        
        print("\n--- Error Handling Demo ---", file=sys.stderr)
        demonstrate_error_handling()
        
        print("\n--- Dashboard Widgets Demo ---", file=sys.stderr)
        demonstrate_dashboard_widgets()
        
        print("\nExample completed!", file=sys.stderr)
        
    except KeyboardInterrupt:
        print("\nExample interrupted by user", file=sys.stderr)
    except Exception as e:
        print(f"\nError running example: {e}", file=sys.stderr)
        sys.exit(1)