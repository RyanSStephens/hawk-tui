#!/usr/bin/env node
/**
 * Hawk TUI Node.js Client Test Suite
 * 
 * Comprehensive tests for all layers of the Hawk TUI Node.js client library.
 * Tests functionality without requiring the TUI to be running.
 */

const assert = require('assert');
const { performance } = require('perf_hooks');
const hawk = require('./hawk');

// Test state
let testResults = [];
let originalStdoutWrite;
let capturedMessages = [];

// Capture stdout to verify messages
function captureStdout() {
    originalStdoutWrite = process.stdout.write;
    process.stdout.write = function(chunk) {
        if (chunk.includes('"jsonrpc":"2.0"')) {
            try {
                capturedMessages.push(JSON.parse(chunk.trim()));
            } catch (e) {
                // Ignore malformed JSON
            }
        }
        return originalStdoutWrite.call(this, chunk);
    };
}

function restoreStdout() {
    if (originalStdoutWrite) {
        process.stdout.write = originalStdoutWrite;
    }
}

function clearCapturedMessages() {
    capturedMessages = [];
}

// Test utilities
function test(name, fn) {
    console.log(`\n🧪 ${name}`);
    try {
        fn();
        console.log('✅ PASS');
        testResults.push({ name, status: 'PASS' });
    } catch (error) {
        console.log(`❌ FAIL: ${error.message}`);
        testResults.push({ name, status: 'FAIL', error: error.message });
    }
}

async function asyncTest(name, fn) {
    console.log(`\n🧪 ${name}`);
    try {
        await fn();
        console.log('✅ PASS');
        testResults.push({ name, status: 'PASS' });
    } catch (error) {
        console.log(`❌ FAIL: ${error.message}`);
        testResults.push({ name, status: 'FAIL', error: error.message });
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Tests
async function runTests() {
    console.log('🚀 Starting Hawk TUI Node.js Client Test Suite\n');
    
    // Setup
    captureStdout();
    hawk.configure({ enabled: true, batchTimeout: 50 });
    
    // Test Layer 0: Configuration
    test('Configuration - Basic config', () => {
        hawk.configure({ appName: 'test-app', debug: true });
        assert(hawk.isEnabled() === true);
    });
    
    test('Configuration - Enable/Disable', () => {
        hawk.disable();
        assert(hawk.isEnabled() === false);
        hawk.enable();
        assert(hawk.isEnabled() === true);
    });
    
    // Test Layer 1: Simple Functions
    await asyncTest('Layer 1 - Basic logging', async () => {
        clearCapturedMessages();
        hawk.configure({ debug: false }); // Disable debug to reduce noise
        hawk.log('Test message', { level: 'INFO' });
        await sleep(100); // Wait for batch
        
        assert(capturedMessages.length > 0);
        const logMessage = capturedMessages.find(m => m.params.message === 'Test message');
        assert(logMessage !== undefined, 'Could not find test message');
        assert(logMessage.method === 'hawk.log');
        assert(logMessage.params.level === 'INFO');
    });
    
    await asyncTest('Layer 1 - Metrics', async () => {
        clearCapturedMessages();
        hawk.metric('test_metric', 42.5, { unit: 'ms', type: 'gauge' });
        await sleep(100);
        
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.metric');
        assert(message.params.name === 'test_metric');
        assert(message.params.value === 42.5);
        assert(message.params.unit === 'ms');
        assert(message.params.type === 'gauge');
    });
    
    await asyncTest('Layer 1 - Configuration values', async () => {
        clearCapturedMessages();
        const value = hawk.config('test_port', { default: 3000, type: 'number' });
        await sleep(100);
        
        assert(value === 3000);
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.config');
        assert(message.params.key === 'test_port');
        assert(message.params.value === 3000);
    });
    
    await asyncTest('Layer 1 - Progress tracking', async () => {
        clearCapturedMessages();
        hawk.progress('file_upload', 75, 100, { message: 'Uploading...' });
        await sleep(100);
        
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.progress');
        assert(message.params.name === 'file_upload');
        assert(message.params.current === 75);
        assert(message.params.total === 100);
        assert(message.params.percentage === 75);
        assert(message.params.message === 'Uploading...');
    });
    
    await asyncTest('Layer 1 - Events', async () => {
        clearCapturedMessages();
        hawk.event('user_login', { user_id: 123 }, { level: 'SUCCESS' });
        await sleep(100);
        
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.event');
        assert(message.params.name === 'user_login');
        assert(message.params.data.user_id === 123);
        assert(message.params.level === 'SUCCESS');
    });
    
    // Test Layer 2: Monitor Class
    await asyncTest('Layer 2 - Monitor class initialization', async () => {
        const monitor = new hawk.Monitor('test-service');
        assert(monitor.appName === 'test-service');
        assert(monitor.counters instanceof Map);
        assert(monitor.gauges instanceof Map);
    });
    
    await asyncTest('Layer 2 - Monitor counters', async () => {
        clearCapturedMessages();
        const monitor = new hawk.Monitor('test-service');
        const value = monitor.counter('requests', 5);
        await sleep(100);
        
        assert(value === 5);
        const value2 = monitor.counter('requests', 3);
        assert(value2 === 8); // Cumulative
        
        assert(capturedMessages.length > 0);
        const counterMessage = capturedMessages.find(m => m.params.name === 'requests');
        assert(counterMessage !== undefined, 'Could not find counter message');
        assert(counterMessage.method === 'hawk.metric');
        assert(counterMessage.params.type === 'counter');
    });
    
    await asyncTest('Layer 2 - Monitor gauges', async () => {
        clearCapturedMessages();
        const monitor = new hawk.Monitor('test-service');
        const value = monitor.gauge('cpu_usage', 67.5);
        await sleep(100);
        
        assert(value === 67.5);
        assert(capturedMessages.length > 0);
        const gaugeMessage = capturedMessages.find(m => m.params.name === 'cpu_usage');
        assert(gaugeMessage !== undefined, 'Could not find gauge message');
        assert(gaugeMessage.method === 'hawk.metric');
        assert(gaugeMessage.params.type === 'gauge');
        assert(gaugeMessage.params.value === 67.5);
    });
    
    await asyncTest('Layer 2 - Monitor timers', async () => {
        clearCapturedMessages();
        const monitor = new hawk.Monitor('test-service');
        const timer = monitor.timer('test_operation');
        
        assert(typeof timer.stop === 'function');
        
        await sleep(50);
        const duration = timer.stop();
        await sleep(100);
        
        assert(duration >= 45); // Should be around 50ms
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.metric');
        assert(message.params.type === 'histogram');
        assert(message.params.unit === 'ms');
    });
    
    await asyncTest('Layer 2 - Monitor uptime', async () => {
        clearCapturedMessages();
        const monitor = new hawk.Monitor('test-service');
        await sleep(100);
        const uptime = monitor.uptime();
        await sleep(100);
        
        assert(uptime >= 90); // Should be around 100ms
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.metric');
        assert(message.params.name === 'uptime_seconds');
    });
    
    // Test Layer 2: Dashboard Class
    await asyncTest('Layer 2 - Dashboard creation', async () => {
        clearCapturedMessages();
        const dashboard = new hawk.Dashboard('Test Dashboard');
        const widget = dashboard.addWidget('test', 'Test Widget', { value: 42 });
        await sleep(100);
        
        assert(dashboard.name === 'Test Dashboard');
        assert(widget.type === 'test');
        assert(widget.name === 'Test Widget');
        assert(widget.data.value === 42);
        
        assert(capturedMessages.length > 0);
        const message = capturedMessages[0];
        assert(message.method === 'hawk.dashboard');
        assert(message.params.dashboard_name === 'Test Dashboard');
    });
    
    // Test decorators
    await asyncTest('Function monitoring decorator', async () => {
        clearCapturedMessages();
        hawk.configure({ debug: true }); // Enable debug for function monitoring
        
        const testFunction = hawk.monitor(async function namedTestFunction() {
            await sleep(50);
            return 'result';
        });
        
        const result = await testFunction();
        await sleep(100);
        
        assert(result === 'result', `Expected 'result' but got ${result}`);
        
        // Check for execution logs
        const logMessages = capturedMessages.filter(m => m.method === 'hawk.log');
        const executingMessage = logMessages.find(m => m.params.message.includes('Executing'));
        const completedMessage = logMessages.find(m => m.params.message.includes('completed'));
        
        assert(executingMessage !== undefined, 'Could not find executing message');
        assert(completedMessage !== undefined, 'Could not find completed message');
    });
    
    await asyncTest('Function monitoring decorator with error', async () => {
        clearCapturedMessages();
        
        const errorFunction = hawk.monitor(function errorFunction() {
            throw new Error('Test error');
        });
        
        try {
            errorFunction();
            assert(false, 'Should have thrown error');
        } catch (error) {
            assert(error.message === 'Test error', `Expected 'Test error' but got ${error.message}`);
        }
        
        await sleep(100);
        
        // Should have error metrics and logs
        const metricMessages = capturedMessages.filter(m => m.method === 'hawk.metric');
        const errorMetrics = metricMessages.filter(m => m.params.name && m.params.name.includes('errors'));
        const logMessages = capturedMessages.filter(m => m.method === 'hawk.log');
        const errorLogs = logMessages.filter(m => m.params.message && m.params.message.includes('failed'));
        
        assert(errorMetrics.length > 0 || errorLogs.length > 0, 'Should have error metrics or logs');
    });
    
    // Test context management
    await asyncTest('Context management - function form', async () => {
        clearCapturedMessages();
        
        const result = await hawk.context('Test Context', async () => {
            await sleep(50);
            return 'context result';
        });
        
        await sleep(100);
        
        assert(result === 'context result');
        
        const logMessages = capturedMessages.filter(m => m.method === 'hawk.log');
        const enterMessage = logMessages.find(m => m.params.message.includes('Entering context'));
        const exitMessage = logMessages.find(m => m.params.message.includes('Exiting context'));
        
        assert(enterMessage !== undefined);
        assert(exitMessage !== undefined);
    });
    
    await asyncTest('Context management - object form', async () => {
        clearCapturedMessages();
        
        const ctx = hawk.context('Object Context');
        assert(typeof ctx.start === 'function');
        assert(typeof ctx.end === 'function');
        
        ctx.start();
        await sleep(50);
        ctx.end();
        await sleep(100);
        
        const logMessages = capturedMessages.filter(m => m.method === 'hawk.log');
        assert(logMessages.length >= 2);
    });
    
    // Test rate limiting
    await asyncTest('Rate limiting', async () => {
        hawk.configure({ maxMessageRate: 5 }); // Very low rate
        clearCapturedMessages();
        
        // Send more messages than rate limit
        for (let i = 0; i < 10; i++) {
            hawk.log(`Message ${i}`);
        }
        
        await sleep(100);
        
        // Should have fewer messages than sent due to rate limiting
        assert(capturedMessages.length <= 5);
        
        // Reset rate limit
        hawk.configure({ maxMessageRate: 1000 });
    });
    
    // Test batching
    await asyncTest('Message batching', async () => {
        hawk.configure({ batchSize: 3, batchTimeout: 1000 });
        clearCapturedMessages();
        
        // Send exactly batch size worth of messages
        hawk.log('Message 1');
        hawk.log('Message 2');
        hawk.log('Message 3');
        
        await sleep(50); // Should batch immediately
        
        assert(capturedMessages.length === 3);
        
        // Reset config
        hawk.configure({ batchSize: 100, batchTimeout: 50 });
    });
    
    // Test graceful fallback
    test('Graceful fallback when disabled', () => {
        hawk.disable();
        
        // These should not throw errors
        hawk.log('Should not crash');
        hawk.metric('test', 42);
        hawk.event('test', {});
        
        hawk.enable();
    });
    
    // Test auto mode (partial - can't fully test console hijacking in tests)
    test('Auto mode configuration', () => {
        // Should not throw
        hawk.auto({ debug: false });
        
        // Should be able to call again without issues
        hawk.auto({ appName: 'test-auto' });
    });
    
    // Test utilities
    test('Utility functions', () => {
        assert(typeof hawk.isEnabled === 'function');
        assert(typeof hawk.disable === 'function');
        assert(typeof hawk.enable === 'function');
        assert(typeof hawk.flush === 'function');
        assert(typeof hawk.configure === 'function');
    });
    
    // Test version
    test('Version information', () => {
        assert(typeof hawk.version === 'string');
        assert(hawk.version.length > 0);
    });
    
    // Test export structure
    test('Module exports', () => {
        const expectedExports = [
            'auto', 'log', 'metric', 'config', 'progress', 'event',
            'Monitor', 'Dashboard', 'middleware', 'monitor', 'timed', 'context',
            'configure', 'isEnabled', 'disable', 'enable', 'flush', 'version'
        ];
        
        for (const exportName of expectedExports) {
            assert(hawk[exportName] !== undefined, `Missing export: ${exportName}`);
        }
    });
    
    // Cleanup
    restoreStdout();
    
    // Report results
    console.log('\n📊 Test Results Summary:');
    console.log('========================');
    
    const passed = testResults.filter(r => r.status === 'PASS').length;
    const failed = testResults.filter(r => r.status === 'FAIL').length;
    
    console.log(`✅ Passed: ${passed}`);
    console.log(`❌ Failed: ${failed}`);
    console.log(`📋 Total:  ${testResults.length}`);
    
    if (failed > 0) {
        console.log('\n💥 Failed Tests:');
        testResults.filter(r => r.status === 'FAIL').forEach(test => {
            console.log(`   ${test.name}: ${test.error}`);
        });
        process.exit(1);
    } else {
        console.log('\n🎉 All tests passed!');
        process.exit(0);
    }
}

// Handle errors
process.on('uncaughtException', (error) => {
    console.error('\n💥 Uncaught Exception:', error);
    process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
    console.error('\n💥 Unhandled Rejection:', reason);
    process.exit(1);
});

// Run tests if this file is executed directly
if (require.main === module) {
    runTests().catch(console.error);
}

module.exports = { runTests };