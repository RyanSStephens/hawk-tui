#!/usr/bin/env node
/**
 * Hawk TUI Node.js Client Demo
 * 
 * This demo shows all the capabilities of the Hawk TUI Node.js client library.
 * Run with: node demo.js | hawk
 */

const hawk = require('./hawk');

// Configure Hawk
hawk.configure({
    appName: 'nodejs-demo',
    debug: false,
    batchTimeout: 500
});

console.log('Starting Hawk TUI Node.js Demo...');
console.log('Run this with: node demo.js | hawk');
console.log('');

// Demo different layers of the API
async function runDemo() {
    // Layer 0: Magic mode demonstration
    console.log('=== Layer 0: Magic Mode ===');
    hawk.auto({ debug: false });
    
    console.log('This is a regular console.log - auto-detected!');
    console.info('This is console.info - auto-detected!');
    console.warn('This is console.warn - auto-detected!');
    console.error('This is console.error - auto-detected!');
    
    await sleep(1000);

    // Layer 1: Simple functions
    console.log('\n=== Layer 1: Simple Functions ===');
    
    hawk.log('Demo application started', { level: 'SUCCESS' });
    hawk.log('Loading configuration...', { level: 'INFO' });
    hawk.log('Warning: High memory usage detected', { level: 'WARN' });
    
    // Metrics
    for (let i = 0; i < 10; i++) {
        hawk.metric('requests_per_second', 100 + Math.random() * 50);
        hawk.metric('cpu_usage_percent', 20 + Math.random() * 30);
        hawk.metric('memory_usage_mb', 500 + Math.random() * 200);
        hawk.metric('active_connections', Math.floor(10 + Math.random() * 20));
        await sleep(200);
    }

    // Configuration
    const port = hawk.config('server_port', { 
        default: 3000, 
        type: 'number',
        description: 'HTTP server port' 
    });
    
    const dbUrl = hawk.config('database_url', { 
        default: 'mongodb://localhost:27017/myapp',
        type: 'string',
        description: 'Database connection URL'
    });

    hawk.log(`Server configured on port ${port}`, { level: 'INFO' });

    // Progress tracking
    console.log('\n=== Progress Tracking ===');
    const totalSteps = 20;
    for (let i = 0; i <= totalSteps; i++) {
        hawk.progress('file_processing', i, totalSteps, {
            message: `Processing file ${i}/${totalSteps}`,
            status: i === totalSteps ? 'completed' : 'in_progress'
        });
        await sleep(100);
    }

    // Events
    hawk.event('user_login', { 
        user_id: 12345, 
        username: 'john_doe',
        ip: '192.168.1.100' 
    }, { level: 'INFO', category: 'authentication' });

    hawk.event('payment_processed', { 
        amount: 99.99, 
        currency: 'USD',
        transaction_id: 'txn_123456' 
    }, { level: 'SUCCESS', category: 'payment' });

    await sleep(1000);

    // Layer 2: Classes and structured integration
    console.log('\n=== Layer 2: Monitor Class ===');
    
    const monitor = new hawk.Monitor('demo-service');
    
    // Simulate some application activity
    for (let i = 0; i < 15; i++) {
        monitor.counter('api_requests_total', 1, { endpoint: '/api/users' });
        monitor.counter('api_requests_total', Math.random() > 0.8 ? 1 : 0, { endpoint: '/api/orders' });
        
        monitor.gauge('response_time_ms', 50 + Math.random() * 200, { endpoint: '/api/users' });
        monitor.gauge('queue_size', Math.floor(Math.random() * 10), { service: 'background_jobs' });
        
        // Timer example
        const timer = monitor.timer('database_query');
        await sleep(50 + Math.random() * 100); // Simulate DB query
        timer.stop();
        
        monitor.uptime();
        await sleep(200);
    }

    // Dashboard demonstration
    console.log('\n=== Dashboard Demo ===');
    
    const dashboard = new hawk.Dashboard('System Overview');
    
    dashboard.addMetric('CPU Usage', () => 25 + Math.random() * 20, {
        format: '{:.1f}%',
        refreshRate: 1000
    });
    
    dashboard.addMetric('Memory Usage', () => 60 + Math.random() * 15, {
        format: '{:.1f}%',
        refreshRate: 1000
    });
    
    dashboard.addChart('Network I/O', () => {
        return Array.from({ length: 10 }, (_, i) => ({
            timestamp: Date.now() - (9 - i) * 1000,
            value: Math.random() * 100
        }));
    }, { type: 'line', maxPoints: 20 });
    
    dashboard.addTable('Active Processes', () => {
        return [
            { pid: 1234, name: 'node', cpu: '15.2%', memory: '120MB' },
            { pid: 5678, name: 'hawk', cpu: '2.1%', memory: '45MB' },
            { pid: 9012, name: 'nginx', cpu: '1.8%', memory: '80MB' }
        ];
    }, { columns: ['PID', 'Name', 'CPU', 'Memory'], maxRows: 10 });

    await sleep(2000);

    // Context management
    console.log('\n=== Context Management ===');
    
    // Function form
    await hawk.context('Database Migration', async () => {
        hawk.log('Starting user table migration');
        await sleep(500);
        hawk.log('Migrating 10,000 user records');
        await sleep(1000);
        hawk.log('User table migration completed');
    });
    
    // Object form
    const ctx = hawk.context('File Upload Processing');
    ctx.start();
    
    hawk.log('Validating file format');
    await sleep(300);
    hawk.log('Scanning for viruses');
    await sleep(500);
    hawk.log('Uploading to cloud storage');
    await sleep(800);
    hawk.log('File upload completed successfully');
    
    ctx.end();

    // Decorator simulation (in real usage, this would be class methods)
    console.log('\n=== Function Monitoring ===');
    
    const monitoredFunction = hawk.monitor(async function processPayment(amount, currency) {
        hawk.log(`Processing payment: ${amount} ${currency}`);
        
        // Simulate payment processing
        await sleep(200 + Math.random() * 500);
        
        if (Math.random() > 0.1) { // 90% success rate
            hawk.log(`Payment processed successfully: ${amount} ${currency}`, { level: 'SUCCESS' });
            return { success: true, transaction_id: `txn_${Date.now()}` };
        } else {
            throw new Error('Payment processing failed');
        }
    });
    
    // Run monitored function several times
    for (let i = 0; i < 5; i++) {
        try {
            await monitoredFunction(99.99, 'USD');
        } catch (error) {
            hawk.log(`Payment failed: ${error.message}`, { level: 'ERROR' });
        }
        await sleep(300);
    }

    // Simulate some error conditions
    console.log('\n=== Error Simulation ===');
    
    hawk.log('Simulating various error conditions...', { level: 'WARN' });
    
    try {
        throw new Error('Database connection timeout');
    } catch (error) {
        hawk.log(`Database error: ${error.message}`, { 
            level: 'ERROR',
            metadata: { 
                error_type: 'connection_timeout',
                retry_count: 3,
                last_attempt: new Date().toISOString()
            }
        });
    }
    
    // Simulate rate limiting
    hawk.log('Testing rate limiting...', { level: 'DEBUG' });
    for (let i = 0; i < 10; i++) {
        hawk.metric('rate_limit_test', i);
    }

    console.log('\n=== Demo Complete ===');
    hawk.log('Node.js demo completed successfully!', { level: 'SUCCESS' });
    
    // Give time for final messages to flush
    await sleep(1000);
    hawk.flush();
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Handle errors gracefully
process.on('uncaughtException', (error) => {
    console.error('Uncaught Exception:', error);
    hawk.flush();
    process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
    console.error('Unhandled Rejection:', reason);
    hawk.flush();
    process.exit(1);
});

// Run the demo
if (require.main === module) {
    runDemo().catch(console.error);
}