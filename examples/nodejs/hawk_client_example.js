#!/usr/bin/env node

/**
 * Hawk TUI Node.js Client Example
 * 
 * This example demonstrates how to send messages to the Hawk TUI using the JSON-RPC protocol.
 * In a real implementation, this would be part of a client library that abstracts the protocol details.
 */

class HawkClient {
    constructor(appName = 'nodejs-app') {
        this.appName = appName;
        this.sessionId = `sess_${Date.now()}`;
        this.sequence = 0;
    }

    _sendMessage(method, params, messageId = null) {
        this.sequence++;
        
        const message = {
            jsonrpc: '2.0',
            method: method,
            params: params
        };
        
        if (messageId !== null) {
            message.id = messageId;
        }
        
        // Add Hawk metadata
        message.hawk_meta = {
            app_name: this.appName,
            session_id: this.sessionId,
            sequence: this.sequence
        };
        
        // Send to stdout (this would go to the TUI process)
        console.log(JSON.stringify(message));
    }

    _sendBatch(messages) {
        console.log(JSON.stringify(messages));
    }

    // Logging methods
    log(message, level = 'INFO', options = {}) {
        const params = {
            message: message,
            level: level,
            timestamp: new Date().toISOString()
        };
        
        if (options.context) params.context = options.context;
        if (options.tags) params.tags = options.tags;
        if (options.component) params.component = options.component;
        
        this._sendMessage('hawk.log', params);
    }

    debug(message, options = {}) {
        this.log(message, 'DEBUG', options);
    }

    info(message, options = {}) {
        this.log(message, 'INFO', options);
    }

    warn(message, options = {}) {
        this.log(message, 'WARN', options);
    }

    error(message, options = {}) {
        this.log(message, 'ERROR', options);
    }

    success(message, options = {}) {
        this.log(message, 'SUCCESS', options);
    }

    // Metrics methods
    metric(name, value, options = {}) {
        const params = {
            name: name,
            value: value,
            type: options.type || 'gauge',
            timestamp: new Date().toISOString()
        };
        
        if (options.unit) params.unit = options.unit;
        if (options.tags) params.tags = options.tags;
        
        this._sendMessage('hawk.metric', params);
    }

    counter(name, value = 1, options = {}) {
        this.metric(name, value, { ...options, type: 'counter' });
    }

    gauge(name, value, options = {}) {
        this.metric(name, value, { ...options, type: 'gauge' });
    }

    histogram(name, value, options = {}) {
        this.metric(name, value, { ...options, type: 'histogram' });
    }

    // Configuration methods
    config(key, options = {}) {
        const params = {
            key: key,
            type: options.type || 'string'
        };
        
        if (options.value !== undefined) params.value = options.value;
        if (options.description) params.description = options.description;
        if (options.default !== undefined) params.default = options.default;
        if (options.min !== undefined) params.min = options.min;
        if (options.max !== undefined) params.max = options.max;
        if (options.options) params.options = options.options;
        if (options.restartRequired) params.restart_required = options.restartRequired;
        if (options.category) params.category = options.category;
        
        this._sendMessage('hawk.config', params);
    }

    // Progress tracking methods
    progress(id, label, current, total, options = {}) {
        const params = {
            id: id,
            label: label,
            current: current,
            total: total,
            status: options.status || 'in_progress'
        };
        
        if (options.unit) params.unit = options.unit;
        if (options.details) params.details = options.details;
        
        this._sendMessage('hawk.progress', params);
    }

    progressComplete(id, label, total, unit = '') {
        this.progress(id, label, total, total, { status: 'completed', unit });
    }

    progressError(id, label, current, total, errorMessage, unit = '') {
        this.progress(id, label, current, total, { 
            status: 'error', 
            unit, 
            details: errorMessage 
        });
    }

    // Dashboard methods
    dashboardWidget(widgetId, type, title, data, options = {}) {
        const params = {
            widget_id: widgetId,
            type: type,
            title: title,
            data: data
        };
        
        if (options.layout) params.layout = options.layout;
        if (options.config) params.config = options.config;
        
        this._sendMessage('hawk.dashboard', params);
    }

    statusGrid(widgetId, title, services, options = {}) {
        this.dashboardWidget(widgetId, 'status_grid', title, services, options);
    }

    metricChart(widgetId, title, chartData, options = {}) {
        this.dashboardWidget(widgetId, 'metric_chart', title, chartData, options);
    }

    table(widgetId, title, headers, rows, options = {}) {
        const tableData = { headers, rows };
        this.dashboardWidget(widgetId, 'table', title, tableData, options);
    }

    // Event methods
    event(type, title, options = {}) {
        const params = {
            type: type,
            title: title,
            severity: options.severity || 'info',
            timestamp: new Date().toISOString()
        };
        
        if (options.message) params.message = options.message;
        if (options.data) params.data = options.data;
        
        this._sendMessage('hawk.event', params);
    }
}

// Simulate an Express.js web application
function simulateExpressApp() {
    const hawk = new HawkClient('express-api');
    
    hawk.info('Starting Express.js API server', { component: 'server' });
    
    // Configuration
    hawk.config('port', { 
        value: 3000, 
        type: 'integer', 
        description: 'HTTP server port',
        min: 1,
        max: 65535,
        default: 3000
    });
    
    hawk.config('env', {
        value: 'development',
        type: 'enum',
        description: 'Environment mode',
        options: ['development', 'staging', 'production'],
        default: 'development'
    });
    
    hawk.config('cors.enabled', {
        value: true,
        type: 'boolean',
        description: 'Enable CORS middleware',
        default: true
    });
    
    // Initial metrics
    hawk.gauge('server.uptime', 0, { unit: 'seconds' });
    hawk.counter('requests.total', 0);
    hawk.gauge('requests.active', 0);
    
    // Service status dashboard
    hawk.statusGrid('services', 'Service Health', {
        'Express Server': { status: 'starting', response_time: '0ms' },
        'MongoDB': { status: 'healthy', response_time: '8ms' },
        'Redis': { status: 'healthy', response_time: '1ms' },
        'Elasticsearch': { status: 'healthy', response_time: '15ms' }
    });
    
    hawk.success('Express server started on port 3000', { component: 'server' });
    hawk.event('server_started', 'API Server Started', {
        message: 'Express.js API server is now accepting requests',
        data: { port: 3000, env: 'development' }
    });
    
    // Simulate API requests and metrics
    let uptime = 0;
    let totalRequests = 0;
    
    const interval = setInterval(() => {
        uptime += 1;
        
        // Simulate varying API load
        const activeRequests = Math.max(0, 10 + Math.floor(Math.random() * 20));
        const requestsPerSecond = Math.max(0, 50 + Math.floor(Math.random() * 100));
        const avgResponseTime = 50 + Math.floor(Math.random() * 100);
        
        totalRequests += requestsPerSecond;
        
        // Update metrics
        hawk.gauge('server.uptime', uptime, { unit: 'seconds' });
        hawk.gauge('requests.active', activeRequests);
        hawk.gauge('requests.per_second', requestsPerSecond, { unit: 'req/s' });
        hawk.gauge('response.avg_time', avgResponseTime, { unit: 'ms' });
        hawk.counter('requests.total', totalRequests);
        
        // Memory and CPU metrics
        const memUsage = 45 + Math.random() * 20;
        const cpuUsage = 30 + Math.random() * 40;
        
        hawk.gauge('system.memory_usage', memUsage, { unit: '%' });
        hawk.gauge('system.cpu_usage', cpuUsage, { unit: '%' });
        
        // Update service status
        hawk.statusGrid('services', 'Service Health', {
            'Express Server': { 
                status: 'healthy', 
                response_time: `${avgResponseTime}ms`,
                last_checked: new Date().toISOString()
            },
            'MongoDB': { 
                status: 'healthy', 
                response_time: '8ms',
                last_checked: new Date().toISOString()
            },
            'Redis': { 
                status: 'healthy', 
                response_time: '1ms',
                last_checked: new Date().toISOString()
            },
            'Elasticsearch': { 
                status: avgResponseTime > 120 ? 'degraded' : 'healthy', 
                response_time: '15ms',
                last_checked: new Date().toISOString()
            }
        });
        
        // Occasional log messages
        if (uptime % 10 === 0) {
            hawk.info(`Server running for ${uptime} seconds, processed ${totalRequests} requests`, {
                context: { 
                    uptime_seconds: uptime,
                    total_requests: totalRequests,
                    avg_rps: Math.floor(totalRequests / uptime)
                },
                component: 'stats'
            });
        }
        
        // Simulate occasional warnings
        if (avgResponseTime > 120) {
            hawk.warn(`High response time detected: ${avgResponseTime}ms`, {
                context: { response_time: avgResponseTime, threshold: 120 },
                component: 'performance'
            });
        }
        
        if (memUsage > 80) {
            hawk.warn(`High memory usage: ${memUsage.toFixed(1)}%`, {
                context: { memory_usage: memUsage, threshold: 80 },
                component: 'system'
            });
        }
        
        // API endpoint metrics with tags
        const endpoints = ['/api/users', '/api/posts', '/api/auth', '/api/admin'];
        const methods = ['GET', 'POST', 'PUT', 'DELETE'];
        
        for (const endpoint of endpoints) {
            for (const method of methods) {
                if (Math.random() < 0.3) { // 30% chance per endpoint/method
                    const requests = Math.floor(Math.random() * 20);
                    const responseTime = 20 + Math.random() * 200;
                    
                    hawk.counter('api.requests', requests, {
                        tags: { endpoint, method }
                    });
                    
                    hawk.histogram('api.response_time', responseTime, {
                        tags: { endpoint, method },
                        unit: 'ms'
                    });
                }
            }
        }
        
        // Stop after 30 seconds
        if (uptime >= 30) {
            clearInterval(interval);
            simulateShutdown(hawk);
        }
    }, 1000);
}

function simulateShutdown(hawk) {
    hawk.info('Initiating graceful shutdown', { component: 'server' });
    
    let shutdownStep = 0;
    const shutdownSteps = [
        'Stopping new request acceptance',
        'Closing database connections',
        'Flushing logs and metrics',
        'Cleaning up resources',
        'Server shutdown complete'
    ];
    
    const shutdownInterval = setInterval(() => {
        hawk.progress('shutdown', 'Graceful Shutdown', shutdownStep + 1, shutdownSteps.length, {
            details: shutdownSteps[shutdownStep]
        });
        
        hawk.info(shutdownSteps[shutdownStep], { component: 'shutdown' });
        
        shutdownStep++;
        
        if (shutdownStep >= shutdownSteps.length) {
            clearInterval(shutdownInterval);
            hawk.progressComplete('shutdown', 'Shutdown Complete', shutdownSteps.length);
            hawk.event('server_stopped', 'Server Shutdown', {
                message: 'Express.js server shutdown completed successfully',
                severity: 'success'
            });
            hawk.success('Server shutdown completed', { component: 'server' });
        }
    }, 500);
}

function demonstrateRealTimeData() {
    const hawk = new HawkClient('realtime-demo');
    
    hawk.info('Starting real-time data demonstration');
    
    // Create a real-time chart
    const chartData = {
        series: [
            {
                name: 'Websocket Connections',
                color: '#007acc',
                data: []
            },
            {
                name: 'Messages per Second',
                color: '#ff6b6b',
                data: []
            }
        ]
    };
    
    let dataPoints = 0;
    const maxDataPoints = 20;
    
    const updateChart = setInterval(() => {
        const now = Date.now() / 1000;
        const connections = 100 + Math.sin(dataPoints * 0.2) * 30 + Math.random() * 20;
        const messagesPerSec = 50 + Math.cos(dataPoints * 0.3) * 25 + Math.random() * 15;
        
        // Add new data points
        chartData.series[0].data.push({ x: now, y: Math.floor(connections) });
        chartData.series[1].data.push({ x: now, y: Math.floor(messagesPerSec) });
        
        // Keep only the last N data points
        if (chartData.series[0].data.length > maxDataPoints) {
            chartData.series[0].data.shift();
            chartData.series[1].data.shift();
        }
        
        hawk.metricChart('realtime_chart', 'Real-time WebSocket Metrics', chartData, {
            layout: { row: 0, col: 0, width: 12, height: 6 }
        });
        
        // Also send as individual metrics
        hawk.gauge('websocket.connections', Math.floor(connections));
        hawk.gauge('websocket.messages_per_second', Math.floor(messagesPerSec), { unit: 'msg/s' });
        
        dataPoints++;
        
        if (dataPoints >= 50) {
            clearInterval(updateChart);
            hawk.info('Real-time demonstration completed');
        }
    }, 200);
}

function demonstrateErrorScenarios() {
    const hawk = new HawkClient('error-demo');
    
    hawk.info('Demonstrating error scenarios and recovery');
    
    // Simulate database connection issues
    hawk.error('Database connection lost', {
        context: {
            error_code: 'ECONNREFUSED',
            host: 'mongodb://localhost:27017',
            retry_count: 0
        },
        component: 'database'
    });
    
    // Simulate retry logic
    let retryCount = 0;
    const maxRetries = 3;
    
    const retryInterval = setInterval(() => {
        retryCount++;
        
        hawk.progress('db_reconnect', 'Database Reconnection', retryCount, maxRetries, {
            details: `Attempt ${retryCount}/${maxRetries} - Connecting to MongoDB`
        });
        
        hawk.info(`Database reconnection attempt ${retryCount}`, {
            context: { attempt: retryCount, max_attempts: maxRetries },
            component: 'database'
        });
        
        if (retryCount === maxRetries) {
            clearInterval(retryInterval);
            hawk.progressComplete('db_reconnect', 'Reconnection Successful', maxRetries);
            hawk.success('Database connection restored', { component: 'database' });
            
            hawk.event('database_recovered', 'Database Connection Restored', {
                message: 'Successfully reconnected to MongoDB after connection failure',
                severity: 'success',
                data: { attempts: retryCount, duration: `${retryCount * 2}s` }
            });
        } else {
            hawk.warn(`Reconnection attempt ${retryCount} failed, retrying...`, {
                component: 'database'
            });
        }
    }, 2000);
}

// Main execution
if (require.main === module) {
    console.error('Hawk TUI Node.js Client Example');
    console.error('================================');
    console.error('Note: Messages are sent to stdout in JSON-RPC format');
    console.error('In a real scenario, these would be processed by the Hawk TUI');
    console.error('');
    
    try {
        // Start the Express simulation
        simulateExpressApp();
        
        // After 10 seconds, start real-time data demo
        setTimeout(() => {
            console.error('\n--- Real-time Data Demo ---');
            demonstrateRealTimeData();
        }, 10000);
        
        // After 20 seconds, start error scenarios demo
        setTimeout(() => {
            console.error('\n--- Error Scenarios Demo ---');
            demonstrateErrorScenarios();
        }, 20000);
        
    } catch (error) {
        console.error(`\nError running example: ${error.message}`);
        process.exit(1);
    }
}