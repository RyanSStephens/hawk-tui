# Node.js Client

The Hawk TUI Node.js client provides seamless integration with Node.js and JavaScript applications.

## Installation

```bash
npm install hawk-tui
```

## Quick Start

```javascript
const hawk = require('hawk-tui');

// Zero-configuration setup
hawk.auto();

// Your existing code works unchanged
console.log('Server starting...');
console.error('This is an error message');
```

## API Reference

### Core Functions

#### `hawk.auto(options?)`
Automatically captures console output, logs, and metrics.

```javascript
const hawk = require('hawk-tui');
hawk.auto();

// Now all console.log, console.error, etc. are captured
console.log('This appears in the TUI');
console.error('So does this error');
```

**Options:**
```javascript
hawk.auto({
    appName: 'My App',
    captureConsole: true,
    captureProcess: true,
    autoMetrics: true
});
```

#### `hawk.log(message, options?)`
Send a log message to the TUI.

```javascript
hawk.log('Server started');
hawk.log('Database error', { level: 'ERROR' });
hawk.log('User action', { 
    level: 'INFO', 
    context: { userId: 123, action: 'login' }
});
```

**Levels:** DEBUG, INFO, WARN, ERROR, SUCCESS

#### `hawk.metric(name, value, options?)`
Send a metric to the TUI.

```javascript
hawk.metric('requests_total', 1);
hawk.metric('response_time', 125.5, { type: 'histogram' });
hawk.metric('active_users', 42, { type: 'gauge', labels: { region: 'us-east' } });
```

**Types:** counter, gauge, histogram

### Progress Tracking

#### `hawk.progress(name, options?)`
Create a progress tracker for long-running operations.

```javascript
// With known total
const progress = hawk.progress('Processing files', { total: 1000 });
for (let i = 0; i < 1000; i++) {
    await processFile(i);
    progress.update(i + 1);
}

// With unknown total
const progress = hawk.progress('Loading data');
for (const item of dataStream) {
    await processItem(item);
    progress.increment();
}

// With callback completion
await hawk.progress('Migration', { total: tables.length }, async (p) => {
    for (let i = 0; i < tables.length; i++) {
        await migrateTable(tables[i]);
        p.update(i + 1);
    }
});
```

### Configuration Management

#### `hawk.config(name, defaultValue, options?)`
Register a configuration parameter.

```javascript
// Basic configuration
const port = hawk.config('port', 8080, { type: 'number' });
const debug = hawk.config('debug', false, { type: 'boolean' });

// With description
const logLevel = hawk.config('log_level', 'INFO', {
    type: 'string',
    description: 'Application log level',
    options: ['DEBUG', 'INFO', 'WARN', 'ERROR']
});

// Use in your application
if (debug) {
    console.log('Debug mode enabled');
}
```

### Context Management

#### `hawk.context(name, callback?)`
Group related operations under a named context.

```javascript
// With callback (automatic cleanup)
await hawk.context('Database Migration', async () => {
    await hawk.context('Users Table', async () => {
        await migrateUsers();
    });
    await hawk.context('Products Table', async () => {
        await migrateProducts();
    });
});

// Manual context control
const ctx = hawk.context('API Request');
try {
    await processRequest();
} finally {
    ctx.end();
}
```

### Status Updates

#### `hawk.status(name, value, options?)`
Update application status information.

```javascript
hawk.status('Database', 'Connected', { status: 'ok' });
hawk.status('Cache', 'Redis unavailable', { status: 'error' });
hawk.status('Queue', `${queue.size} jobs pending`, { status: 'warning' });
```

**Status levels:** ok, warning, error

### Events and Alerts

#### `hawk.event(name, data?, options?)`
Send significant events.

```javascript
hawk.event('User Registration', { userId: 123, email: 'user@example.com' });
hawk.event('System Restart', null, { level: 'warning' });
hawk.event('Critical Error', { error: err.message }, { level: 'error' });
```

### Timers and Benchmarks

#### `hawk.timer(name, callback?)`
Time operations and send metrics automatically.

```javascript
// With callback
const result = await hawk.timer('database_query', async () => {
    return await db.query('SELECT * FROM users');
});

// Manual timing
const timer = hawk.timer('api_request');
try {
    const result = await processRequest();
    return result;
} finally {
    timer.stop();
}
```

### Decorators and Middleware

#### `hawk.monitor(options?)`
Function decorator for automatic monitoring.

```javascript
// As decorator
const monitoredFunction = hawk.monitor(async function processRequest(req) {
    // Function timing and error handling is automatic
    return await handleRequest(req);
}, { name: 'process_request' });

// As wrapper
async function processData(data) {
    // Original function
}
const monitored = hawk.monitor(processData);
```

## Express.js Integration

```javascript
const express = require('express');
const hawk = require('hawk-tui');

const app = express();
hawk.auto();

// Request logging middleware
app.use((req, res, next) => {
    const start = Date.now();
    
    hawk.metric('http_requests_total', 1, {
        labels: { method: req.method, route: req.route?.path || req.path }
    });
    
    res.on('finish', () => {
        const duration = Date.now() - start;
        hawk.metric('http_request_duration_ms', duration, {
            type: 'histogram',
            labels: { 
                method: req.method, 
                status: res.statusCode,
                route: req.route?.path || req.path
            }
        });
    });
    
    next();
});

// Context middleware
app.use((req, res, next) => {
    req.hawkContext = hawk.context(`${req.method} ${req.path}`);
    res.on('finish', () => req.hawkContext.end());
    next();
});

// Route with monitoring
app.get('/api/users', async (req, res) => {
    try {
        const users = await hawk.timer('database_query', async () => {
            return await User.findAll();
        });
        
        hawk.log(`Retrieved ${users.length} users`);
        res.json({ users });
    } catch (error) {
        hawk.log(`Error fetching users: ${error.message}`, { level: 'ERROR' });
        res.status(500).json({ error: 'Internal server error' });
    }
});

app.listen(3000, () => {
    hawk.log('Server started on port 3000', { level: 'INFO' });
    hawk.status('Server', 'Running', { status: 'ok' });
});
```

## Next.js Integration

```javascript
// middleware.js
import { NextResponse } from 'next/server';
import hawk from 'hawk-tui';

hawk.auto();

export function middleware(request) {
    const start = Date.now();
    
    hawk.metric('nextjs_requests_total', 1, {
        labels: { path: request.nextUrl.pathname }
    });
    
    const response = NextResponse.next();
    
    response.headers.set('x-hawk-request-id', Date.now().toString());
    
    // Log completion (this is approximate in middleware)
    setTimeout(() => {
        const duration = Date.now() - start;
        hawk.metric('nextjs_request_duration_ms', duration, {
            type: 'histogram',
            labels: { path: request.nextUrl.pathname }
        });
    }, 0);
    
    return response;
}

// pages/api/users.js
import hawk from 'hawk-tui';

export default async function handler(req, res) {
    await hawk.context(`API ${req.method} /api/users`, async () => {
        try {
            const users = await hawk.timer('database_query', async () => {
                return await fetchUsers();
            });
            
            hawk.log(`Retrieved ${users.length} users`);
            res.status(200).json({ users });
        } catch (error) {
            hawk.log(`Error: ${error.message}`, { level: 'ERROR' });
            res.status(500).json({ error: 'Internal server error' });
        }
    });
}
```

## Socket.io Integration

```javascript
const io = require('socket.io')(server);
const hawk = require('hawk-tui');

hawk.auto();

io.on('connection', (socket) => {
    hawk.metric('websocket_connections', 1);
    hawk.log(`Client connected: ${socket.id}`);
    
    socket.on('message', async (data) => {
        await hawk.context(`WebSocket Message: ${data.type}`, async () => {
            hawk.metric('websocket_messages_total', 1, {
                labels: { type: data.type }
            });
            
            try {
                const result = await processMessage(data);
                socket.emit('response', result);
                hawk.log(`Processed message: ${data.type}`);
            } catch (error) {
                hawk.log(`Message error: ${error.message}`, { level: 'ERROR' });
                socket.emit('error', { message: error.message });
            }
        });
    });
    
    socket.on('disconnect', () => {
        hawk.metric('websocket_disconnections', 1);
        hawk.log(`Client disconnected: ${socket.id}`);
    });
});
```

## Worker Threads

```javascript
// main.js
const { Worker, isMainThread, parentPort } = require('worker_threads');
const hawk = require('hawk-tui');

if (isMainThread) {
    hawk.auto();
    
    const worker = new Worker(__filename);
    
    worker.on('message', (data) => {
        if (data.type === 'hawk') {
            // Forward worker messages to main TUI
            hawk.log(data.message, data.options);
        }
    });
    
    worker.postMessage({ task: 'process', data: largeDataset });
} else {
    // Worker thread
    parentPort.on('message', async ({ task, data }) => {
        if (task === 'process') {
            const progress = createWorkerProgress('Processing data', data.length);
            
            for (let i = 0; i < data.length; i++) {
                await processItem(data[i]);
                progress.update(i + 1);
            }
            
            parentPort.postMessage({
                type: 'hawk',
                message: 'Worker completed processing',
                options: { level: 'SUCCESS' }
            });
        }
    });
    
    function createWorkerProgress(name, total) {
        let current = 0;
        parentPort.postMessage({
            type: 'hawk',
            message: `Started: ${name}`,
            options: { level: 'INFO' }
        });
        
        return {
            update(value) {
                current = value;
                if (current % 100 === 0) {  // Update every 100 items
                    parentPort.postMessage({
                        type: 'hawk',
                        message: `Progress: ${current}/${total} (${Math.round(current/total*100)}%)`,
                        options: { level: 'INFO' }
                    });
                }
            }
        };
    }
}
```

## Database Integration

### MongoDB with Mongoose

```javascript
const mongoose = require('mongoose');
const hawk = require('hawk-tui');

hawk.auto();

// Monitor database connection
mongoose.connection.on('connected', () => {
    hawk.status('MongoDB', 'Connected', { status: 'ok' });
    hawk.log('MongoDB connected');
});

mongoose.connection.on('error', (err) => {
    hawk.status('MongoDB', 'Error', { status: 'error' });
    hawk.log(`MongoDB error: ${err.message}`, { level: 'ERROR' });
});

mongoose.connection.on('disconnected', () => {
    hawk.status('MongoDB', 'Disconnected', { status: 'warning' });
    hawk.log('MongoDB disconnected', { level: 'WARN' });
});

// Monitor queries
const originalExec = mongoose.Query.prototype.exec;
mongoose.Query.prototype.exec = function() {
    const start = Date.now();
    const model = this.model.modelName;
    const operation = this.op;
    
    return originalExec.call(this).then(
        result => {
            const duration = Date.now() - start;
            hawk.metric('mongodb_query_duration_ms', duration, {
                type: 'histogram',
                labels: { model, operation }
            });
            hawk.metric('mongodb_queries_total', 1, {
                labels: { model, operation, status: 'success' }
            });
            return result;
        },
        error => {
            const duration = Date.now() - start;
            hawk.metric('mongodb_query_duration_ms', duration, {
                type: 'histogram',
                labels: { model, operation }
            });
            hawk.metric('mongodb_queries_total', 1, {
                labels: { model, operation, status: 'error' }
            });
            hawk.log(`MongoDB query error: ${error.message}`, { level: 'ERROR' });
            throw error;
        }
    );
};
```

### PostgreSQL with pg

```javascript
const { Pool } = require('pg');
const hawk = require('hawk-tui');

hawk.auto();

const pool = new Pool({
    user: 'username',
    host: 'localhost',
    database: 'mydb',
    password: 'password',
    port: 5432,
});

// Monitor connection pool
pool.on('connect', () => {
    hawk.metric('postgres_connections_active', pool.totalCount);
    hawk.metric('postgres_connections_idle', pool.idleCount);
});

pool.on('error', (err) => {
    hawk.log(`PostgreSQL pool error: ${err.message}`, { level: 'ERROR' });
    hawk.status('PostgreSQL', 'Pool Error', { status: 'error' });
});

// Wrap query method
const originalQuery = pool.query.bind(pool);
pool.query = async function(text, params) {
    const start = Date.now();
    
    try {
        const result = await originalQuery(text, params);
        const duration = Date.now() - start;
        
        hawk.metric('postgres_query_duration_ms', duration, { type: 'histogram' });
        hawk.metric('postgres_queries_total', 1, { labels: { status: 'success' } });
        
        return result;
    } catch (error) {
        const duration = Date.now() - start;
        
        hawk.metric('postgres_query_duration_ms', duration, { type: 'histogram' });
        hawk.metric('postgres_queries_total', 1, { labels: { status: 'error' } });
        hawk.log(`PostgreSQL query error: ${error.message}`, { level: 'ERROR' });
        
        throw error;
    }
};
```

## Configuration

### Environment Variables

```javascript
const hawk = require('hawk-tui');

// Configure via environment
hawk.configure({
    appName: process.env.HAWK_APP_NAME || 'My Node.js App',
    logLevel: process.env.HAWK_LOG_LEVEL || 'INFO',
    metricsEnabled: process.env.HAWK_METRICS_ENABLED === 'true'
});

hawk.auto();
```

### Configuration File

```javascript
// config.json
{
    "appName": "My Node.js Application",
    "logLevel": "INFO",
    "autoDetect": {
        "console": true,
        "metrics": true,
        "process": true
    },
    "dashboard": {
        "widgets": ["logs", "metrics", "status"],
        "refreshRate": "1000ms"
    }
}
```

```javascript
const hawk = require('hawk-tui');
const config = require('./config.json');

hawk.configure(config);
hawk.auto();
```

## Error Handling

The Node.js client handles errors gracefully:

```javascript
const hawk = require('hawk-tui');

// This is safe even if Hawk TUI is not available
hawk.auto();

process.on('uncaughtException', (error) => {
    hawk.log(`Uncaught exception: ${error.message}`, { level: 'ERROR' });
    // Your error handling continues
});

process.on('unhandledRejection', (reason, promise) => {
    hawk.log(`Unhandled rejection: ${reason}`, { level: 'ERROR' });
    // Your error handling continues
});
```

## Performance Tips

1. **Batch operations** for high-frequency events:
```javascript
const batch = hawk.createBatch();
batch.metric('requests', 1);
batch.metric('errors', 0);
batch.log('Request processed');
await batch.send();
```

2. **Use sampling** for high-volume metrics:
```javascript
// Sample 10% of requests
if (Math.random() < 0.1) {
    hawk.metric('detailed_request_info', requestData);
}
```

3. **Avoid blocking operations**:
```javascript
// Good: async metrics
setImmediate(() => {
    hawk.metric('background_metric', calculateMetric());
});

// Bad: blocking main thread
hawk.metric('expensive_metric', expensiveCalculation()); // Don't do this
```

## TypeScript Support

```typescript
import * as hawk from 'hawk-tui';

interface UserAction {
    userId: number;
    action: string;
    timestamp: Date;
}

hawk.auto();

async function logUserAction(action: UserAction): Promise<void> {
    hawk.log('User action recorded', {
        level: 'INFO' as const,
        context: {
            userId: action.userId,
            action: action.action
        }
    });
    
    hawk.metric('user_actions_total', 1, {
        labels: { action: action.action }
    });
}

// Type-safe configuration
const config: hawk.Config = {
    appName: 'TypeScript App',
    autoDetect: {
        console: true,
        metrics: true
    }
};

hawk.configure(config);
```

## Debugging

Enable debug logging:

```javascript
const hawk = require('hawk-tui');

// Enable debug mode
hawk.setDebug(true);
hawk.auto();

// Or via environment
// DEBUG=hawk* node app.js
```

## Examples

See the [examples/nodejs/](../../examples/nodejs/) directory for complete working examples:

- [Basic Usage](../../examples/nodejs/demo.js)
- [Express Server](../../examples/nodejs/express_server.js)
- [WebSocket Server](../../examples/nodejs/websocket_server.js)
- [Database Integration](../../examples/nodejs/database_example.js)