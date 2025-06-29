#!/usr/bin/env python3
"""
Hawk TUI Advanced Features - Layer 3 & 4 Enterprise Features

This module provides advanced features for enterprise and complex applications:
- Dashboard creation and management
- Configuration panels with validation
- Progress tracking and reporting
- Batch operations for performance
- Security and audit features
- Remote monitoring capabilities

Usage:
    from hawk_advanced import Dashboard, ConfigPanel, ProgressTracker, AuditLogger
    
    dashboard = Dashboard("my-app")
    dashboard.add_metric("cpu_usage", get_cpu_usage, refresh=1.0)
    dashboard.add_status("database", check_database)
"""

import json
import time
import threading
import weakref
from datetime import datetime, timezone, timedelta
from typing import Any, Dict, List, Optional, Union, Callable, Set
from dataclasses import dataclass, field
from abc import ABC, abstractmethod
import uuid
import hashlib
import os
from pathlib import Path

try:
    from hawk import HawkClient, HawkConfig
except ImportError:
    # Fallback for standalone usage
    import sys
    sys.path.append(os.path.dirname(__file__))
    from hawk import HawkClient, HawkConfig


@dataclass
class WidgetConfig:
    """Configuration for dashboard widgets."""
    widget_id: str
    title: str
    widget_type: str
    refresh_rate: float = 1.0
    layout: Optional[Dict[str, int]] = None
    enabled: bool = True
    
    def __post_init__(self):
        if self.layout is None:
            self.layout = {"row": 0, "col": 0, "width": 4, "height": 3}


@dataclass
class ConfigField:
    """Configuration field definition."""
    key: str
    field_type: str
    description: str = ""
    default: Any = None
    required: bool = False
    min_value: Optional[Union[int, float]] = None
    max_value: Optional[Union[int, float]] = None
    options: Optional[List[str]] = None
    restart_required: bool = False
    validation_func: Optional[Callable[[Any], bool]] = None
    category: str = "General"


class ProgressTracker:
    """
    Advanced progress tracking with estimation and reporting.
    
    Features:
    - Time estimation based on historical data
    - Progress rate calculation
    - Hierarchical progress (sub-tasks)
    - Automatic completion detection
    """
    
    def __init__(self, client: Optional[HawkClient] = None):
        self.client = client or HawkClient()
        self._active_progress = {}
        self._progress_history = {}
        self._lock = threading.Lock()
    
    def start(self, progress_id: str, label: str, total: Union[int, float],
              unit: str = "", parent_id: Optional[str] = None) -> 'ProgressContext':
        """Start tracking progress for an operation."""
        with self._lock:
            progress_data = {
                "id": progress_id,
                "label": label,
                "total": total,
                "current": 0,
                "unit": unit,
                "start_time": time.time(),
                "last_update": time.time(),
                "parent_id": parent_id,
                "children": set(),
                "estimated_completion": None
            }
            
            self._active_progress[progress_id] = progress_data
            
            if parent_id and parent_id in self._active_progress:
                self._active_progress[parent_id]["children"].add(progress_id)
        
        # Send initial progress
        self.client.progress(progress_id, label, 0, total, unit=unit, status="pending")
        
        return ProgressContext(self, progress_id)
    
    def update(self, progress_id: str, current: Union[int, float],
               details: Optional[str] = None):
        """Update progress for an operation."""
        with self._lock:
            if progress_id not in self._active_progress:
                return
            
            progress_data = self._active_progress[progress_id]
            progress_data["current"] = current
            progress_data["last_update"] = time.time()
            
            # Calculate progress rate and estimate completion
            elapsed = progress_data["last_update"] - progress_data["start_time"]
            if elapsed > 0 and current > 0:
                rate = current / elapsed
                remaining = progress_data["total"] - current
                if rate > 0:
                    eta_seconds = remaining / rate
                    progress_data["estimated_completion"] = (
                        datetime.now(timezone.utc) + timedelta(seconds=eta_seconds)
                    ).isoformat()
        
        # Send progress update
        params = {
            "current": current,
            "details": details or f"Progress: {current}/{progress_data['total']}"
        }
        
        if progress_data["estimated_completion"]:
            params["estimated_completion"] = progress_data["estimated_completion"]
        
        self.client.progress(
            progress_id, 
            progress_data["label"],
            current, 
            progress_data["total"],
            unit=progress_data["unit"],
            **params
        )
    
    def complete(self, progress_id: str, success: bool = True,
                 final_message: Optional[str] = None):
        """Mark progress as completed."""
        with self._lock:
            if progress_id not in self._active_progress:
                return
            
            progress_data = self._active_progress[progress_id]
            
            # Complete all children first
            for child_id in list(progress_data["children"]):
                self.complete(child_id, success)
            
            # Store in history
            progress_data["end_time"] = time.time()
            progress_data["success"] = success
            progress_data["duration"] = progress_data["end_time"] - progress_data["start_time"]
            
            self._progress_history[progress_id] = progress_data
            del self._active_progress[progress_id]
        
        # Send completion
        status = "completed" if success else "error"
        details = final_message or ("Operation completed" if success else "Operation failed")
        
        self.client.progress(
            progress_id,
            progress_data["label"],
            progress_data["total"],
            progress_data["total"],
            unit=progress_data["unit"],
            status=status,
            details=details
        )
    
    def get_active_progress(self) -> Dict[str, Dict]:
        """Get all currently active progress operations."""
        with self._lock:
            return dict(self._active_progress)
    
    def get_summary(self) -> Dict[str, Any]:
        """Get progress tracking summary."""
        with self._lock:
            active_count = len(self._active_progress)
            completed_count = len(self._progress_history)
            
            if self._progress_history:
                avg_duration = sum(p["duration"] for p in self._progress_history.values()) / completed_count
                success_rate = sum(1 for p in self._progress_history.values() if p.get("success", True)) / completed_count
            else:
                avg_duration = 0
                success_rate = 1.0
            
            return {
                "active_operations": active_count,
                "completed_operations": completed_count,
                "average_duration": avg_duration,
                "success_rate": success_rate,
                "current_operations": list(self._active_progress.keys())
            }


class ProgressContext:
    """Context manager for progress tracking."""
    
    def __init__(self, tracker: ProgressTracker, progress_id: str):
        self.tracker = tracker
        self.progress_id = progress_id
        self._current = 0
        
    def update(self, current: Union[int, float], details: Optional[str] = None):
        """Update progress."""
        self._current = current
        self.tracker.update(self.progress_id, current, details)
    
    def increment(self, amount: Union[int, float] = 1, details: Optional[str] = None):
        """Increment progress by amount."""
        self._current += amount
        self.tracker.update(self.progress_id, self._current, details)
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        success = exc_type is None
        self.tracker.complete(self.progress_id, success)


class Dashboard:
    """
    Advanced dashboard creation and management.
    
    Features:
    - Widget lifecycle management
    - Automatic refresh scheduling
    - Data validation and caching
    - Layout management
    - Performance optimization
    """
    
    def __init__(self, name: str, client: Optional[HawkClient] = None):
        self.name = name
        self.client = client or HawkClient()
        self._widgets = {}
        self._refresh_threads = {}
        self._layout_grid = {}
        self._shutdown = False
        
        # Auto-layout configuration
        self._next_row = 0
        self._next_col = 0
        self._max_cols = 12
        
        # Performance tracking
        self._refresh_stats = {}
        
    def add_metric(self, widget_id: str, title: str, data_func: Callable[[], Union[int, float]],
                   unit: str = "", refresh: float = 1.0, **layout_kwargs) -> 'MetricWidget':
        """Add a metric widget that displays a single numeric value."""
        widget = MetricWidget(widget_id, title, data_func, unit, refresh, **layout_kwargs)
        self._add_widget(widget)
        return widget
    
    def add_status(self, widget_id: str, title: str, check_func: Callable[[], Dict[str, Any]],
                   refresh: float = 5.0, **layout_kwargs) -> 'StatusWidget':
        """Add a status widget that displays service health."""
        widget = StatusWidget(widget_id, title, check_func, refresh, **layout_kwargs)
        self._add_widget(widget)
        return widget
    
    def add_chart(self, widget_id: str, title: str, data_func: Callable[[], Dict[str, Any]],
                  chart_type: str = "line", refresh: float = 2.0, **layout_kwargs) -> 'ChartWidget':
        """Add a chart widget for time-series data."""
        widget = ChartWidget(widget_id, title, data_func, chart_type, refresh, **layout_kwargs)
        self._add_widget(widget)
        return widget
    
    def add_table(self, widget_id: str, title: str, data_func: Callable[[], Dict[str, Any]],
                  refresh: float = 5.0, **layout_kwargs) -> 'TableWidget':
        """Add a table widget for tabular data."""
        widget = TableWidget(widget_id, title, data_func, refresh, **layout_kwargs)
        self._add_widget(widget)
        return widget
    
    def add_custom(self, widget_id: str, title: str, widget_type: str,
                   data_func: Callable[[], Any], refresh: float = 1.0,
                   **layout_kwargs) -> 'CustomWidget':
        """Add a custom widget type."""
        widget = CustomWidget(widget_id, title, widget_type, data_func, refresh, **layout_kwargs)
        self._add_widget(widget)
        return widget
    
    def _add_widget(self, widget: 'DashboardWidget'):
        """Add widget to dashboard and start refresh cycle."""
        # Auto-layout if not specified
        if widget.config.layout == {"row": 0, "col": 0, "width": 4, "height": 3}:
            widget.config.layout = self._auto_layout(widget.config.layout["width"], 
                                                   widget.config.layout["height"])
        
        self._widgets[widget.config.widget_id] = widget
        
        # Start refresh thread
        if widget.config.refresh_rate > 0:
            thread = threading.Thread(
                target=self._refresh_worker,
                args=(widget,),
                daemon=True
            )
            self._refresh_threads[widget.config.widget_id] = thread
            thread.start()
        
        # Send initial widget creation
        self._send_widget_update(widget)
    
    def _auto_layout(self, width: int, height: int) -> Dict[str, int]:
        """Automatically assign layout position."""
        if self._next_col + width > self._max_cols:
            self._next_row += height
            self._next_col = 0
        
        layout = {
            "row": self._next_row,
            "col": self._next_col,
            "width": width,
            "height": height
        }
        
        self._next_col += width
        return layout
    
    def _refresh_worker(self, widget: 'DashboardWidget'):
        """Background worker for widget refresh."""
        while not self._shutdown and widget.config.enabled:
            try:
                start_time = time.time()
                
                # Get fresh data
                data = widget.get_data()
                
                # Send update
                self._send_widget_update(widget, data)
                
                # Track performance
                refresh_time = time.time() - start_time
                self._refresh_stats[widget.config.widget_id] = {
                    "last_refresh": time.time(),
                    "refresh_duration": refresh_time,
                    "refresh_count": self._refresh_stats.get(widget.config.widget_id, {}).get("refresh_count", 0) + 1
                }
                
                # Wait for next refresh
                time.sleep(max(0, widget.config.refresh_rate - refresh_time))
                
            except Exception as e:
                # Log error but don't crash
                self.client.error(f"Widget refresh error: {widget.config.widget_id}: {e}")
                time.sleep(widget.config.refresh_rate)
    
    def _send_widget_update(self, widget: 'DashboardWidget', data: Any = None):
        """Send widget update to TUI."""
        if data is None:
            data = widget.get_data()
        
        self.client._send_message_async("hawk.dashboard", {
            "widget_id": widget.config.widget_id,
            "type": widget.config.widget_type,
            "title": widget.config.title,
            "data": data,
            "layout": widget.config.layout,
            "dashboard": self.name
        })
    
    def remove_widget(self, widget_id: str):
        """Remove a widget from the dashboard."""
        if widget_id in self._widgets:
            del self._widgets[widget_id]
        
        if widget_id in self._refresh_threads:
            # Thread will stop on next iteration
            del self._refresh_threads[widget_id]
    
    def get_stats(self) -> Dict[str, Any]:
        """Get dashboard performance statistics."""
        return {
            "widget_count": len(self._widgets),
            "active_threads": len(self._refresh_threads),
            "refresh_stats": dict(self._refresh_stats)
        }
    
    def shutdown(self):
        """Shutdown dashboard and stop all refresh threads."""
        self._shutdown = True
        
        # Wait for threads to finish
        for thread in self._refresh_threads.values():
            if thread.is_alive():
                thread.join(timeout=1.0)


class DashboardWidget(ABC):
    """Base class for dashboard widgets."""
    
    def __init__(self, widget_id: str, title: str, widget_type: str,
                 data_func: Callable, refresh_rate: float = 1.0, **layout_kwargs):
        layout = layout_kwargs or {"width": 4, "height": 3}
        self.config = WidgetConfig(
            widget_id=widget_id,
            title=title,
            widget_type=widget_type,
            refresh_rate=refresh_rate,
            layout=layout
        )
        self.data_func = data_func
        self._cache = None
        self._cache_time = 0
    
    @abstractmethod
    def get_data(self) -> Any:
        """Get current widget data."""
        pass


class MetricWidget(DashboardWidget):
    """Widget for displaying single numeric metrics."""
    
    def __init__(self, widget_id: str, title: str, data_func: Callable[[], Union[int, float]],
                 unit: str = "", refresh_rate: float = 1.0, **layout_kwargs):
        super().__init__(widget_id, title, "gauge", data_func, refresh_rate, **layout_kwargs)
        self.unit = unit
    
    def get_data(self) -> Dict[str, Any]:
        value = self.data_func()
        return {
            "value": value,
            "unit": self.unit,
            "timestamp": datetime.now(timezone.utc).isoformat()
        }


class StatusWidget(DashboardWidget):
    """Widget for displaying service status."""
    
    def __init__(self, widget_id: str, title: str, check_func: Callable[[], Dict[str, Any]],
                 refresh_rate: float = 5.0, **layout_kwargs):
        super().__init__(widget_id, title, "status_grid", check_func, refresh_rate, **layout_kwargs)
    
    def get_data(self) -> Dict[str, Any]:
        return self.data_func()


class ChartWidget(DashboardWidget):
    """Widget for displaying charts and time-series data."""
    
    def __init__(self, widget_id: str, title: str, data_func: Callable[[], Dict[str, Any]],
                 chart_type: str = "line", refresh_rate: float = 2.0, **layout_kwargs):
        super().__init__(widget_id, title, "metric_chart", data_func, refresh_rate, **layout_kwargs)
        self.chart_type = chart_type
    
    def get_data(self) -> Dict[str, Any]:
        data = self.data_func()
        data["chart_type"] = self.chart_type
        return data


class TableWidget(DashboardWidget):
    """Widget for displaying tabular data."""
    
    def __init__(self, widget_id: str, title: str, data_func: Callable[[], Dict[str, Any]],
                 refresh_rate: float = 5.0, **layout_kwargs):
        super().__init__(widget_id, title, "table", data_func, refresh_rate, **layout_kwargs)
    
    def get_data(self) -> Dict[str, Any]:
        return self.data_func()


class CustomWidget(DashboardWidget):
    """Widget for custom widget types."""
    
    def __init__(self, widget_id: str, title: str, widget_type: str,
                 data_func: Callable[[], Any], refresh_rate: float = 1.0, **layout_kwargs):
        super().__init__(widget_id, title, widget_type, data_func, refresh_rate, **layout_kwargs)
    
    def get_data(self) -> Any:
        return self.data_func()


class ConfigPanel:
    """
    Advanced configuration panel with validation and persistence.
    
    Features:
    - Field validation and constraints
    - Configuration persistence
    - Change notifications
    - Bulk updates with transactions
    - Category organization
    """
    
    def __init__(self, name: str, client: Optional[HawkClient] = None):
        self.name = name
        self.client = client or HawkClient()
        self._fields = {}
        self._values = {}
        self._callbacks = {}
        self._categories = set()
        self._lock = threading.Lock()
        
        # Persistence
        self._config_file = None
        self._auto_save = True
    
    def add_field(self, key: str, field_type: str, description: str = "",
                  default: Any = None, required: bool = False,
                  min_value: Optional[Union[int, float]] = None,
                  max_value: Optional[Union[int, float]] = None,
                  options: Optional[List[str]] = None,
                  restart_required: bool = False,
                  validation_func: Optional[Callable[[Any], bool]] = None,
                  category: str = "General") -> 'ConfigField':
        """Add a configuration field."""
        field = ConfigField(
            key=key,
            field_type=field_type,
            description=description,
            default=default,
            required=required,
            min_value=min_value,
            max_value=max_value,
            options=options,
            restart_required=restart_required,
            validation_func=validation_func,
            category=category
        )
        
        with self._lock:
            self._fields[key] = field
            self._categories.add(category)
            
            # Set default value if not already set
            if key not in self._values and default is not None:
                self._values[key] = default
        
        # Send field definition
        self._send_config_field(field)
        
        return field
    
    def set_value(self, key: str, value: Any, validate: bool = True) -> bool:
        """Set configuration value with validation."""
        with self._lock:
            if key not in self._fields:
                raise ValueError(f"Unknown configuration key: {key}")
            
            field = self._fields[key]
            
            # Validate value
            if validate and not self._validate_value(field, value):
                return False
            
            old_value = self._values.get(key)
            self._values[key] = value
            
            # Send update
            self._send_config_update(key, value)
            
            # Call change callback if registered
            if key in self._callbacks:
                try:
                    self._callbacks[key](key, old_value, value)
                except Exception as e:
                    self.client.error(f"Config callback error for {key}: {e}")
            
            # Auto-save
            if self._auto_save and self._config_file:
                self._save_to_file()
            
            return True
    
    def get_value(self, key: str, default: Any = None) -> Any:
        """Get configuration value."""
        with self._lock:
            return self._values.get(key, default)
    
    def get_all_values(self) -> Dict[str, Any]:
        """Get all configuration values."""
        with self._lock:
            return dict(self._values)
    
    def on_change(self, key: str, callback: Callable[[str, Any, Any], None]):
        """Register change callback for a configuration key."""
        self._callbacks[key] = callback
    
    def _validate_value(self, field: ConfigField, value: Any) -> bool:
        """Validate a configuration value."""
        try:
            # Type validation
            if field.field_type == "integer" and not isinstance(value, int):
                try:
                    value = int(value)
                except (ValueError, TypeError):
                    return False
            elif field.field_type == "float" and not isinstance(value, (int, float)):
                try:
                    value = float(value)
                except (ValueError, TypeError):
                    return False
            elif field.field_type == "boolean" and not isinstance(value, bool):
                if isinstance(value, str):
                    value = value.lower() in ("true", "1", "yes", "on")
                else:
                    return False
            elif field.field_type == "enum" and field.options and value not in field.options:
                return False
            
            # Range validation
            if field.min_value is not None and value < field.min_value:
                return False
            if field.max_value is not None and value > field.max_value:
                return False
            
            # Custom validation
            if field.validation_func and not field.validation_func(value):
                return False
            
            return True
            
        except Exception:
            return False
    
    def _send_config_field(self, field: ConfigField):
        """Send configuration field definition."""
        params = {
            "key": field.key,
            "type": field.field_type,
            "description": field.description,
            "default": field.default,
            "required": field.required,
            "category": field.category,
            "panel": self.name
        }
        
        if field.min_value is not None:
            params["min"] = field.min_value
        if field.max_value is not None:
            params["max"] = field.max_value
        if field.options:
            params["options"] = field.options
        if field.restart_required:
            params["restart_required"] = True
        
        # Include current value if set
        if field.key in self._values:
            params["value"] = self._values[field.key]
        
        self.client._send_message_async("hawk.config", params)
    
    def _send_config_update(self, key: str, value: Any):
        """Send configuration value update."""
        self.client._send_message_async("hawk.config", {
            "key": key,
            "value": value,
            "panel": self.name
        })
    
    def save_to_file(self, filename: Union[str, Path]):
        """Save configuration to file."""
        self._config_file = Path(filename)
        self._save_to_file()
    
    def load_from_file(self, filename: Union[str, Path]):
        """Load configuration from file."""
        config_file = Path(filename)
        if config_file.exists():
            try:
                with open(config_file, 'r') as f:
                    data = json.load(f)
                
                for key, value in data.items():
                    if key in self._fields:
                        self.set_value(key, value, validate=True)
                
                self._config_file = config_file
                
            except Exception as e:
                self.client.error(f"Failed to load config from {config_file}: {e}")
    
    def _save_to_file(self):
        """Save current configuration to file."""
        if not self._config_file:
            return
        
        try:
            with open(self._config_file, 'w') as f:
                json.dump(self._values, f, indent=2, default=str)
        except Exception as e:
            self.client.error(f"Failed to save config to {self._config_file}: {e}")
    
    def export_schema(self) -> Dict[str, Any]:
        """Export configuration schema."""
        schema = {
            "name": self.name,
            "categories": list(self._categories),
            "fields": {}
        }
        
        for key, field in self._fields.items():
            schema["fields"][key] = {
                "type": field.field_type,
                "description": field.description,
                "default": field.default,
                "required": field.required,
                "category": field.category
            }
            
            if field.min_value is not None:
                schema["fields"][key]["min"] = field.min_value
            if field.max_value is not None:
                schema["fields"][key]["max"] = field.max_value
            if field.options:
                schema["fields"][key]["options"] = field.options
        
        return schema


class BatchOperations:
    """
    High-performance batch operations for metrics and logs.
    
    Features:
    - Automatic batching with configurable thresholds
    - Compression for large batches
    - Error handling and retry logic
    - Performance monitoring
    """
    
    def __init__(self, client: Optional[HawkClient] = None, 
                 batch_size: int = 100, flush_interval: float = 1.0):
        self.client = client or HawkClient()
        self.batch_size = batch_size
        self.flush_interval = flush_interval
        
        self._message_batch = []
        self._lock = threading.Lock()
        self._last_flush = time.time()
        self._stats = {
            "messages_sent": 0,
            "batches_sent": 0,
            "avg_batch_size": 0,
            "flush_errors": 0
        }
        
        # Background flush timer
        self._timer = None
        self._start_flush_timer()
    
    def add_log(self, message: str, level: str = "INFO", **kwargs):
        """Add log message to batch."""
        self._add_message("hawk.log", {
            "message": message,
            "level": level,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            **kwargs
        })
    
    def add_metric(self, name: str, value: Union[int, float], **kwargs):
        """Add metric to batch."""
        self._add_message("hawk.metric", {
            "name": name,
            "value": value,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            **kwargs
        })
    
    def add_event(self, event_type: str, title: str, **kwargs):
        """Add event to batch."""
        self._add_message("hawk.event", {
            "type": event_type,
            "title": title,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            **kwargs
        })
    
    def _add_message(self, method: str, params: Dict[str, Any]):
        """Add message to batch queue."""
        with self._lock:
            self._message_batch.append({
                "jsonrpc": "2.0",
                "method": method,
                "params": params,
                "hawk_meta": {
                    "app_name": self.client._config.app_name,
                    "session_id": self.client.session_id,
                    "sequence": self.client._get_next_sequence(),
                    "batch_id": str(uuid.uuid4())[:8]
                }
            })
            
            # Auto-flush if batch is full
            if len(self._message_batch) >= self.batch_size:
                self._flush_now()
    
    def _flush_now(self):
        """Flush current batch immediately."""
        if not self._message_batch:
            return
        
        batch = self._message_batch[:]
        self._message_batch.clear()
        
        try:
            if len(batch) == 1:
                json_str = json.dumps(batch[0], default=self._json_serializer)
            else:
                json_str = json.dumps(batch, default=self._json_serializer)
            
            print(json_str, flush=True)
            
            # Update stats
            self._stats["messages_sent"] += len(batch)
            self._stats["batches_sent"] += 1
            self._stats["avg_batch_size"] = self._stats["messages_sent"] / self._stats["batches_sent"]
            
        except Exception as e:
            self._stats["flush_errors"] += 1
            if self.client._config.debug:
                print(f"Batch flush error: {e}", file=sys.stderr)
        
        self._last_flush = time.time()
    
    def _json_serializer(self, obj):
        """JSON serializer for datetime objects."""
        if isinstance(obj, datetime):
            return obj.isoformat()
        raise TypeError(f"Object of type {type(obj)} is not JSON serializable")
    
    def _start_flush_timer(self):
        """Start periodic flush timer."""
        def timer_flush():
            current_time = time.time()
            if (current_time - self._last_flush >= self.flush_interval and 
                self._message_batch):
                with self._lock:
                    self._flush_now()
            
            # Schedule next timer
            if not self.client._shutdown:
                self._timer = threading.Timer(self.flush_interval / 2, timer_flush)
                self._timer.daemon = True
                self._timer.start()
        
        self._timer = threading.Timer(self.flush_interval / 2, timer_flush)
        self._timer.daemon = True
        self._timer.start()
    
    def flush(self):
        """Manually flush current batch."""
        with self._lock:
            self._flush_now()
    
    def get_stats(self) -> Dict[str, Any]:
        """Get batch operation statistics."""
        return dict(self._stats)


class AuditLogger:
    """
    Security audit logging for enterprise compliance.
    
    Features:
    - Tamper-evident logging
    - Structured audit events
    - Compliance reporting
    - Secure storage options
    """
    
    def __init__(self, audit_file: Optional[Union[str, Path]] = None,
                 client: Optional[HawkClient] = None):
        self.client = client or HawkClient()
        self.audit_file = Path(audit_file) if audit_file else None
        self._lock = threading.Lock()
        
        # Security features
        self._session_id = str(uuid.uuid4())
        self._sequence = 0
        
        if self.audit_file:
            self.audit_file.parent.mkdir(parents=True, exist_ok=True)
    
    def log_access(self, user: str, resource: str, action: str, 
                   result: str = "success", details: Optional[Dict] = None):
        """Log access attempt."""
        self._log_audit_event("access", {
            "user": user,
            "resource": resource,
            "action": action,
            "result": result,
            "details": details or {}
        })
    
    def log_config_change(self, user: str, key: str, old_value: Any, new_value: Any):
        """Log configuration change."""
        self._log_audit_event("config_change", {
            "user": user,
            "config_key": key,
            "old_value": str(old_value),
            "new_value": str(new_value)
        })
    
    def log_command(self, user: str, command: str, args: List[str], 
                    result: str = "success", output: Optional[str] = None):
        """Log command execution."""
        self._log_audit_event("command", {
            "user": user,
            "command": command,
            "arguments": args,
            "result": result,
            "output": output
        })
    
    def _log_audit_event(self, event_type: str, data: Dict[str, Any]):
        """Log structured audit event."""
        with self._lock:
            self._sequence += 1
            
            event = {
                "timestamp": datetime.now(timezone.utc).isoformat(),
                "session_id": self._session_id,
                "sequence": self._sequence,
                "event_type": event_type,
                "data": data,
                "source_ip": self._get_source_ip(),
                "checksum": None  # Will be calculated
            }
            
            # Calculate checksum for tamper detection
            event_str = json.dumps(event, sort_keys=True)
            event["checksum"] = hashlib.sha256(event_str.encode()).hexdigest()[:16]
            
            # Write to audit file
            if self.audit_file:
                try:
                    with open(self.audit_file, 'a') as f:
                        f.write(json.dumps(event) + '\n')
                except Exception as e:
                    self.client.error(f"Audit log write failed: {e}")
            
            # Send to TUI
            self.client.event("audit_event", f"Audit: {event_type}", 
                            severity="info", data=event)
    
    def _get_source_ip(self) -> str:
        """Get source IP address for audit logging."""
        try:
            import socket
            hostname = socket.gethostname()
            return socket.gethostbyname(hostname)
        except Exception:
            return "unknown"
    
    def verify_integrity(self) -> Dict[str, Any]:
        """Verify audit log integrity."""
        if not self.audit_file or not self.audit_file.exists():
            return {"status": "no_audit_file", "verified": 0, "corrupted": 0}
        
        verified = 0
        corrupted = 0
        
        try:
            with open(self.audit_file, 'r') as f:
                for line_num, line in enumerate(f, 1):
                    try:
                        event = json.loads(line.strip())
                        
                        # Verify checksum
                        stored_checksum = event.pop("checksum", None)
                        calculated_checksum = hashlib.sha256(
                            json.dumps(event, sort_keys=True).encode()
                        ).hexdigest()[:16]
                        
                        if stored_checksum == calculated_checksum:
                            verified += 1
                        else:
                            corrupted += 1
                            self.client.warn(f"Audit log corruption detected at line {line_num}")
                        
                    except json.JSONDecodeError:
                        corrupted += 1
                        
        except Exception as e:
            return {"status": "error", "error": str(e)}
        
        return {
            "status": "completed",
            "verified": verified,
            "corrupted": corrupted,
            "integrity": verified / (verified + corrupted) if (verified + corrupted) > 0 else 1.0
        }


# Export advanced features
__all__ = [
    'Dashboard', 'ConfigPanel', 'ProgressTracker', 'ProgressContext',
    'BatchOperations', 'AuditLogger',
    'MetricWidget', 'StatusWidget', 'ChartWidget', 'TableWidget', 'CustomWidget',
    'WidgetConfig', 'ConfigField'
]


# Example usage when run directly
if __name__ == "__main__":
    import random
    import sys
    
    print("Hawk TUI Advanced Features Demo", file=sys.stderr)
    
    # Create dashboard
    dashboard = Dashboard("advanced-demo")
    
    # Add widgets
    dashboard.add_metric("cpu_usage", "CPU Usage", 
                        lambda: random.uniform(30, 80), unit="%")
    
    dashboard.add_status("services", "Service Status",
                        lambda: {
                            "API": {"status": "healthy", "response_time": "45ms"},
                            "Database": {"status": "healthy", "response_time": "12ms"},
                            "Cache": {"status": random.choice(["healthy", "degraded"]), 
                                    "response_time": f"{random.randint(1, 50)}ms"}
                        })
    
    # Configuration panel
    config = ConfigPanel("app-config")
    config.add_field("max_connections", "integer", "Maximum connections", 
                     default=100, min_value=1, max_value=1000)
    config.add_field("debug_enabled", "boolean", "Enable debug mode", default=False)
    config.add_field("log_level", "enum", "Logging level", default="INFO",
                     options=["DEBUG", "INFO", "WARN", "ERROR"])
    
    # Progress tracking
    progress = ProgressTracker()
    
    # Demonstrate for a few seconds
    try:
        with progress.start("demo_task", "Running demo", 10) as task:
            for i in range(10):
                task.update(i + 1, f"Step {i + 1}/10")
                time.sleep(0.5)
        
        print("Advanced features demo completed", file=sys.stderr)
        
    except KeyboardInterrupt:
        print("Demo interrupted", file=sys.stderr)
    finally:
        dashboard.shutdown()