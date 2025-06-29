# Real-World Examples & Use Cases

## 1. Web API Server Monitoring

### Basic Python Flask App
```python
from flask import Flask
import hawk

app = Flask(__name__)
hawk.auto()  # One line - auto-detects HTTP requests, logs, errors

@app.route('/api/users')
def get_users():
    hawk.log("Fetching users from database")
    users = db.get_users()
    hawk.metric("api_requests", tags={"endpoint": "/users", "method": "GET"})
    return users

if __name__ == '__main__':
    app.run()
```

**TUI Shows:**
- Real-time HTTP request logs
- Response time metrics
- Error rate charts
- Active connection count

### Advanced Node.js Express App
```javascript
const express = require('express');
const hawk = require('hawk-tui');

const app = express();

// Setup monitoring dashboard
const tui = new hawk.TUI("E-commerce API");
const dashboard = tui.dashboard("API Overview");

dashboard.add_metric("Requests/sec", () => getRequestRate());
dashboard.add_metric("Response Time", () => getAvgResponseTime());
dashboard.add_chart("Error Rates", () => getErrorRates());
dashboard.add_table("Active Sessions", () => getActiveSessions());

// Middleware for automatic request tracking
app.use(hawk.middleware({
    track_requests: true,
    track_errors: true,
    track_performance: true
}));

app.get('/api/products', async (req, res) => {
    const timer = hawk.timer("product_query");
    
    try {
        const products = await db.products.findAll();
        hawk.counter("successful_requests").inc();
        res.json(products);
    } catch (error) {
        hawk.log(error.message, {level: "ERROR", context: "product_query"});
        hawk.counter("failed_requests", {endpoint: "products"}).inc();
        res.status(500).json({error: "Internal server error"});
    } finally {
        timer.stop();
    }
});

app.listen(3000, () => {
    hawk.log("Server started on port 3000", {level: "SUCCESS"});
});
```

## 2. Database Migration Tool

### Python Migration Script
```python
import hawk
from sqlalchemy import create_engine

# Setup migration monitoring
tui = hawk.TUI("Database Migration")
migration_panel = tui.panel("Migration Progress")

def migrate_table(table_name, records_count):
    """Migrate a single table with progress tracking"""
    progress = migration_panel.progress(f"Migrating {table_name}", total=records_count)
    
    for i, record in enumerate(get_records(table_name)):
        try:
            migrate_record(record)
            progress.update(i + 1)
            
            if i % 100 == 0:  # Log every 100 records
                hawk.log(f"Migrated {i+1}/{records_count} records from {table_name}")
                
        except Exception as e:
            hawk.log(f"Failed to migrate record {record.id}: {e}", level="ERROR")
            hawk.counter("migration_errors", tags={"table": table_name}).inc()
    
    progress.complete()
    hawk.log(f"✅ Completed migration of {table_name}", level="SUCCESS")

def main():
    tables = [
        ("users", 10000),
        ("orders", 50000), 
        ("products", 5000),
        ("reviews", 25000)
    ]
    
    for table_name, count in tables:
        migrate_table(table_name, count)
    
    hawk.log("🎉 All migrations completed!", level="SUCCESS")

if __name__ == "__main__":
    main()
```

**TUI Shows:**
- Progress bars for each table
- Real-time error logs
- Migration speed metrics
- Memory usage monitoring

## 3. Data Processing Pipeline

### Python ETL Job
```python
import hawk
import pandas as pd
from datetime import datetime

class DataPipeline:
    def __init__(self):
        self.tui = hawk.TUI("Data Processing Pipeline")
        self.setup_dashboard()
    
    def setup_dashboard(self):
        # Main dashboard
        dash = self.tui.dashboard("Pipeline Status")
        dash.add_metric("Records Processed", lambda: self.records_processed)
        dash.add_metric("Processing Rate", lambda: self.processing_rate)
        dash.add_chart("Quality Score", lambda: self.get_quality_metrics())
        
        # Configuration panel
        config = self.tui.config_panel("Settings")
        config.add_field("batch_size", type="int", default=1000, min=100, max=10000)
        config.add_field("quality_threshold", type="float", default=0.95, min=0.0, max=1.0)
        config.add_field("retry_attempts", type="int", default=3, min=1, max=10)
        
        # Control panel
        controls = self.tui.command_panel("Controls")
        controls.add_command("Pause Processing", self.pause_processing)
        controls.add_command("Resume Processing", self.resume_processing)
        controls.add_command("Export Results", self.export_results)
    
    def process_file(self, filepath):
        """Process a single data file"""
        hawk.log(f"Starting processing of {filepath}")
        
        # Load data
        with hawk.context("Data Loading"):
            df = pd.read_csv(filepath)
            hawk.metric("input_records", len(df))
        
        # Clean data
        with hawk.context("Data Cleaning"):
            df_clean = self.clean_data(df)
            dropped = len(df) - len(df_clean)
            if dropped > 0:
                hawk.log(f"Dropped {dropped} invalid records", level="WARN")
        
        # Transform data
        with hawk.context("Data Transformation"):
            df_transformed = self.transform_data(df_clean)
            
        # Quality check
        quality_score = self.check_quality(df_transformed)
        hawk.metric("quality_score", quality_score)
        
        if quality_score < hawk.config("quality_threshold"):
            hawk.log(f"Quality score {quality_score:.2f} below threshold", level="ERROR")
            return False
        
        # Save results
        output_path = f"processed_{filepath}"
        df_transformed.to_csv(output_path, index=False)
        hawk.log(f"✅ Saved processed data to {output_path}", level="SUCCESS")
        
        return True
    
    def run_pipeline(self, input_files):
        """Run the complete pipeline"""
        hawk.log("🚀 Starting data pipeline", level="SUCCESS")
        
        total_files = len(input_files)
        progress = hawk.progress("Overall Progress", total=total_files)
        
        for i, filepath in enumerate(input_files):
            success = self.process_file(filepath)
            progress.update(i + 1)
            
            if success:
                hawk.counter("files_processed").inc()
            else:
                hawk.counter("files_failed").inc()
        
        hawk.log("🎉 Pipeline completed!", level="SUCCESS")
```

## 4. DevOps Deployment Monitor

### Kubernetes Deployment Script
```python
import hawk
import subprocess
import yaml
from kubernetes import client, config

class DeploymentMonitor:
    def __init__(self):
        self.tui = hawk.TUI("Kubernetes Deployment")
        self.k8s_client = client.AppsV1Api()
        self.setup_monitoring()
    
    def setup_monitoring(self):
        # Service status dashboard
        dash = self.tui.dashboard("Cluster Status")
        dash.add_table("Deployments", self.get_deployment_status)
        dash.add_table("Pods", self.get_pod_status)
        dash.add_metric("Healthy Pods", self.count_healthy_pods)
        
        # Deployment controls
        controls = self.tui.command_panel("Deployment Controls")
        controls.add_command("Deploy Latest", self.deploy_latest)
        controls.add_command("Rollback", self.rollback_deployment)
        controls.add_command("Scale Up", lambda: self.scale_deployment(replicas=5))
        controls.add_command("Scale Down", lambda: self.scale_deployment(replicas=2))
    
    def deploy_application(self, app_name, image_tag):
        """Deploy application with real-time monitoring"""
        hawk.log(f"🚀 Starting deployment of {app_name}:{image_tag}")
        
        # Update deployment
        with hawk.context("Updating Deployment"):
            self.update_deployment(app_name, image_tag)
        
        # Monitor rollout
        with hawk.context("Monitoring Rollout"):
            success = self.wait_for_rollout(app_name)
        
        if success:
            hawk.log(f"✅ Deployment of {app_name} successful!", level="SUCCESS")
            hawk.metric("successful_deployments").inc()
        else:
            hawk.log(f"❌ Deployment of {app_name} failed!", level="ERROR")
            hawk.metric("failed_deployments").inc()
            
        return success
    
    def wait_for_rollout(self, app_name, timeout=300):
        """Wait for deployment rollout with progress tracking"""
        progress = hawk.progress(f"Rolling out {app_name}", indeterminate=True)
        
        for attempt in range(timeout):
            deployment = self.k8s_client.read_namespaced_deployment(
                name=app_name, namespace="default"
            )
            
            ready_replicas = deployment.status.ready_replicas or 0
            desired_replicas = deployment.spec.replicas
            
            hawk.log(f"Pods ready: {ready_replicas}/{desired_replicas}")
            
            if ready_replicas == desired_replicas:
                progress.complete()
                return True
                
            time.sleep(1)
        
        progress.fail()
        return False
```

## 5. Machine Learning Training Monitor

### PyTorch Training Script
```python
import hawk
import torch
import torch.nn as nn
from torch.utils.data import DataLoader

class TrainingMonitor:
    def __init__(self, model_name):
        self.tui = hawk.TUI(f"Training: {model_name}")
        self.setup_dashboard()
        
    def setup_dashboard(self):
        # Training metrics
        dash = self.tui.dashboard("Training Metrics")
        dash.add_chart("Loss", self.get_loss_history)
        dash.add_chart("Accuracy", self.get_accuracy_history)
        dash.add_metric("Learning Rate", lambda: self.current_lr)
        dash.add_metric("Epoch", lambda: self.current_epoch)
        
        # System metrics
        system = self.tui.dashboard("System Resources")
        system.add_metric("GPU Memory", self.get_gpu_memory)
        system.add_metric("GPU Utilization", self.get_gpu_utilization)
        system.add_chart("Temperature", self.get_temperature)
        
        # Controls
        controls = self.tui.command_panel("Training Controls")
        controls.add_command("Pause Training", self.pause_training)
        controls.add_command("Save Checkpoint", self.save_checkpoint)
        controls.add_command("Adjust Learning Rate", self.adjust_lr)

def train_model(model, dataloader, optimizer, criterion, epochs=100):
    monitor = TrainingMonitor("ResNet-50")
    
    for epoch in range(epochs):
        monitor.current_epoch = epoch
        epoch_loss = 0.0
        correct_predictions = 0
        total_samples = 0
        
        # Progress bar for epoch
        epoch_progress = hawk.progress(f"Epoch {epoch+1}/{epochs}", total=len(dataloader))
        
        for batch_idx, (data, targets) in enumerate(dataloader):
            # Forward pass
            outputs = model(data)
            loss = criterion(outputs, targets)
            
            # Backward pass
            optimizer.zero_grad()
            loss.backward()
            optimizer.step()
            
            # Update metrics
            epoch_loss += loss.item()
            _, predicted = torch.max(outputs.data, 1)
            total_samples += targets.size(0)
            correct_predictions += (predicted == targets).sum().item()
            
            # Log batch progress
            if batch_idx % 10 == 0:
                current_loss = loss.item()
                hawk.metric("batch_loss", current_loss)
                hawk.log(f"Batch {batch_idx}: Loss = {current_loss:.4f}")
            
            epoch_progress.update(batch_idx + 1)
        
        # Calculate epoch metrics
        avg_loss = epoch_loss / len(dataloader)
        accuracy = 100.0 * correct_predictions / total_samples
        
        hawk.metric("epoch_loss", avg_loss)
        hawk.metric("epoch_accuracy", accuracy)
        
        hawk.log(f"Epoch {epoch+1}: Loss = {avg_loss:.4f}, Accuracy = {accuracy:.2f}%")
        
        # Save checkpoint every 10 epochs
        if (epoch + 1) % 10 == 0:
            checkpoint_path = f"checkpoint_epoch_{epoch+1}.pth"
            torch.save(model.state_dict(), checkpoint_path)
            hawk.log(f"💾 Saved checkpoint: {checkpoint_path}", level="SUCCESS")
```

## 6. CI/CD Pipeline Monitor

### GitHub Actions Equivalent
```python
import hawk
import subprocess
import json

class CIPipeline:
    def __init__(self):
        self.tui = hawk.TUI("CI/CD Pipeline")
        self.setup_dashboard()
    
    def setup_dashboard(self):
        dash = self.tui.dashboard("Pipeline Status")
        dash.add_table("Build Steps", self.get_build_steps)
        dash.add_metric("Build Time", lambda: self.build_duration)
        dash.add_chart("Test Results", self.get_test_results)
        
    def run_pipeline(self):
        """Run complete CI/CD pipeline"""
        steps = [
            ("Checkout Code", self.checkout_code),
            ("Install Dependencies", self.install_deps),
            ("Run Linting", self.run_linting),
            ("Run Tests", self.run_tests),
            ("Build Application", self.build_app),
            ("Deploy to Staging", self.deploy_staging),
            ("Run Integration Tests", self.run_integration_tests),
            ("Deploy to Production", self.deploy_production)
        ]
        
        pipeline_progress = hawk.progress("Pipeline Progress", total=len(steps))
        
        for i, (step_name, step_func) in enumerate(steps):
            hawk.log(f"🔄 Running: {step_name}")
            
            try:
                with hawk.context(step_name):
                    result = step_func()
                    
                if result:
                    hawk.log(f"✅ {step_name} completed", level="SUCCESS")
                else:
                    hawk.log(f"❌ {step_name} failed", level="ERROR")
                    return False
                    
            except Exception as e:
                hawk.log(f"💥 {step_name} crashed: {e}", level="ERROR")
                return False
            
            pipeline_progress.update(i + 1)
        
        hawk.log("🎉 Pipeline completed successfully!", level="SUCCESS")
        return True
    
    def run_tests(self):
        """Run test suite with detailed reporting"""
        hawk.log("Running test suite...")
        
        # Run pytest with JSON output
        result = subprocess.run([
            "python", "-m", "pytest", 
            "--json-report", "--json-report-file=test_results.json",
            "tests/"
        ], capture_output=True, text=True)
        
        # Parse results
        with open("test_results.json", "r") as f:
            test_data = json.load(f)
        
        # Report metrics
        hawk.metric("tests_passed", test_data["summary"]["passed"])
        hawk.metric("tests_failed", test_data["summary"]["failed"])
        hawk.metric("test_duration", test_data["summary"]["duration"])
        
        # Log individual test failures
        for test in test_data["tests"]:
            if test["outcome"] == "failed":
                hawk.log(f"FAILED: {test['nodeid']}", level="ERROR")
                hawk.log(f"  {test['call']['longrepr']}", level="ERROR")
        
        return result.returncode == 0
```

These examples show how Hawk TUI can transform ordinary command-line applications into rich, interactive experiences with minimal code changes. The key is providing immediate value while scaling to enterprise needs.