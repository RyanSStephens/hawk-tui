#!/usr/bin/env python3
"""
Hawk TUI Python Client Library - Dead Simple Integration

This is the main client library that provides layered complexity:
- Layer 0: hawk.auto() - Magic mode with zero configuration
- Layer 1: Simple functions like hawk.log(), hawk.metric()
- Layer 2: Structured integration with classes and decorators
- Layer 3: Advanced features in hawk_advanced.py

Usage Examples:

Layer 0 - Magic mode:
    import hawk
    hawk.auto()  # That's it! Auto-detects and monitors everything

Layer 1 - Simple functions:
    import hawk
    hawk.log("Server started")
    hawk.metric("requests_per_second", 145)
    hawk.config("port", default=8080)

Layer 2 - Decorators and context managers:
    @hawk.monitor
    def my_function():
        pass
    
    with hawk.context("Database Migration"):
        migrate_tables()

Requirements:
- Python 3.7+
- No external dependencies
- Thread-safe
- Works with or without Hawk TUI running
"""

import json
import sys
import time
import threading
import logging
import os
import atexit
import functools
import warnings
from datetime import datetime, timezone
from typing import Any, Dict, List, Optional, Union, Callable, ContextManager
from contextlib import contextmanager
from dataclasses import dataclass, asdict
from queue import Queue, Empty
import traceback


# Global state management
_global_client = None
_auto_enabled = False
_context_stack = []


@dataclass
class HawkConfig:
    """Configuration for Hawk TUI client."""
    app_name: str = "python-app"
    buffer_size: int = 100
    flush_interval: float = 0.1  # seconds
    auto_detect: bool = True
    graceful_fallback: bool = True
    thread_safe: bool = True
    debug: bool = False


class HawkMessage:
    """Represents a single Hawk TUI message."""
    
    def __init__(self, method: str, params: Dict[str, Any], message_id: Optional[str] = None):
        self.method = method
        self.params = params
        self.message_id = message_id
        self.timestamp = datetime.now(timezone.utc)
    
    def to_jsonrpc(self, app_name: str, session_id: str, sequence: int) -> Dict[str, Any]:
        """Convert to JSON-RPC format."""
        message = {
            "jsonrpc": "2.0",
            "method": self.method,
            "params": self.params
        }
        
        if self.message_id:
            message["id"] = self.message_id
        
        # Add Hawk metadata
        message["hawk_meta"] = {
            "app_name": app_name,
            "session_id": session_id,
            "sequence": sequence,
            "timestamp": self.timestamp.isoformat()
        }
        
        return message


class HawkClient:
    """
    Core Hawk TUI client with buffering, batching, and thread safety.
    
    This client handles:
    - Message buffering and batching for performance
    - Thread-safe operations
    - Graceful fallback when TUI is not available
    - Auto-detection of common patterns
    """
    
    def __init__(self, config: Optional[HawkConfig] = None):
        self._config = config or HawkConfig()
        self.session_id = f"sess_{int(time.time())}"
        self.sequence = 0
        self.is_available = True
        
        # Thread safety
        self._lock = threading.RLock()
        self._sequence_lock = threading.Lock()
        
        # Message buffering
        self._buffer = Queue(maxsize=self._config.buffer_size * 2)
        self._flush_timer = None
        self._shutdown = False
        
        # Start background thread for message processing
        if self._config.thread_safe:
            self._worker_thread = threading.Thread(target=self._message_worker, daemon=True)
            self._worker_thread.start()
        
        # Register cleanup
        atexit.register(self.shutdown)
        
        # Auto-detection setup
        if self._config.auto_detect:
            self._setup_auto_detection()
    
    def _get_next_sequence(self) -> int:
        """Thread-safe sequence number generation."""
        with self._sequence_lock:
            self.sequence += 1
            return self.sequence
    
    def _setup_auto_detection(self):
        """Set up automatic detection of logging and other patterns."""
        try:
            # Intercept Python logging
            self._setup_logging_integration()
        except Exception as e:
            if self._config.debug:
                print(f"Warning: Could not set up auto-detection: {e}", file=sys.stderr)
    
    def _setup_logging_integration(self):
        """Integrate with Python's logging module."""
        class HawkLogHandler(logging.Handler):
            def __init__(self, client):
                super().__init__()
                self.client = client
            
            def emit(self, record):
                try:
                    level_map = {
                        logging.DEBUG: "DEBUG",
                        logging.INFO: "INFO", 
                        logging.WARNING: "WARN",
                        logging.ERROR: "ERROR",
                        logging.CRITICAL: "ERROR"
                    }
                    
                    level = level_map.get(record.levelno, "INFO")
                    message = self.format(record)
                    
                    context = {
                        "logger": record.name,
                        "module": record.module,
                        "function": record.funcName,
                        "line": record.lineno
                    }
                    
                    if record.exc_info:
                        context["exception"] = self.format(record)
                    
                    self.client._send_message_async("hawk.log", {
                        "message": message,
                        "level": level,
                        "context": context,
                        "component": record.name
                    })
                except Exception:
                    pass  # Never let Hawk break the application
        
        # Add handler to root logger
        handler = HawkLogHandler(self)
        handler.setLevel(logging.INFO)
        logging.getLogger().addHandler(handler)
    
    def _message_worker(self):
        """Background worker thread for processing messages."""
        batch = []
        last_flush = time.time()
        
        while not self._shutdown:
            try:
                # Try to get a message with timeout
                try:
                    message = self._buffer.get(timeout=0.05)
                    batch.append(message)
                except Empty:
                    pass
                
                # Flush conditions
                current_time = time.time()
                should_flush = (
                    len(batch) >= self._config.buffer_size or
                    (batch and current_time - last_flush >= self._config.flush_interval) or
                    self._shutdown
                )
                
                if should_flush and batch:
                    self._flush_batch(batch)
                    batch = []
                    last_flush = current_time
                    
            except Exception as e:
                if self._config.debug:
                    print(f"Hawk worker error: {e}", file=sys.stderr)
                # Clear the batch to prevent infinite errors
                batch = []
    
    def _flush_batch(self, messages: List[HawkMessage]):
        """Flush a batch of messages to stdout."""
        if not messages:
            return
        
        try:
            if len(messages) == 1:
                # Single message
                msg = messages[0]
                jsonrpc = msg.to_jsonrpc(
                    self._config.app_name, 
                    self.session_id, 
                    self._get_next_sequence()
                )
                json_str = json.dumps(jsonrpc, default=self._json_serializer)
            else:
                # Batch of messages
                batch = []
                for msg in messages:
                    jsonrpc = msg.to_jsonrpc(
                        self._config.app_name,
                        self.session_id,
                        self._get_next_sequence()
                    )
                    batch.append(jsonrpc)
                json_str = json.dumps(batch, default=self._json_serializer)
            
            print(json_str, flush=True)
            
        except Exception as e:
            if self._config.debug:
                print(f"Hawk flush error: {e}", file=sys.stderr)
            if not self._config.graceful_fallback:
                raise
    
    def _json_serializer(self, obj):
        """Custom JSON serializer for special types."""
        if isinstance(obj, datetime):
            return obj.isoformat()
        elif hasattr(obj, '__dict__'):
            return obj.__dict__
        elif hasattr(obj, '_asdict'):
            return obj._asdict()
        raise TypeError(f"Object of type {type(obj)} is not JSON serializable")
    
    def _send_message_async(self, method: str, params: Dict[str, Any], message_id: Optional[str] = None):
        """Send message asynchronously via buffer."""
        if self._shutdown:
            return
        
        try:
            message = HawkMessage(method, params, message_id)
            
            if self._config.thread_safe:
                # Use buffer and worker thread
                try:
                    self._buffer.put_nowait(message)
                except:
                    # Buffer full, drop oldest and try again
                    try:
                        self._buffer.get_nowait()
                        self._buffer.put_nowait(message)
                    except:
                        pass  # Give up gracefully
            else:
                # Direct flush
                self._flush_batch([message])
                
        except Exception as e:
            if self._config.debug:
                print(f"Hawk send error: {e}", file=sys.stderr)
            if not self._config.graceful_fallback:
                raise
    
    def _send_message_sync(self, method: str, params: Dict[str, Any], message_id: Optional[str] = None):
        """Send message synchronously."""
        message = HawkMessage(method, params, message_id)
        self._flush_batch([message])
    
    def shutdown(self):
        """Gracefully shutdown the client."""
        if self._shutdown:
            return
        
        self._shutdown = True
        
        # Flush any remaining messages
        if hasattr(self, '_buffer'):
            remaining = []
            while True:
                try:
                    remaining.append(self._buffer.get_nowait())
                except Empty:
                    break
            
            if remaining:
                self._flush_batch(remaining)
        
        # Wait for worker thread
        if hasattr(self, '_worker_thread') and self._worker_thread.is_alive():
            self._worker_thread.join(timeout=1.0)
    
    # Core API methods
    def log(self, message: str, level: str = "INFO", 
            context: Optional[Dict[str, Any]] = None,
            tags: Optional[List[str]] = None,
            component: Optional[str] = None,
            **kwargs):
        """Send a log message."""
        params = {
            "message": str(message),
            "level": level.upper(),
            "timestamp": datetime.now(timezone.utc).isoformat()
        }
        
        # Add current context from stack
        if _context_stack:
            if context is None:
                context = {}
            context["hawk_context"] = " > ".join(_context_stack)
        
        if context:
            params["context"] = context
        if tags:
            params["tags"] = tags
        if component:
            params["component"] = component
        
        # Add any extra kwargs
        params.update(kwargs)
        
        self._send_message_async("hawk.log", params)
    
    def metric(self, name: str, value: Union[int, float],
               metric_type: str = "gauge",
               unit: Optional[str] = None,
               tags: Optional[Dict[str, Any]] = None,
               **kwargs):
        """Send a metric value."""
        params = {
            "name": name,
            "value": float(value),
            "type": metric_type,
            "timestamp": datetime.now(timezone.utc).isoformat()
        }
        
        if unit:
            params["unit"] = unit
        if tags:
            params["tags"] = tags
        
        params.update(kwargs)
        
        self._send_message_async("hawk.metric", params)
    
    def config(self, key: str, value: Any = None,
               config_type: Optional[str] = None,
               description: Optional[str] = None,
               default: Any = None,
               **kwargs):
        """Define or get a configuration parameter."""
        params = {"key": key}
        
        if value is not None:
            params["value"] = value
        if config_type:
            params["type"] = config_type
        elif value is not None:
            # Auto-detect type
            if isinstance(value, bool):
                params["type"] = "boolean"
            elif isinstance(value, int):
                params["type"] = "integer"
            elif isinstance(value, float):
                params["type"] = "float"
            else:
                params["type"] = "string"
        
        if description:
            params["description"] = description
        if default is not None:
            params["default"] = default
        
        params.update(kwargs)
        
        self._send_message_async("hawk.config", params)
        
        # Return current value or default
        return value if value is not None else default
    
    def progress(self, progress_id: str, label: str,
                current: Union[int, float], total: Union[int, float],
                unit: str = "", status: str = "in_progress",
                details: Optional[str] = None,
                **kwargs):
        """Update progress for a long-running operation."""
        params = {
            "id": progress_id,
            "label": label,
            "current": float(current),
            "total": float(total),
            "status": status
        }
        
        if unit:
            params["unit"] = unit
        if details:
            params["details"] = details
        
        params.update(kwargs)
        
        self._send_message_async("hawk.progress", params)
    
    def event(self, event_type: str, title: str,
              message: Optional[str] = None,
              severity: str = "info",
              data: Optional[Dict[str, Any]] = None,
              **kwargs):
        """Send an application event."""
        params = {
            "type": event_type,
            "title": title,
            "severity": severity,
            "timestamp": datetime.now(timezone.utc).isoformat()
        }
        
        if message:
            params["message"] = message
        if data:
            params["data"] = data
        
        params.update(kwargs)
        
        self._send_message_async("hawk.event", params)
    
    # Convenience methods
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
    
    def counter(self, name: str, value: Union[int, float] = 1, **kwargs):
        """Send a counter metric."""
        self.metric(name, value, metric_type="counter", **kwargs)
    
    def gauge(self, name: str, value: Union[int, float], **kwargs):
        """Send a gauge metric."""
        self.metric(name, value, metric_type="gauge", **kwargs)
    
    def histogram(self, name: str, value: Union[int, float], **kwargs):
        """Send a histogram metric."""
        self.metric(name, value, metric_type="histogram", **kwargs)


# Global convenience functions for Layer 1 API
def _get_global_client() -> HawkClient:
    """Get or create the global client instance."""
    global _global_client
    if _global_client is None:
        _global_client = HawkClient()
    return _global_client


def auto(app_name: Optional[str] = None, **config_kwargs):
    """
    Layer 0: Magic mode - Enable Hawk TUI with zero configuration.
    
    This function:
    - Automatically detects your application patterns
    - Integrates with Python logging
    - Sets up metrics collection
    - Provides instant visualization
    
    Args:
        app_name: Application name (auto-detected if not provided)
        **config_kwargs: Additional configuration options
    
    Example:
        import hawk
        hawk.auto()  # That's it!
    """
    global _global_client, _auto_enabled
    
    if _auto_enabled:
        return
    
    # Auto-detect app name
    if app_name is None:
        app_name = os.path.basename(sys.argv[0]) if sys.argv else "python-app"
        if app_name.endswith('.py'):
            app_name = app_name[:-3]
    
    # Create config with auto-detection enabled
    config = HawkConfig(
        app_name=app_name,
        auto_detect=True,
        **config_kwargs
    )
    
    _global_client = HawkClient(config)
    _auto_enabled = True
    
    # Send startup event
    _global_client.event("app_started", "Application Started", 
                        f"{app_name} initialized with Hawk TUI")


# Layer 1: Simple function API
def log(message: str, level: str = "INFO", **kwargs):
    """Send a log message. Works with or without auto() being called."""
    _get_global_client().log(message, level, **kwargs)


def metric(name: str, value: Union[int, float], **kwargs):
    """Send a metric value. Works with or without auto() being called."""
    _get_global_client().metric(name, value, **kwargs)


def config(key: str, value: Any = None, **kwargs):
    """Define or get a configuration parameter."""
    client = _get_global_client()
    return client.config(key, value, **kwargs)


def progress(progress_id: str, label: str, current: Union[int, float], 
             total: Union[int, float], **kwargs):
    """Update progress for a long-running operation."""
    _get_global_client().progress(progress_id, label, current, total, **kwargs)


def event(event_type: str, title: str, **kwargs):
    """Send an application event."""
    _get_global_client().event(event_type, title, **kwargs)


# Convenience logging functions
def debug(message: str, **kwargs):
    """Send a debug log message."""
    log(message, "DEBUG", **kwargs)


def info(message: str, **kwargs):
    """Send an info log message."""
    log(message, "INFO", **kwargs)


def warn(message: str, **kwargs):
    """Send a warning log message."""
    log(message, "WARN", **kwargs)


def error(message: str, **kwargs):
    """Send an error log message."""
    log(message, "ERROR", **kwargs)


def success(message: str, **kwargs):
    """Send a success log message."""
    log(message, "SUCCESS", **kwargs)


# Convenience metric functions
def counter(name: str, value: Union[int, float] = 1, **kwargs):
    """Send a counter metric."""
    metric(name, value, metric_type="counter", **kwargs)


def gauge(name: str, value: Union[int, float], **kwargs):
    """Send a gauge metric."""
    metric(name, value, metric_type="gauge", **kwargs)


def histogram(name: str, value: Union[int, float], **kwargs):
    """Send a histogram metric."""
    metric(name, value, metric_type="histogram", **kwargs)


# Layer 2: Decorators and Context Managers
def monitor(func: Optional[Callable] = None, *, 
           name: Optional[str] = None,
           log_calls: bool = True,
           log_errors: bool = True,
           track_time: bool = True):
    """
    Decorator to automatically monitor function calls.
    
    Args:
        func: Function to monitor (when used as @monitor)
        name: Custom name for monitoring (defaults to function name)
        log_calls: Whether to log function calls
        log_errors: Whether to log exceptions
        track_time: Whether to track execution time
    
    Example:
        @hawk.monitor
        def my_function():
            return "result"
        
        @hawk.monitor(name="api_call", track_time=True)
        def api_call():
            return make_request()
    """
    def decorator(f):
        monitor_name = name or f.__name__
        
        @functools.wraps(f)
        def wrapper(*args, **kwargs):
            client = _get_global_client()
            
            if log_calls:
                client.debug(f"Calling {monitor_name}", component="monitor")
            
            start_time = time.time() if track_time else None
            
            try:
                result = f(*args, **kwargs)
                
                if track_time:
                    duration = time.time() - start_time
                    client.histogram(f"{monitor_name}.duration", duration, unit="seconds")
                
                if log_calls:
                    client.debug(f"Completed {monitor_name}", component="monitor")
                
                return result
                
            except Exception as e:
                if track_time:
                    duration = time.time() - start_time
                    client.histogram(f"{monitor_name}.error_duration", duration, unit="seconds")
                
                if log_errors:
                    client.error(f"Error in {monitor_name}: {str(e)}", 
                               context={"exception": str(e), "traceback": traceback.format_exc()},
                               component="monitor")
                
                client.counter(f"{monitor_name}.errors")
                raise
        
        return wrapper
    
    if func is None:
        # Used with arguments: @monitor(name="custom")
        return decorator
    else:
        # Used without arguments: @monitor
        return decorator(func)


def timed(name: Optional[str] = None, unit: str = "seconds"):
    """
    Decorator to time function execution.
    
    Args:
        name: Metric name (defaults to function name)
        unit: Time unit for the metric
    
    Example:
        @hawk.timed("api_request_time")
        def make_api_request():
            return requests.get("https://api.example.com")
    """
    def decorator(func):
        metric_name = name or f"{func.__name__}.duration"
        
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            start_time = time.time()
            try:
                result = func(*args, **kwargs)
                duration = time.time() - start_time
                histogram(metric_name, duration, unit=unit)
                return result
            except Exception:
                duration = time.time() - start_time
                histogram(f"{metric_name}.error", duration, unit=unit)
                raise
        
        return wrapper
    return decorator


@contextmanager
def context(name: str, log_entry: bool = True, log_exit: bool = True):
    """
    Context manager for grouping related operations.
    
    Args:
        name: Context name
        log_entry: Whether to log context entry
        log_exit: Whether to log context exit
    
    Example:
        with hawk.context("Database Migration"):
            migrate_tables()
            update_schema()
    """
    global _context_stack
    
    _context_stack.append(name)
    
    client = _get_global_client()
    
    if log_entry:
        client.info(f"Starting {name}", component="context")
    
    start_time = time.time()
    
    try:
        yield
        
        if log_exit:
            duration = time.time() - start_time
            client.success(f"Completed {name} in {duration:.2f}s", 
                          context={"duration": duration},
                          component="context")
        
    except Exception as e:
        duration = time.time() - start_time
        client.error(f"Failed {name} after {duration:.2f}s: {str(e)}",
                    context={"duration": duration, "error": str(e)},
                    component="context")
        raise
    
    finally:
        _context_stack.pop()


@contextmanager
def batch():
    """
    Context manager for batching multiple operations.
    
    Example:
        with hawk.batch():
            for item in items:
                hawk.log(f"Processing {item}")
                hawk.counter("items_processed")
    """
    client = _get_global_client()
    
    # Store original flush interval
    original_interval = client._config.flush_interval
    
    # Delay flushing during batch
    client._config.flush_interval = 60.0  # Very long interval
    
    try:
        yield
    finally:
        # Restore original interval and force flush
        client._config.flush_interval = original_interval
        client._flush_timer = None


# Utility functions
def is_available() -> bool:
    """Check if Hawk TUI is available and responsive."""
    try:
        # Simple test - try to send a message
        test_message = json.dumps({
            "jsonrpc": "2.0",
            "method": "hawk.ping",
            "params": {}
        })
        print(test_message, flush=True)
        return True
    except Exception:
        return False


def shutdown():
    """Gracefully shutdown Hawk TUI client."""
    global _global_client
    if _global_client:
        _global_client.shutdown()
        _global_client = None


# Module-level configuration
def configure(**kwargs):
    """Configure the global Hawk client."""
    global _global_client
    if _global_client:
        # Update existing client config
        for key, value in kwargs.items():
            if hasattr(_global_client._config, key):
                setattr(_global_client._config, key, value)
    else:
        # Create new client with config
        config = HawkConfig(**kwargs)
        _global_client = HawkClient(config)


# Export public API
__all__ = [
    # Layer 0: Magic mode
    'auto',
    
    # Layer 1: Simple functions
    'log', 'metric', 'config', 'progress', 'event',
    'debug', 'info', 'warn', 'error', 'success',
    'counter', 'gauge', 'histogram',
    
    # Layer 2: Decorators and context managers
    'monitor', 'timed', 'context', 'batch',
    
    # Utilities
    'is_available', 'shutdown', 'configure',
    
    # Classes for advanced usage
    'HawkClient', 'HawkConfig'
]


# Example usage when run directly
if __name__ == "__main__":
    # Demonstrate the layered API
    print("Hawk TUI Python Client - Layered API Demo", file=sys.stderr)
    
    # Layer 0: Magic mode
    auto("demo-app")
    
    # Layer 1: Simple functions
    info("Application started")
    counter("demo.runs")
    gauge("demo.value", 42)
    config("demo.enabled", True, description="Enable demo mode")
    
    # Layer 2: Decorators and context managers
    @monitor
    def demo_function():
        time.sleep(0.1)
        return "success"
    
    with context("Demo Operations"):
        result = demo_function()
        success(f"Demo completed: {result}")
    
    # Clean shutdown
    time.sleep(0.5)  # Let messages flush
    shutdown()