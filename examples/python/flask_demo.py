#!/usr/bin/env python3
"""
Hawk TUI Flask Demo - Real-World Web Application Example

This demo shows how to integrate Hawk TUI into a Flask web application with:
- HTTP request monitoring
- Database query tracking
- Error handling and logging
- Metrics collection and reporting
- Configuration management
- Real-time dashboards

To run this demo:
1. Install Flask: pip install flask
2. Run: python flask_demo.py
3. In another terminal: hawk -- python flask_demo.py (to see TUI)
4. Visit http://localhost:5000 to interact with the app

The demo includes:
- User management endpoints
- Database simulation
- Error scenarios
- Performance monitoring
- Configuration panel
- Real-time metrics dashboard
"""

import os
import sys
import time
import random
import threading
import sqlite3
from datetime import datetime, timedelta
from typing import Dict, Any, List, Optional
from functools import wraps
import json

# Add the current directory to Python path for imports
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

try:
    from flask import Flask, request, jsonify, render_template_string, abort
except ImportError:
    print("Flask is required for this demo. Install with: pip install flask")
    sys.exit(1)

import hawk
from hawk_advanced import Dashboard, ConfigPanel, ProgressTracker, BatchOperations, AuditLogger


# Initialize Hawk TUI with auto-detection
hawk.auto("flask-demo-app")


class FlaskHawkIntegration:
    """
    Flask integration class that provides comprehensive monitoring.
    
    Features:
    - Automatic request/response monitoring
    - Database query tracking
    - Error handling and logging
    - Performance metrics
    - Security audit logging
    """
    
    def __init__(self, app: Flask):
        self.app = app
        self.dashboard = Dashboard("flask-dashboard")
        self.config_panel = ConfigPanel("flask-config")
        self.progress_tracker = ProgressTracker()
        self.batch_ops = BatchOperations(batch_size=50, flush_interval=2.0)
        self.audit_logger = AuditLogger("./logs/flask_audit.log")
        
        # Request tracking
        self.request_count = 0
        self.error_count = 0
        self.response_times = []
        self.active_requests = 0
        
        # Database simulation
        self.db_queries = 0
        self.db_errors = 0
        self.db_response_times = []
        
        # Thread safety
        self.lock = threading.Lock()
        
        # Setup monitoring
        self._setup_request_monitoring()
        self._setup_dashboard()
        self._setup_configuration()
        self._setup_error_handling()
        
        hawk.success("Flask Hawk integration initialized")
    
    def _setup_request_monitoring(self):
        """Set up automatic request monitoring."""
        
        @self.app.before_request
        def before_request():
            request.start_time = time.time()
            
            with self.lock:
                self.active_requests += 1
            
            # Log request details
            hawk.info(f"{request.method} {request.path}", 
                     context={
                         "method": request.method,
                         "path": request.path,
                         "remote_addr": request.remote_addr,
                         "user_agent": request.headers.get("User-Agent", "unknown")[:100]
                     },
                     component="http")
            
            # Audit logging for sensitive endpoints
            if request.path.startswith('/admin') or request.path.startswith('/api'):
                self.audit_logger.log_access(
                    user=request.headers.get("X-User-ID", "anonymous"),
                    resource=request.path,
                    action=request.method,
                    details={"ip": request.remote_addr}
                )
        
        @self.app.after_request
        def after_request(response):
            duration = time.time() - request.start_time
            
            with self.lock:
                self.active_requests -= 1
                self.request_count += 1
                self.response_times.append(duration)
                
                # Keep only last 100 response times for metrics
                if len(self.response_times) > 100:
                    self.response_times = self.response_times[-100:]
                
                if response.status_code >= 400:
                    self.error_count += 1
            
            # Log response
            level = "ERROR" if response.status_code >= 500 else "WARN" if response.status_code >= 400 else "INFO"
            hawk.log(f"{request.method} {request.path} -> {response.status_code} ({duration:.3f}s)",
                    level=level,
                    context={
                        "status_code": response.status_code,
                        "duration": duration,
                        "content_length": response.content_length
                    },
                    component="http")
            
            # Send metrics via batch operations for better performance
            self.batch_ops.add_metric("http.request_duration", duration, 
                                    tags={"method": request.method, "status": str(response.status_code)})
            self.batch_ops.add_metric("http.requests_total", 1, metric_type="counter")
            
            return response
    
    def _setup_dashboard(self):
        """Set up real-time dashboard widgets."""
        
        # Request metrics
        self.dashboard.add_metric("active_requests", "Active Requests",
                                lambda: self.active_requests, refresh=0.5)
        
        self.dashboard.add_metric("total_requests", "Total Requests",
                                lambda: self.request_count, refresh=2.0)
        
        self.dashboard.add_metric("error_rate", "Error Rate %",
                                lambda: (self.error_count / max(self.request_count, 1)) * 100,
                                unit="%", refresh=5.0)
        
        self.dashboard.add_metric("avg_response_time", "Avg Response Time",
                                lambda: sum(self.response_times) / len(self.response_times) * 1000 
                                       if self.response_times else 0,
                                unit="ms", refresh=1.0)
        
        # Database metrics
        self.dashboard.add_metric("db_queries", "Database Queries",
                                lambda: self.db_queries, refresh=2.0)
        
        # System status
        self.dashboard.add_status("system_status", "System Status",
                                lambda: self._get_system_status(), refresh=5.0)
        
        # Request rate chart
        self.dashboard.add_chart("request_rate", "Request Rate",
                               lambda: self._get_request_rate_data(),
                               chart_type="line", refresh=2.0)
    
    def _setup_configuration(self):
        """Set up configuration panel."""
        
        # Add configuration fields
        self.config_panel.add_field("max_request_size", "integer",
                                  "Maximum request size in MB",
                                  default=16, min_value=1, max_value=100,
                                  category="HTTP")
        
        self.config_panel.add_field("debug_mode", "boolean",
                                  "Enable debug mode", default=False,
                                  restart_required=True, category="Debug")
        
        self.config_panel.add_field("log_level", "enum",
                                  "Logging level", default="INFO",
                                  options=["DEBUG", "INFO", "WARN", "ERROR"],
                                  category="Logging")
        
        self.config_panel.add_field("request_timeout", "integer",
                                  "Request timeout in seconds",
                                  default=30, min_value=5, max_value=300,
                                  category="HTTP")
        
        self.config_panel.add_field("db_pool_size", "integer",
                                  "Database connection pool size",
                                  default=10, min_value=1, max_value=50,
                                  category="Database")
        
        # Configuration change handlers
        self.config_panel.on_change("debug_mode", self._handle_debug_change)
        self.config_panel.on_change("log_level", self._handle_log_level_change)
    
    def _setup_error_handling(self):
        """Set up comprehensive error handling."""
        
        @self.app.errorhandler(404)
        def not_found(error):
            hawk.warn(f"404 Not Found: {request.path}",
                     context={"path": request.path, "method": request.method},
                     component="http")
            return jsonify({"error": "Not found"}), 404
        
        @self.app.errorhandler(500)
        def internal_error(error):
            hawk.error(f"500 Internal Server Error: {str(error)}",
                      context={"path": request.path, "method": request.method},
                      component="http")
            return jsonify({"error": "Internal server error"}), 500
        
        @self.app.errorhandler(Exception)
        def handle_exception(error):
            hawk.error(f"Unhandled exception: {str(error)}",
                      context={
                          "path": request.path,
                          "method": request.method,
                          "exception_type": type(error).__name__
                      },
                      component="app")
            return jsonify({"error": "An unexpected error occurred"}), 500
    
    def _get_system_status(self) -> Dict[str, Any]:
        """Get current system status for dashboard."""
        try:
            # Simulate system checks
            cpu_usage = random.uniform(20, 80)
            memory_usage = random.uniform(40, 90)
            disk_usage = random.uniform(10, 95)
            
            # Database health check
            db_healthy = self.db_errors / max(self.db_queries, 1) < 0.1 if self.db_queries > 0 else True
            
            return {
                "Flask Server": {
                    "status": "healthy" if self.error_count / max(self.request_count, 1) < 0.1 else "degraded",
                    "requests": self.request_count,
                    "errors": self.error_count
                },
                "Database": {
                    "status": "healthy" if db_healthy else "degraded",
                    "queries": self.db_queries,
                    "errors": self.db_errors
                },
                "System": {
                    "status": "healthy" if cpu_usage < 90 and memory_usage < 95 else "degraded",
                    "cpu": f"{cpu_usage:.1f}%",
                    "memory": f"{memory_usage:.1f}%",
                    "disk": f"{disk_usage:.1f}%"
                }
            }
        except Exception as e:
            hawk.error(f"Error getting system status: {e}")
            return {"error": str(e)}
    
    def _get_request_rate_data(self) -> Dict[str, Any]:
        """Get request rate data for chart."""
        try:
            now = time.time()
            
            # Generate sample data for demonstration
            data_points = []
            for i in range(20):
                timestamp = now - (19 - i) * 5  # 5-second intervals
                # Simulate varying request rate
                rate = max(0, 10 + 5 * random.sin(i * 0.5) + random.uniform(-2, 2))
                data_points.append({"x": timestamp, "y": rate})
            
            return {
                "series": [{
                    "name": "Requests/sec",
                    "color": "#4ecdc4",
                    "data": data_points
                }]
            }
        except Exception as e:
            hawk.error(f"Error generating chart data: {e}")
            return {"series": []}
    
    def _handle_debug_change(self, key: str, old_value: Any, new_value: Any):
        """Handle debug mode configuration change."""
        self.audit_logger.log_config_change("system", key, old_value, new_value)
        if new_value:
            hawk.warn("Debug mode enabled - this should not be used in production!")
        else:
            hawk.info("Debug mode disabled")
    
    def _handle_log_level_change(self, key: str, old_value: Any, new_value: Any):
        """Handle log level configuration change."""
        self.audit_logger.log_config_change("system", key, old_value, new_value)
        hawk.info(f"Log level changed from {old_value} to {new_value}")
    
    def simulate_database_query(self, query_type: str, duration: Optional[float] = None) -> Dict[str, Any]:
        """Simulate database query with monitoring."""
        if duration is None:
            duration = random.uniform(0.001, 0.1)  # 1ms to 100ms
        
        with self.lock:
            self.db_queries += 1
        
        # Simulate query execution
        start_time = time.time()
        time.sleep(duration)
        actual_duration = time.time() - start_time
        
        self.db_response_times.append(actual_duration)
        if len(self.db_response_times) > 100:
            self.db_response_times = self.db_response_times[-100:]
        
        # Log database query
        hawk.debug(f"Database query: {query_type} ({actual_duration:.3f}s)",
                  context={"query_type": query_type, "duration": actual_duration},
                  component="database")
        
        # Send metrics
        self.batch_ops.add_metric("db.query_duration", actual_duration,
                                tags={"query_type": query_type})
        
        # Simulate occasional database errors
        if random.random() < 0.05:  # 5% error rate
            with self.lock:
                self.db_errors += 1
            error_msg = f"Database error during {query_type}"
            hawk.error(error_msg, component="database")
            raise Exception(error_msg)
        
        return {"success": True, "duration": actual_duration, "rows": random.randint(1, 100)}


# Create Flask app
app = Flask(__name__)
app.config['SECRET_KEY'] = 'demo-secret-key'

# Initialize Hawk integration
hawk_integration = FlaskHawkIntegration(app)

# In-memory "database" for demo
users_db = [
    {"id": 1, "name": "Alice Smith", "email": "alice@example.com", "active": True},
    {"id": 2, "name": "Bob Johnson", "email": "bob@example.com", "active": True},
    {"id": 3, "name": "Charlie Brown", "email": "charlie@example.com", "active": False},
]


# Utility decorators
def monitor_endpoint(name: Optional[str] = None):
    """Decorator to monitor endpoint performance."""
    def decorator(f):
        endpoint_name = name or f.__name__
        
        @wraps(f)
        def wrapper(*args, **kwargs):
            with hawk.context(f"Endpoint: {endpoint_name}"):
                return f(*args, **kwargs)
        return wrapper
    return decorator


def require_auth(f):
    """Simple authentication decorator for demo."""
    @wraps(f)
    def wrapper(*args, **kwargs):
        auth_header = request.headers.get('Authorization')
        if not auth_header or not auth_header.startswith('Bearer '):
            hawk.warn("Unauthorized access attempt",
                     context={"path": request.path, "ip": request.remote_addr},
                     component="auth")
            return jsonify({"error": "Authentication required"}), 401
        
        # Log successful authentication
        hawk.debug("User authenticated", component="auth")
        return f(*args, **kwargs)
    return wrapper


# Routes
@app.route('/')
def index():
    """Home page with API documentation."""
    hawk.info("Home page accessed")
    
    html = """
    <!DOCTYPE html>
    <html>
    <head>
        <title>Flask Hawk TUI Demo</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; }
            .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
            .method { font-weight: bold; color: #007bff; }
            code { background: #e9ecef; padding: 2px 4px; border-radius: 3px; }
        </style>
    </head>
    <body>
        <h1>Flask Hawk TUI Demo</h1>
        <p>This demo shows Hawk TUI integration with Flask. Check your TUI to see real-time monitoring!</p>
        
        <h2>Available Endpoints:</h2>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/</code> - This home page
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/api/users</code> - List all users
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/api/users/&lt;id&gt;</code> - Get specific user
        </div>
        
        <div class="endpoint">
            <span class="method">POST</span> <code>/api/users</code> - Create new user (requires auth)
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/api/stats</code> - Application statistics
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/simulate/load</code> - Simulate high load
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/simulate/error</code> - Trigger error for testing
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <code>/admin/config</code> - View configuration (requires auth)
        </div>
        
        <h3>Authentication:</h3>
        <p>Protected endpoints require: <code>Authorization: Bearer demo-token</code></p>
        
        <h3>Example Commands:</h3>
        <pre>
curl http://localhost:5000/api/users
curl http://localhost:5000/api/users/1
curl -H "Authorization: Bearer demo-token" http://localhost:5000/admin/config
curl http://localhost:5000/simulate/load
        </pre>
    </body>
    </html>
    """
    return html


@app.route('/api/users')
@monitor_endpoint("list_users")
def list_users():
    """List all users."""
    try:
        # Simulate database query
        query_result = hawk_integration.simulate_database_query("SELECT * FROM users")
        
        hawk.info(f"Retrieved {len(users_db)} users")
        return jsonify({"users": users_db, "count": len(users_db)})
        
    except Exception as e:
        hawk.error(f"Failed to list users: {str(e)}")
        return jsonify({"error": "Failed to retrieve users"}), 500


@app.route('/api/users/<int:user_id>')
@monitor_endpoint("get_user")
def get_user(user_id):
    """Get specific user by ID."""
    try:
        # Simulate database query
        query_result = hawk_integration.simulate_database_query("SELECT * FROM users WHERE id = ?")
        
        user = next((u for u in users_db if u["id"] == user_id), None)
        if not user:
            hawk.warn(f"User not found: {user_id}")
            return jsonify({"error": "User not found"}), 404
        
        hawk.info(f"Retrieved user: {user['name']}")
        return jsonify(user)
        
    except Exception as e:
        hawk.error(f"Failed to get user {user_id}: {str(e)}")
        return jsonify({"error": "Failed to retrieve user"}), 500


@app.route('/api/users', methods=['POST'])
@require_auth
@monitor_endpoint("create_user")
def create_user():
    """Create new user."""
    try:
        data = request.get_json()
        if not data or not data.get('name') or not data.get('email'):
            hawk.warn("Invalid user creation request", context={"data": data})
            return jsonify({"error": "Name and email are required"}), 400
        
        # Simulate database insertion with progress tracking
        with hawk_integration.progress_tracker.start("create_user", "Creating user", 3) as progress:
            progress.update(1, "Validating data")
            time.sleep(0.1)
            
            # Simulate database query
            progress.update(2, "Inserting into database")
            query_result = hawk_integration.simulate_database_query("INSERT INTO users")
            
            new_user = {
                "id": max([u["id"] for u in users_db]) + 1,
                "name": data["name"],
                "email": data["email"],
                "active": data.get("active", True)
            }
            users_db.append(new_user)
            
            progress.update(3, "User created successfully")
        
        hawk.success(f"Created new user: {new_user['name']}")
        hawk_integration.audit_logger.log_command(
            user="demo-user",
            command="create_user",
            args=[new_user["name"], new_user["email"]],
            result="success"
        )
        
        return jsonify(new_user), 201
        
    except Exception as e:
        hawk.error(f"Failed to create user: {str(e)}")
        return jsonify({"error": "Failed to create user"}), 500


@app.route('/api/stats')
@monitor_endpoint("get_stats")
def get_stats():
    """Get application statistics."""
    try:
        # Simulate complex statistics calculation
        with hawk.context("Statistics Calculation"):
            hawk.debug("Calculating application statistics")
            
            # Simulate database queries for statistics
            hawk_integration.simulate_database_query("SELECT COUNT(*) FROM users")
            hawk_integration.simulate_database_query("SELECT COUNT(*) FROM requests")
            hawk_integration.simulate_database_query("SELECT AVG(response_time) FROM requests")
            
            stats = {
                "total_users": len(users_db),
                "active_users": len([u for u in users_db if u["active"]]),
                "total_requests": hawk_integration.request_count,
                "error_count": hawk_integration.error_count,
                "error_rate": (hawk_integration.error_count / max(hawk_integration.request_count, 1)) * 100,
                "avg_response_time": (sum(hawk_integration.response_times) / len(hawk_integration.response_times)) 
                                   if hawk_integration.response_times else 0,
                "database_queries": hawk_integration.db_queries,
                "database_errors": hawk_integration.db_errors,
                "uptime": time.time() - app.start_time if hasattr(app, 'start_time') else 0
            }
        
        hawk.info("Statistics retrieved successfully")
        return jsonify(stats)
        
    except Exception as e:
        hawk.error(f"Failed to get statistics: {str(e)}")
        return jsonify({"error": "Failed to retrieve statistics"}), 500


@app.route('/simulate/load')
@monitor_endpoint("simulate_load")
def simulate_load():
    """Simulate high load for testing."""
    try:
        hawk.warn("Starting load simulation")
        
        # Simulate multiple concurrent operations
        with hawk_integration.progress_tracker.start("load_simulation", "Simulating high load", 50) as progress:
            for i in range(50):
                # Simulate various operations
                if i % 10 == 0:
                    hawk_integration.simulate_database_query("SELECT * FROM heavy_table", 
                                                           duration=random.uniform(0.05, 0.2))
                elif i % 5 == 0:
                    hawk_integration.simulate_database_query("UPDATE stats SET value = ?")
                else:
                    hawk_integration.simulate_database_query("SELECT id FROM users")
                
                # Send some metrics
                hawk.counter("load_simulation.operations")
                hawk.gauge("load_simulation.memory_usage", random.uniform(50, 95))
                
                progress.update(i + 1, f"Completed operation {i + 1}/50")
                time.sleep(0.02)  # Small delay to simulate work
        
        hawk.success("Load simulation completed")
        return jsonify({"message": "Load simulation completed", "operations": 50})
        
    except Exception as e:
        hawk.error(f"Load simulation failed: {str(e)}")
        return jsonify({"error": "Load simulation failed"}), 500


@app.route('/simulate/error')
@monitor_endpoint("simulate_error")
def simulate_error():
    """Trigger an error for testing error handling."""
    error_type = request.args.get('type', 'generic')
    
    hawk.warn(f"Simulating {error_type} error as requested")
    
    if error_type == 'database':
        try:
            hawk_integration.simulate_database_query("INVALID SQL QUERY")
        except Exception as e:
            return jsonify({"error": "Database error simulated", "details": str(e)}), 500
    
    elif error_type == 'timeout':
        hawk.error("Simulating timeout error")
        time.sleep(2)  # Simulate long operation
        return jsonify({"error": "Request timed out"}), 408
    
    elif error_type == 'auth':
        hawk.error("Simulating authentication error")
        return jsonify({"error": "Authentication failed"}), 401
    
    else:
        # Generic error
        hawk.error("Simulating generic application error")
        raise Exception(f"Simulated {error_type} error for testing")


@app.route('/admin/config')
@require_auth
@monitor_endpoint("admin_config")
def admin_config():
    """View application configuration (admin only)."""
    try:
        hawk.info("Admin configuration accessed")
        
        config_data = {
            "current_config": hawk_integration.config_panel.get_all_values(),
            "schema": hawk_integration.config_panel.export_schema(),
            "dashboard_stats": hawk_integration.dashboard.get_stats(),
            "batch_stats": hawk_integration.batch_ops.get_stats(),
            "audit_integrity": hawk_integration.audit_logger.verify_integrity()
        }
        
        return jsonify(config_data)
        
    except Exception as e:
        hawk.error(f"Failed to get admin config: {str(e)}")
        return jsonify({"error": "Failed to retrieve configuration"}), 500


# Background monitoring task
def background_monitoring():
    """Background task for continuous monitoring."""
    while True:
        try:
            # Send system metrics every 5 seconds
            hawk.gauge("system.cpu_usage", random.uniform(20, 80))
            hawk.gauge("system.memory_usage", random.uniform(40, 90))
            hawk.gauge("system.disk_usage", random.uniform(10, 95))
            
            # Application metrics
            hawk.gauge("app.active_users", len([u for u in users_db if u["active"]]))
            hawk.gauge("app.total_users", len(users_db))
            
            time.sleep(5)
            
        except Exception as e:
            hawk.error(f"Background monitoring error: {str(e)}")
            time.sleep(10)


if __name__ == '__main__':
    # Record start time for uptime calculation
    app.start_time = time.time()
    
    # Start background monitoring thread
    monitor_thread = threading.Thread(target=background_monitoring, daemon=True)
    monitor_thread.start()
    
    # Log application startup
    hawk.event("app_startup", "Flask Demo Application Started",
              message="Flask application is starting up with Hawk TUI integration",
              data={"port": 5000, "debug": False})
    
    hawk.success("Flask demo application starting on http://localhost:5000")
    hawk.info("Visit http://localhost:5000 for API documentation and testing")
    hawk.info("Use 'Authorization: Bearer demo-token' header for protected endpoints")
    
    try:
        # Run Flask app
        app.run(host='0.0.0.0', port=5000, debug=False, threaded=True)
    except KeyboardInterrupt:
        hawk.info("Application shutdown requested")
    except Exception as e:
        hawk.error(f"Application crashed: {str(e)}")
    finally:
        # Cleanup
        hawk_integration.dashboard.shutdown()
        hawk.event("app_shutdown", "Flask Demo Application Stopped")
        hawk.success("Flask demo application stopped")
        hawk.shutdown()