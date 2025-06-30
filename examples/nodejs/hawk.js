#!/usr/bin/env node
/**
 * Hawk TUI Node.js Client Library - Universal TUI Framework
 * 
 * This library provides layered complexity similar to the Python client:
 * - Layer 0: hawk.auto() - Magic mode with zero configuration
 * - Layer 1: Simple functions like hawk.log(), hawk.metric()
 * - Layer 2: Structured integration with classes and middleware
 * - Layer 3: Advanced features and enterprise capabilities
 * 
 * Usage Examples:
 * 
 * Layer 0 - Magic mode:
 *     const hawk = require('hawk-tui-client');
 *     hawk.auto(); // Auto-detects and monitors everything
 * 
 * Layer 1 - Simple functions:
 *     hawk.log('Server started');
 *     hawk.metric('requests_per_second', 145);
 *     hawk.config('port', { default: 8080 });
 * 
 * Layer 2 - Middleware and classes:
 *     app.use(hawk.middleware());
 *     const monitor = new hawk.Monitor('my-app');
 * 
 * Requirements:
 * - Node.js 14+
 * - No external dependencies
 * - Thread-safe
 * - Works with or without Hawk TUI running
 */

const fs = require('fs');
const path = require('path');
const { performance } = require('perf_hooks');
const { EventEmitter } = require('events');

// Global configuration
let globalConfig = {
    enabled: true,
    appName: 'nodejs-app',
    debug: false,
    batchSize: 100,
    batchTimeout: 1000,
    maxMessageRate: 1000, // messages per second
    gracefulFallback: true
};

// Message batching and rate limiting
let messageBatch = [];
let batchTimer = null;
let messageCount = 0;
let lastReset = Date.now();

// Auto-detection state
let autoEnabled = false;
let originalConsole = {};

/**
 * Core message sending functionality
 */
function sendMessage(method, params) {
    if (!globalConfig.enabled) return;

    // Rate limiting
    const now = Date.now();
    if (now - lastReset >= 1000) {
        messageCount = 0;
        lastReset = now;
    }
    
    if (messageCount >= globalConfig.maxMessageRate) {
        if (globalConfig.debug) {
            console.warn('[Hawk] Rate limit exceeded, dropping message');
        }
        return;
    }
    
    messageCount++;

    const message = {
        jsonrpc: '2.0',
        method: method,
        params: {
            timestamp: new Date().toISOString(),
            app_name: globalConfig.appName,
            ...params
        }
    };

    if (globalConfig.debug) {
        console.error('[Hawk Debug]', JSON.stringify(message, null, 2));
    }

    // Add to batch
    messageBatch.push(message);

    // Send batch if it's full or start timer if it's the first message
    if (messageBatch.length >= globalConfig.batchSize) {
        flushBatch();
    } else if (messageBatch.length === 1 && !batchTimer) {
        batchTimer = setTimeout(flushBatch, globalConfig.batchTimeout);
    }
}

function flushBatch() {
    if (messageBatch.length === 0) return;

    try {
        for (const message of messageBatch) {
            process.stdout.write(JSON.stringify(message) + '\n');
        }
    } catch (error) {
        if (globalConfig.gracefulFallback && globalConfig.debug) {
            console.error('[Hawk] Failed to send messages:', error.message);
        }
    }

    messageBatch = [];
    if (batchTimer) {
        clearTimeout(batchTimer);
        batchTimer = null;
    }
}

/**
 * Layer 0: Magic Mode - Auto-detection
 */
function auto(options = {}) {
    if (autoEnabled) return;
    autoEnabled = true;

    // Update global config
    Object.assign(globalConfig, options);

    // Hook into console methods
    const logLevels = {
        log: 'INFO',
        info: 'INFO',
        warn: 'WARN',
        error: 'ERROR',
        debug: 'DEBUG'
    };

    Object.keys(logLevels).forEach(method => {
        originalConsole[method] = console[method];
        console[method] = function(...args) {
            // Send to Hawk TUI
            log(args.join(' '), { level: logLevels[method] });
            
            // Call original console method
            originalConsole[method].apply(console, args);
        };
    });

    // Hook into process events
    process.on('uncaughtException', (error) => {
        log(`Uncaught Exception: ${error.message}`, { 
            level: 'ERROR', 
            stack: error.stack 
        });
        originalConsole.error('Uncaught Exception:', error);
    });

    process.on('unhandledRejection', (reason, promise) => {
        log(`Unhandled Rejection: ${reason}`, { 
            level: 'ERROR', 
            promise_info: promise.toString() 
        });
        originalConsole.error('Unhandled Rejection:', reason);
    });

    // Auto-detect common patterns
    if (global.process && process.env) {
        // Environment detection
        const nodeEnv = process.env.NODE_ENV || 'development';
        config('node_env', { value: nodeEnv, description: 'Node.js environment' });
        
        const port = process.env.PORT;
        if (port) {
            config('port', { value: parseInt(port), description: 'Server port' });
        }
    }

    log('Hawk TUI auto-detection enabled', { level: 'SUCCESS' });
}

/**
 * Layer 1: Simple Functions
 */
function log(message, options = {}) {
    const params = {
        message: String(message),
        level: options.level || 'INFO',
        context: options.context || null,
        metadata: options.metadata || null,
        ...(options.stack && { stack: options.stack })
    };

    sendMessage('hawk.log', params);
}

function metric(name, value, options = {}) {
    const params = {
        name: String(name),
        value: Number(value),
        type: options.type || 'gauge',
        unit: options.unit || '',
        tags: options.tags || {},
        format: options.format || null
    };

    sendMessage('hawk.metric', params);
}

function config(key, options = {}) {
    const params = {
        key: String(key),
        value: options.value !== undefined ? options.value : options.default,
        type: options.type || 'auto',
        description: options.description || '',
        required: options.required || false,
        validation: options.validation || null
    };

    sendMessage('hawk.config', params);
    return params.value;
}

function progress(name, current, total, options = {}) {
    const params = {
        name: String(name),
        current: Number(current),
        total: Number(total),
        percentage: total > 0 ? Math.round((current / total) * 100) : 0,
        message: options.message || '',
        status: options.status || 'in_progress'
    };

    sendMessage('hawk.progress', params);
}

function event(name, data = {}, options = {}) {
    const params = {
        name: String(name),
        data: data,
        level: options.level || 'INFO',
        category: options.category || 'general',
        action: options.action || null
    };

    sendMessage('hawk.event', params);
}

/**
 * Layer 2: Classes and Middleware
 */
class Monitor {
    constructor(appName, options = {}) {
        this.appName = appName;
        this.options = options;
        this.startTime = Date.now();
        this.counters = new Map();
        this.gauges = new Map();
        this.timers = new Map();
    }

    counter(name, value = 1, tags = {}) {
        const current = this.counters.get(name) || 0;
        const newValue = current + value;
        this.counters.set(name, newValue);
        
        metric(name, newValue, { type: 'counter', tags });
        return newValue;
    }

    gauge(name, value, tags = {}) {
        this.gauges.set(name, value);
        metric(name, value, { type: 'gauge', tags });
        return value;
    }

    timer(name) {
        const startTime = performance.now();
        this.timers.set(name, startTime);

        return {
            stop: () => {
                const endTime = performance.now();
                const duration = endTime - startTime;
                metric(name, duration, { type: 'histogram', unit: 'ms' });
                this.timers.delete(name);
                return duration;
            }
        };
    }

    uptime() {
        const uptime = Date.now() - this.startTime;
        this.gauge('uptime_seconds', Math.floor(uptime / 1000));
        return uptime;
    }
}

class Dashboard {
    constructor(name, options = {}) {
        this.name = name;
        this.options = options;
        this.widgets = [];
    }

    addWidget(type, name, data = {}) {
        const widget = {
            type,
            name,
            data,
            timestamp: new Date().toISOString()
        };

        this.widgets.push(widget);

        sendMessage('hawk.dashboard', {
            dashboard_name: this.name,
            widget: widget
        });

        return widget;
    }

    addMetric(name, getValue, options = {}) {
        return this.addWidget('metric', name, {
            getValue: getValue.toString(),
            refresh_rate: options.refreshRate || 1000,
            format: options.format || null
        });
    }

    addChart(name, getDataPoints, options = {}) {
        return this.addWidget('chart', name, {
            getDataPoints: getDataPoints.toString(),
            chart_type: options.type || 'line',
            max_points: options.maxPoints || 100
        });
    }

    addTable(name, getRows, options = {}) {
        return this.addWidget('table', name, {
            getRows: getRows.toString(),
            columns: options.columns || [],
            max_rows: options.maxRows || 20
        });
    }
}

/**
 * Express.js Middleware
 */
function middleware(options = {}) {
    const monitor = new Monitor('express-app', options);

    return (req, res, next) => {
        const startTime = performance.now();
        const requestId = Math.random().toString(36).substr(2, 9);

        // Log request
        log(`${req.method} ${req.path}`, {
            level: 'INFO',
            metadata: {
                request_id: requestId,
                method: req.method,
                path: req.path,
                user_agent: req.get('User-Agent'),
                ip: req.ip
            }
        });

        // Track metrics
        monitor.counter('http_requests_total', 1, {
            method: req.method,
            route: req.route?.path || req.path
        });

        // Hook into response
        const originalSend = res.send;
        res.send = function(data) {
            const endTime = performance.now();
            const duration = endTime - startTime;

            // Log response
            log(`${req.method} ${req.path} - ${res.statusCode} (${duration.toFixed(2)}ms)`, {
                level: res.statusCode >= 400 ? 'ERROR' : 'INFO',
                metadata: {
                    request_id: requestId,
                    status_code: res.statusCode,
                    duration_ms: duration
                }
            });

            // Track response metrics
            monitor.counter('http_responses_total', 1, {
                method: req.method,
                status_code: res.statusCode.toString()
            });

            monitor.gauge('http_response_time_ms', duration, {
                method: req.method,
                route: req.route?.path || req.path
            });

            originalSend.call(this, data);
        };

        next();
    };
}

/**
 * Decorators and Higher-Order Functions
 */
function monitor(fn) {
    // Handle decorator usage @monitor
    if (arguments.length === 3) {
        const target = arguments[0];
        const propertyKey = arguments[1];
        const descriptor = arguments[2];
        
        if (descriptor && typeof descriptor.value === 'function') {
            const originalMethod = descriptor.value;
            
            descriptor.value = function(...args) {
                const timer = performance.now();
                const functionName = `${target.constructor.name}.${propertyKey}`;
                
                log(`Executing ${functionName}`, { level: 'DEBUG' });
                
                try {
                    const result = originalMethod.apply(this, args);
                    
                    // Handle promises
                    if (result && typeof result.then === 'function') {
                        return result
                            .then(value => {
                                const duration = performance.now() - timer;
                                metric(`${functionName}_duration`, duration, { unit: 'ms' });
                                log(`${functionName} completed (${duration.toFixed(2)}ms)`, { level: 'DEBUG' });
                                return value;
                            })
                            .catch(error => {
                                const duration = performance.now() - timer;
                                metric(`${functionName}_errors`, 1, { type: 'counter' });
                                log(`${functionName} failed: ${error.message}`, { 
                                    level: 'ERROR',
                                    metadata: { duration_ms: duration }
                                });
                                throw error;
                            });
                    }
                    
                    const duration = performance.now() - timer;
                    metric(`${functionName}_duration`, duration, { unit: 'ms' });
                    log(`${functionName} completed (${duration.toFixed(2)}ms)`, { level: 'DEBUG' });
                    return result;
                    
                } catch (error) {
                    const duration = performance.now() - timer;
                    metric(`${functionName}_errors`, 1, { type: 'counter' });
                    log(`${functionName} failed: ${error.message}`, { 
                        level: 'ERROR',
                        metadata: { duration_ms: duration }
                    });
                    throw error;
                }
            };
            
            return descriptor;
        }
    }
    
    // Function decoration: monitor(function)
    if (typeof fn === 'function') {
        return function(...args) {
            const timer = performance.now();
            const functionName = fn.name || 'anonymous';
            
            log(`Executing ${functionName}`, { level: 'DEBUG' });
            
            try {
                const result = fn.apply(this, args);
                
                if (result && typeof result.then === 'function') {
                    return result
                        .then(value => {
                            const duration = performance.now() - timer;
                            metric(`${functionName}_duration`, duration, { unit: 'ms' });
                            log(`${functionName} completed (${duration.toFixed(2)}ms)`, { level: 'DEBUG' });
                            return value;
                        })
                        .catch(error => {
                            const duration = performance.now() - timer;
                            metric(`${functionName}_errors`, 1, { type: 'counter' });
                            log(`${functionName} failed: ${error.message}`, { level: 'ERROR' });
                            throw error;
                        });
                }
                
                const duration = performance.now() - timer;
                metric(`${functionName}_duration`, duration, { unit: 'ms' });
                log(`${functionName} completed (${duration.toFixed(2)}ms)`, { level: 'DEBUG' });
                return result;
                
            } catch (error) {
                const duration = performance.now() - timer;
                metric(`${functionName}_errors`, 1, { type: 'counter' });
                log(`${functionName} failed: ${error.message}`, { level: 'ERROR' });
                throw error;
            }
        };
    }
    
    throw new Error('monitor() requires a function argument');
}

function timed(name) {
    return function(target, propertyKey, descriptor) {
        if (descriptor && typeof descriptor.value === 'function') {
            const originalMethod = descriptor.value;
            
            descriptor.value = function(...args) {
                const timer = performance.now();
                
                try {
                    const result = originalMethod.apply(this, args);
                    
                    if (result && typeof result.then === 'function') {
                        return result.finally(() => {
                            const duration = performance.now() - timer;
                            metric(name, duration, { unit: 'ms', type: 'histogram' });
                        });
                    }
                    
                    const duration = performance.now() - timer;
                    metric(name, duration, { unit: 'ms', type: 'histogram' });
                    return result;
                    
                } catch (error) {
                    const duration = performance.now() - timer;
                    metric(name, duration, { unit: 'ms', type: 'histogram' });
                    throw error;
                }
            };
            
            return descriptor;
        }
    };
}

/**
 * Context Management
 */
function context(name, fn) {
    if (typeof fn === 'function') {
        // Function form: context('name', () => { ... })
        log(`Entering context: ${name}`, { level: 'DEBUG' });
        const timer = performance.now();
        
        try {
            const result = fn();
            
            if (result && typeof result.then === 'function') {
                return result.finally(() => {
                    const duration = performance.now() - timer;
                    log(`Exiting context: ${name} (${duration.toFixed(2)}ms)`, { 
                        level: 'DEBUG',
                        metadata: { duration_ms: duration }
                    });
                });
            }
            
            const duration = performance.now() - timer;
            log(`Exiting context: ${name} (${duration.toFixed(2)}ms)`, { 
                level: 'DEBUG',
                metadata: { duration_ms: duration }
            });
            return result;
            
        } catch (error) {
            const duration = performance.now() - timer;
            log(`Context ${name} failed: ${error.message}`, { 
                level: 'ERROR',
                metadata: { duration_ms: duration }
            });
            throw error;
        }
    }
    
    // Object form: returns context manager
    return {
        name,
        start() {
            this.startTime = performance.now();
            log(`Entering context: ${name}`, { level: 'DEBUG' });
        },
        end() {
            if (this.startTime) {
                const duration = performance.now() - this.startTime;
                log(`Exiting context: ${name} (${duration.toFixed(2)}ms)`, { 
                    level: 'DEBUG',
                    metadata: { duration_ms: duration }
                });
            }
        }
    };
}

/**
 * Utility Functions
 */
function configure(options = {}) {
    Object.assign(globalConfig, options);
    log('Hawk TUI configuration updated', { 
        level: 'DEBUG',
        metadata: { config: globalConfig }
    });
}

function isEnabled() {
    return globalConfig.enabled;
}

function disable() {
    globalConfig.enabled = false;
    log('Hawk TUI disabled', { level: 'DEBUG' });
}

function enable() {
    globalConfig.enabled = true;
    log('Hawk TUI enabled', { level: 'DEBUG' });
}

function flush() {
    flushBatch();
}

// Graceful shutdown
process.on('exit', () => {
    flush();
});

process.on('SIGINT', () => {
    flush();
    process.exit(0);
});

process.on('SIGTERM', () => {
    flush();
    process.exit(0);
});

// Export API
module.exports = {
    // Layer 0: Magic mode
    auto,
    
    // Layer 1: Simple functions
    log,
    metric,
    config,
    progress,
    event,
    
    // Layer 2: Classes and middleware
    Monitor,
    Dashboard,
    middleware,
    
    // Decorators
    monitor,
    timed,
    context,
    
    // Utilities
    configure,
    isEnabled,
    disable,
    enable,
    flush,
    
    // Version info
    version: '1.0.0'
};