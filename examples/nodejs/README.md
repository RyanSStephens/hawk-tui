# Hawk TUI Node.js Client

Universal TUI framework client library for Node.js applications. Transform any Node.js application into a beautiful, interactive terminal interface with minimal code changes.

## Quick Start

### Installation

```bash
npm install hawk-tui-client
```

### Basic Usage

```javascript
const hawk = require('hawk-tui-client');

// Magic mode - auto-detects everything
hawk.auto();

// Your existing application code works unchanged
console.log('Server starting...');  // Appears in TUI
console.error('Database error');    // Automatically captured

// Add custom metrics and logging
hawk.log('Custom message', { level: 'SUCCESS' });
hawk.metric('requests_per_second', 145);
```

### Run with TUI

```bash
node your-app.js | hawk
```

## Features

### Layer 0: Magic Mode (Zero Configuration)

```javascript
const hawk = require('hawk-tui-client');
hawk.auto();

// All console.* calls are automatically captured
console.log('This appears in the TUI!');
console.error('Errors are highlighted!');
console.warn('Warnings are color-coded!');
```

### Layer 1: Simple Functions

```javascript
// Logging with levels
hawk.log('Server started', { level: 'SUCCESS' });
hawk.log('Processing request', { level: 'INFO' });
hawk.log('High memory usage', { level: 'WARN' });
hawk.log('Database error', { level: 'ERROR' });

// Metrics
hawk.metric('requests_per_second', 145);
hawk.metric('cpu_usage', 67.5, { unit: '%' });
hawk.metric('memory_usage', 1024, { unit: 'MB', tags: { service: 'api' } });

// Configuration
const port = hawk.config('port', { default: 3000, type: 'number' });
const dbUrl = hawk.config('database_url', { 
    default: 'mongodb://localhost',
    description: 'Database connection string'
});

// Progress tracking
hawk.progress('file_upload', 75, 100, { message: 'Uploading files...' });

// Events
hawk.event('user_login', { user_id: 123, ip: '192.168.1.1' });
```

### Layer 2: Express.js Middleware

```javascript
const express = require('express');
const hawk = require('hawk-tui-client');

const app = express();

// Add Hawk monitoring middleware
app.use(hawk.middleware());

app.get('/api/users', (req, res) => {
    // Requests are automatically logged and timed
    res.json({ users: [] });
});

app.listen(3000);
```

### Layer 2: Monitor Class

```javascript
const monitor = new hawk.Monitor('my-service');

// Counters
monitor.counter('http_requests', 1, { method: 'GET' });

// Gauges
monitor.gauge('active_connections', 42);

// Timers
const timer = monitor.timer('database_query');
await performDatabaseQuery();
timer.stop(); // Automatically records duration

// Uptime tracking
monitor.uptime();
```

### Layer 2: Dashboard Creation

```javascript
const dashboard = new hawk.Dashboard('System Overview');

// Add metrics
dashboard.addMetric('CPU Usage', () => getCpuUsage(), {
    format: '{:.1f}%',
    refreshRate: 1000
});

// Add charts
dashboard.addChart('Response Times', () => getResponseTimes(), {
    type: 'line',
    maxPoints: 50
});

// Add tables
dashboard.addTable('Active Users', () => getActiveUsers(), {
    columns: ['ID', 'Username', 'Status'],
    maxRows: 20
});
```

### Function Monitoring

```javascript
// Decorator for automatic monitoring
const monitoredFunction = hawk.monitor(async function processPayment(amount) {
    // Function execution is automatically timed and logged
    await processPaymentLogic(amount);
});

// Timing specific operations
const timedFunction = hawk.timed('api_request')(async function callAPI() {
    return await fetch('/api/data');
});
```

### Context Management

```javascript
// Function form
await hawk.context('Database Migration', async () => {
    hawk.log('Starting migration...');
    await migrateTables();
    hawk.log('Migration completed');
});

// Object form for manual control
const ctx = hawk.context('File Processing');
ctx.start();
// ... do work ...
ctx.end();
```

## API Reference

### Configuration

```javascript
hawk.configure({
    enabled: true,           // Enable/disable Hawk
    appName: 'my-app',      // Application identifier
    debug: false,           // Debug mode
    batchSize: 100,         // Messages per batch
    batchTimeout: 1000,     // Batch timeout (ms)
    maxMessageRate: 1000,   // Max messages per second
    gracefulFallback: true  // Continue if TUI unavailable
});
```

### Log Levels

- `DEBUG` - Detailed information for debugging
- `INFO` - General information (default)
- `WARN` - Warning messages
- `ERROR` - Error conditions
- `SUCCESS` - Success notifications

### Metric Types

- `gauge` - Current value (default)
- `counter` - Incrementing counter
- `histogram` - Distribution of values

### Utility Functions

```javascript
hawk.isEnabled()    // Check if Hawk is enabled
hawk.disable()      // Disable Hawk
hawk.enable()       // Enable Hawk
hawk.flush()        // Force send pending messages
```

## TypeScript Support

Full TypeScript definitions are included:

```typescript
import * as hawk from 'hawk-tui-client';

hawk.log('Type-safe logging!', { level: 'INFO' });
hawk.metric('typed_metric', 42, { type: 'gauge' });

const monitor = new hawk.Monitor('typed-service');
monitor.counter('requests', 1);
```

## Express.js Integration

```javascript
const express = require('express');
const hawk = require('hawk-tui-client');

const app = express();

// Enable auto-monitoring
app.use(hawk.middleware());

app.get('/api/health', (req, res) => {
    // Automatically logged: "GET /api/health - 200 (15ms)"
    res.json({ status: 'healthy' });
});

app.listen(3000, () => {
    hawk.log('Server started on port 3000', { level: 'SUCCESS' });
});
```

## Real-World Examples

### Web API Server

```javascript
const express = require('express');
const hawk = require('hawk-tui-client');

hawk.auto();
const app = express();
const monitor = new hawk.Monitor('api-server');

app.use(hawk.middleware());

app.get('/api/users/:id', async (req, res) => {
    const timer = monitor.timer('user_lookup');
    
    try {
        const user = await getUserById(req.params.id);
        monitor.counter('successful_requests', 1);
        res.json(user);
    } catch (error) {
        monitor.counter('failed_requests', 1);
        hawk.log(`User lookup failed: ${error.message}`, { level: 'ERROR' });
        res.status(500).json({ error: 'User not found' });
    } finally {
        timer.stop();
    }
});

app.listen(3000);
```

### Database Operations

```javascript
const hawk = require('hawk-tui-client');

class DatabaseService {
    @hawk.monitor
    async migrateDatabase() {
        const tables = ['users', 'orders', 'products'];
        
        for (let i = 0; i < tables.length; i++) {
            const table = tables[i];
            
            hawk.progress('migration', i, tables.length, {
                message: `Migrating ${table} table...`
            });
            
            await hawk.context(`Migrate ${table}`, async () => {
                await this.migrateTable(table);
            });
        }
        
        hawk.log('Database migration completed!', { level: 'SUCCESS' });
    }
}
```

### Background Job Processing

```javascript
const hawk = require('hawk-tui-client');

class JobProcessor {
    constructor() {
        this.monitor = new hawk.Monitor('job-processor');
        this.setupDashboard();
    }
    
    setupDashboard() {
        const dashboard = new hawk.Dashboard('Job Processing');
        
        dashboard.addMetric('Queue Size', () => this.getQueueSize());
        dashboard.addMetric('Jobs/min', () => this.getJobsPerMinute());
        dashboard.addTable('Recent Jobs', () => this.getRecentJobs());
    }
    
    @hawk.timed('job_processing')
    async processJob(job) {
        hawk.log(`Processing job: ${job.type}`, { 
            level: 'INFO',
            metadata: { job_id: job.id, priority: job.priority }
        });
        
        try {
            await this.executeJob(job);
            this.monitor.counter('jobs_completed', 1, { type: job.type });
            hawk.log(`Job completed: ${job.id}`, { level: 'SUCCESS' });
        } catch (error) {
            this.monitor.counter('jobs_failed', 1, { type: job.type });
            throw error;
        }
    }
}
```

## Testing

```bash
# Run the demo
npm run demo

# Run with Hawk TUI
npm run demo | hawk

# Run tests
npm test
```

## Performance

- **Zero overhead** when TUI is not running
- **Batched messaging** for high-throughput applications
- **Rate limiting** prevents overwhelming the TUI
- **Graceful fallback** if TUI becomes unavailable
- **Thread-safe** for multi-threaded applications

## Requirements

- Node.js 14 or higher
- No external dependencies
- Works with CommonJS and ES modules

## License

AGPL-3.0 for open source use. Commercial licenses available.

For commercial licensing inquiries: license@hawktui.dev

## Links

- **Documentation**: https://hawktui.dev/docs
- **GitHub**: https://github.com/hawk-tui/hawk
- **Issues**: https://github.com/hawk-tui/hawk/issues
- **Examples**: https://github.com/hawk-tui/hawk/tree/main/examples/nodejs