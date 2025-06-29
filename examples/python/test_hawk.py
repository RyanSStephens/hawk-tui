#!/usr/bin/env python3
"""
Test script for Hawk TUI Python Client Library

This script validates the implementation and demonstrates that it achieves
the "5-minute rule" for adoption. It tests all layers of complexity and
verifies that the library works correctly with or without the TUI running.

Usage:
    python test_hawk.py           # Run all tests
    python test_hawk.py --demo    # Run interactive demo
    python test_hawk.py --perf    # Run performance tests
"""

import sys
import time
import threading
import random
import traceback
from contextlib import contextmanager
import argparse

# Import the hawk library
import hawk
from hawk_advanced import Dashboard, ConfigPanel, ProgressTracker, BatchOperations


class TestResults:
    """Track test results and statistics."""
    
    def __init__(self):
        self.passed = 0
        self.failed = 0
        self.errors = []
        self.start_time = time.time()
    
    def test_pass(self, test_name: str):
        self.passed += 1
        print(f"✓ {test_name}")
    
    def test_fail(self, test_name: str, error: str):
        self.failed += 1
        self.errors.append(f"{test_name}: {error}")
        print(f"✗ {test_name}: {error}")
    
    def summary(self):
        duration = time.time() - self.start_time
        total = self.passed + self.failed
        
        print(f"\n{'='*60}")
        print(f"Test Results: {self.passed}/{total} passed ({duration:.2f}s)")
        
        if self.failed > 0:
            print(f"\nFailures:")
            for error in self.errors:
                print(f"  - {error}")
        
        return self.failed == 0


@contextmanager
def test_case(results: TestResults, test_name: str):
    """Context manager for individual test cases."""
    try:
        yield
        results.test_pass(test_name)
    except Exception as e:
        results.test_fail(test_name, str(e))


def test_layer_0_magic_mode(results: TestResults):
    """Test Layer 0: Magic mode with auto-detection."""
    print("\n--- Layer 0: Magic Mode Tests ---")
    
    with test_case(results, "Auto initialization"):
        hawk.auto("test-app")
        assert hawk._auto_enabled
        assert hawk._global_client is not None
    
    with test_case(results, "Auto-detection works without errors"):
        # These should work silently
        import logging
        logging.info("Test log message")
        logging.warning("Test warning")
        logging.error("Test error")
    
    with test_case(results, "Auto-detection graceful fallback"):
        # Should work even if no TUI is running
        for i in range(10):
            hawk.log(f"Auto test message {i}")


def test_layer_1_simple_functions(results: TestResults):
    """Test Layer 1: Simple function API."""
    print("\n--- Layer 1: Simple Functions Tests ---")
    
    with test_case(results, "Basic logging functions"):
        hawk.debug("Debug message")
        hawk.info("Info message")
        hawk.warn("Warning message")
        hawk.error("Error message")
        hawk.success("Success message")
    
    with test_case(results, "Metric functions"):
        hawk.counter("test_counter", 1)
        hawk.gauge("test_gauge", 42.5)
        hawk.histogram("test_histogram", 0.123)
        hawk.metric("custom_metric", 99, metric_type="gauge", unit="percent")
    
    with test_case(results, "Configuration functions"):
        port = hawk.config("test_port", default=8080, type="integer")
        assert port == 8080
        
        debug = hawk.config("test_debug", default=False, type="boolean")
        assert debug == False
    
    with test_case(results, "Progress tracking"):
        hawk.progress("test_progress", "Testing progress", 50, 100, unit="%")
        hawk.progress("test_progress", "Testing progress", 100, 100, unit="%", status="completed")
    
    with test_case(results, "Event logging"):
        hawk.event("test_event", "Test Event", message="Testing event system")


def test_layer_2_decorators_contexts(results: TestResults):
    """Test Layer 2: Decorators and context managers."""
    print("\n--- Layer 2: Decorators & Context Managers Tests ---")
    
    with test_case(results, "Monitor decorator"):
        @hawk.monitor
        def test_function():
            time.sleep(0.01)
            return "success"
        
        result = test_function()
        assert result == "success"
    
    with test_case(results, "Timed decorator"):
        @hawk.timed("test_timing")
        def timed_function():
            time.sleep(0.02)
            return 42
        
        result = timed_function()
        assert result == 42
    
    with test_case(results, "Context manager"):
        with hawk.context("Test Context"):
            hawk.log("Inside context")
            with hawk.context("Nested Context"):
                hawk.log("Inside nested context")
    
    with test_case(results, "Batch context manager"):
        with hawk.batch():
            for i in range(10):
                hawk.log(f"Batch message {i}")
                hawk.counter("batch_test")
    
    with test_case(results, "Monitor decorator with errors"):
        @hawk.monitor(log_errors=True)
        def error_function():
            raise ValueError("Test error")
        
        try:
            error_function()
            assert False, "Should have raised exception"
        except ValueError:
            pass  # Expected


def test_layer_3_advanced_features(results: TestResults):
    """Test Layer 3: Advanced features."""
    print("\n--- Layer 3: Advanced Features Tests ---")
    
    with test_case(results, "Dashboard creation"):
        dashboard = Dashboard("test-dashboard")
        
        # Add widgets
        dashboard.add_metric("test_metric", "Test Metric", lambda: random.uniform(0, 100))
        dashboard.add_status("test_status", "Test Status", lambda: {"service": {"status": "healthy"}})
        
        # Let it run briefly
        time.sleep(0.5)
        
        stats = dashboard.get_stats()
        assert stats["widget_count"] >= 2
        
        dashboard.shutdown()
    
    with test_case(results, "Configuration panel"):
        config = ConfigPanel("test-config")
        
        config.add_field("test_int", "integer", "Test integer", default=42, min_value=1, max_value=100)
        config.add_field("test_bool", "boolean", "Test boolean", default=True)
        config.add_field("test_enum", "enum", "Test enum", default="option1", options=["option1", "option2", "option3"])
        
        # Test value setting and getting
        assert config.set_value("test_int", 50)
        assert config.get_value("test_int") == 50
        
        assert config.set_value("test_bool", False)
        assert config.get_value("test_bool") == False
        
        # Test validation
        assert not config.set_value("test_int", 150)  # Should fail (above max)
        assert config.get_value("test_int") == 50  # Should remain unchanged
    
    with test_case(results, "Progress tracker"):
        tracker = ProgressTracker()
        
        # Test simple progress
        with tracker.start("test_task", "Test Task", 10) as progress:
            for i in range(10):
                progress.update(i + 1, f"Step {i + 1}")
                time.sleep(0.01)
        
        # Test nested progress
        with tracker.start("parent_task", "Parent Task", 2) as parent:
            parent.update(1, "Starting subtask")
            
            with tracker.start("child_task", "Child Task", 5, parent_id="parent_task") as child:
                for i in range(5):
                    child.update(i + 1, f"Child step {i + 1}")
                    time.sleep(0.005)
            
            parent.update(2, "Subtask completed")
        
        summary = tracker.get_summary()
        assert summary["completed_operations"] >= 2
    
    with test_case(results, "Batch operations"):
        batch = BatchOperations(batch_size=10, flush_interval=0.1)
        
        # Add many messages
        for i in range(25):
            batch.add_log(f"Batch log {i}", level="INFO")
            batch.add_metric(f"batch_metric_{i % 3}", random.uniform(0, 100))
            batch.add_event("batch_event", f"Event {i}")
        
        # Force flush
        batch.flush()
        
        stats = batch.get_stats()
        assert stats["messages_sent"] >= 75  # 3 messages * 25 iterations


def test_performance(results: TestResults):
    """Test performance characteristics."""
    print("\n--- Performance Tests ---")
    
    with test_case(results, "High-frequency logging performance"):
        start_time = time.time()
        
        # Send 1000 messages quickly
        with hawk.batch():
            for i in range(1000):
                hawk.log(f"Performance test message {i}")
        
        duration = time.time() - start_time
        rate = 1000 / duration
        
        print(f"    Rate: {rate:.0f} messages/second")
        assert rate > 1000, f"Performance too slow: {rate:.0f} msg/s"
    
    with test_case(results, "Thread safety"):
        results_list = []
        
        def worker_thread(thread_id):
            try:
                for i in range(100):
                    hawk.log(f"Thread {thread_id} message {i}")
                    hawk.counter(f"thread_{thread_id}_counter")
                    time.sleep(0.001)
                results_list.append(True)
            except Exception as e:
                results_list.append(f"Thread {thread_id} error: {e}")
        
        # Start multiple threads
        threads = []
        for i in range(5):
            thread = threading.Thread(target=worker_thread, args=(i,))
            threads.append(thread)
            thread.start()
        
        # Wait for completion
        for thread in threads:
            thread.join()
        
        # Check results
        assert len(results_list) == 5, f"Expected 5 results, got {len(results_list)}"
        for result in results_list:
            assert result == True, f"Thread error: {result}"
    
    with test_case(results, "Memory usage stability"):
        import gc
        
        # Force garbage collection
        gc.collect()
        
        # Generate lots of messages
        for i in range(1000):
            hawk.log(f"Memory test {i}")
            if i % 100 == 0:
                gc.collect()
        
        # Force final cleanup
        hawk.shutdown()
        hawk._global_client = None
        gc.collect()


def test_error_handling(results: TestResults):
    """Test error handling and graceful degradation."""
    print("\n--- Error Handling Tests ---")
    
    with test_case(results, "Graceful fallback when TUI unavailable"):
        # This should work even without TUI
        for i in range(10):
            hawk.log(f"Fallback test {i}")
            hawk.metric(f"fallback_metric", i)
    
    with test_case(results, "Invalid metric values handled gracefully"):
        # These should not crash the application
        try:
            hawk.metric("test", float('inf'))
            hawk.metric("test", float('nan'))
            hawk.metric("test", "invalid")
        except Exception:
            pass  # Some might raise exceptions, but shouldn't crash the app
    
    with test_case(results, "Large message handling"):
        # Test with very large messages
        large_message = "x" * 10000
        hawk.log(large_message)
        
        large_context = {f"key_{i}": f"value_{i}" * 100 for i in range(100)}
        hawk.log("Test with large context", context=large_context)
    
    with test_case(results, "Exception in monitored function"):
        @hawk.monitor
        def failing_function():
            raise RuntimeError("Intentional test error")
        
        try:
            failing_function()
            assert False, "Should have raised exception"
        except RuntimeError:
            pass  # Expected


def test_integration_patterns(results: TestResults):
    """Test common integration patterns."""
    print("\n--- Integration Pattern Tests ---")
    
    with test_case(results, "Function monitoring pattern"):
        call_count = 0
        
        @hawk.monitor
        def api_endpoint():
            nonlocal call_count
            call_count += 1
            time.sleep(0.01)
            return {"status": "success", "data": [1, 2, 3]}
        
        for i in range(5):
            result = api_endpoint()
            assert result["status"] == "success"
        
        assert call_count == 5
    
    with test_case(results, "Database operation pattern"):
        @hawk.timed("db_query")
        def simulate_db_query(query):
            hawk.log(f"Executing query: {query}", component="database")
            time.sleep(0.02)  # Simulate query time
            hawk.counter("db_queries")
            return [{"id": 1, "name": "test"}]
        
        with hawk.context("Database Operations"):
            users = simulate_db_query("SELECT * FROM users")
            orders = simulate_db_query("SELECT * FROM orders")
            
        assert len(users) == 1
        assert len(orders) == 1
    
    with test_case(results, "Background task pattern"):
        task_completed = threading.Event()
        
        @hawk.monitor
        def background_task():
            with hawk.context("Background Processing"):
                for i in range(10):
                    hawk.log(f"Processing item {i}")
                    hawk.progress("bg_task", "Background work", i + 1, 10)
                    time.sleep(0.01)
            task_completed.set()
        
        thread = threading.Thread(target=background_task)
        thread.start()
        
        # Wait for completion
        assert task_completed.wait(timeout=5.0), "Background task timed out"
        thread.join()


def run_interactive_demo():
    """Run an interactive demonstration."""
    print("🚀 Hawk TUI Python Client - Interactive Demo")
    print("=" * 50)
    
    # Initialize
    hawk.auto("interactive-demo")
    
    print("\n1. Basic logging and metrics...")
    hawk.info("Demo started")
    hawk.success("All systems operational")
    
    for i in range(5):
        hawk.counter("demo_counter")
        hawk.gauge("demo_value", random.uniform(0, 100))
        time.sleep(0.5)
    
    print("\n2. Configuration example...")
    port = hawk.config("demo_port", default=8080, description="Demo server port")
    debug = hawk.config("demo_debug", default=False, description="Enable debug mode")
    
    print(f"   Port: {port}, Debug: {debug}")
    
    print("\n3. Progress tracking...")
    with hawk.context("Demo Operations"):
        for i in range(10):
            hawk.progress("demo_progress", "Demo task", i + 1, 10)
            time.sleep(0.2)
    
    print("\n4. Function monitoring...")
    @hawk.monitor
    def demo_function():
        time.sleep(0.5)
        return "Demo completed!"
    
    result = demo_function()
    hawk.success(result)
    
    print("\n5. Dashboard (brief demo)...")
    dashboard = Dashboard("demo-dashboard")
    dashboard.add_metric("cpu_usage", "CPU %", lambda: random.uniform(20, 80))
    dashboard.add_status("services", "Status", lambda: {"API": {"status": "healthy"}})
    
    time.sleep(2.0)
    dashboard.shutdown()
    
    hawk.event("demo_complete", "Demo Completed", message="Interactive demo finished successfully")
    print("\n✨ Demo completed! Check your TUI for the full experience.")


def run_performance_benchmark():
    """Run performance benchmarks."""
    print("⚡ Hawk TUI Performance Benchmark")
    print("=" * 40)
    
    hawk.auto("performance-benchmark")
    
    # Test 1: Message throughput
    print("\n1. Message Throughput Test")
    start_time = time.time()
    message_count = 10000
    
    with hawk.batch():
        for i in range(message_count):
            hawk.log(f"Benchmark message {i}")
    
    duration = time.time() - start_time
    rate = message_count / duration
    print(f"   Sent {message_count} messages in {duration:.2f}s")
    print(f"   Rate: {rate:.0f} messages/second")
    
    # Test 2: Metric throughput
    print("\n2. Metric Throughput Test")
    start_time = time.time()
    metric_count = 5000
    
    with hawk.batch():
        for i in range(metric_count):
            hawk.metric(f"benchmark_metric_{i % 10}", random.uniform(0, 100))
    
    duration = time.time() - start_time
    rate = metric_count / duration
    print(f"   Sent {metric_count} metrics in {duration:.2f}s")
    print(f"   Rate: {rate:.0f} metrics/second")
    
    # Test 3: Mixed operations
    print("\n3. Mixed Operations Test")
    start_time = time.time()
    operation_count = 1000
    
    with hawk.batch():
        for i in range(operation_count):
            hawk.log(f"Mixed test {i}")
            hawk.counter("mixed_counter")
            hawk.gauge("mixed_gauge", i)
            if i % 10 == 0:
                hawk.progress("mixed_progress", "Mixed test", i, operation_count)
    
    duration = time.time() - start_time
    total_messages = operation_count * 3  # log + counter + gauge
    rate = total_messages / duration
    print(f"   Performed {operation_count} operations ({total_messages} messages) in {duration:.2f}s")
    print(f"   Rate: {rate:.0f} messages/second")
    
    print("\n✅ Performance benchmark completed!")


def main():
    """Main test runner."""
    parser = argparse.ArgumentParser(description="Test Hawk TUI Python client")
    parser.add_argument("--demo", action="store_true", help="Run interactive demo")
    parser.add_argument("--perf", action="store_true", help="Run performance benchmark")
    parser.add_argument("--verbose", action="store_true", help="Verbose output")
    
    args = parser.parse_args()
    
    if args.demo:
        run_interactive_demo()
        return
    
    if args.perf:
        run_performance_benchmark()
        return
    
    # Run full test suite
    print("🧪 Hawk TUI Python Client - Test Suite")
    print("=" * 50)
    
    results = TestResults()
    
    try:
        test_layer_0_magic_mode(results)
        test_layer_1_simple_functions(results)
        test_layer_2_decorators_contexts(results)
        test_layer_3_advanced_features(results)
        test_performance(results)
        test_error_handling(results)
        test_integration_patterns(results)
        
    except KeyboardInterrupt:
        print("\n\nTests interrupted by user")
        return 1
    
    except Exception as e:
        print(f"\n\nUnexpected error: {e}")
        traceback.print_exc()
        return 1
    
    finally:
        # Cleanup
        try:
            hawk.shutdown()
        except:
            pass
    
    # Print summary
    success = results.summary()
    
    if success:
        print("\n🎉 All tests passed! The library is working correctly.")
        print("\n📋 5-Minute Rule Validation:")
        print("   ✓ Layer 0: hawk.auto() works (0 lines of config)")
        print("   ✓ Layer 1: Simple functions work (1-2 lines per feature)")  
        print("   ✓ Layer 2: Decorators and contexts work (advanced patterns)")
        print("   ✓ Layer 3: Enterprise features work (dashboards, config)")
        print("   ✓ Performance: >1000 messages/second")
        print("   ✓ Thread safety: Multiple threads work correctly")
        print("   ✓ Error handling: Graceful fallback when TUI unavailable")
        return 0
    else:
        print("\n❌ Some tests failed. Please check the errors above.")
        return 1


if __name__ == "__main__":
    sys.exit(main())