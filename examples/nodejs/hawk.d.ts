/**
 * Hawk TUI Node.js Client Library TypeScript Definitions
 * 
 * Provides comprehensive type definitions for the Hawk TUI client library,
 * enabling full TypeScript support with IntelliSense and type checking.
 */

declare module 'hawk-tui-client' {
    // Configuration interfaces
    interface HawkConfig {
        enabled?: boolean;
        appName?: string;
        debug?: boolean;
        batchSize?: number;
        batchTimeout?: number;
        maxMessageRate?: number;
        gracefulFallback?: boolean;
    }

    interface AutoOptions extends HawkConfig {
        [key: string]: any;
    }

    // Message interfaces
    interface LogOptions {
        level?: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'SUCCESS';
        context?: string | null;
        metadata?: Record<string, any> | null;
        stack?: string;
    }

    interface MetricOptions {
        type?: 'gauge' | 'counter' | 'histogram';
        unit?: string;
        tags?: Record<string, string>;
        format?: string | null;
    }

    interface ConfigOptions {
        value?: any;
        default?: any;
        type?: 'string' | 'number' | 'boolean' | 'array' | 'object' | 'auto';
        description?: string;
        required?: boolean;
        validation?: any;
    }

    interface ProgressOptions {
        message?: string;
        status?: 'in_progress' | 'completed' | 'failed';
    }

    interface EventOptions {
        level?: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'SUCCESS';
        category?: string;
        action?: string | null;
    }

    // Timer interface
    interface Timer {
        stop(): number;
    }

    // Widget interfaces
    interface Widget {
        type: string;
        name: string;
        data: any;
        timestamp: string;
    }

    interface DashboardMetricOptions {
        refreshRate?: number;
        format?: string | null;
    }

    interface DashboardChartOptions {
        type?: 'line' | 'bar' | 'area' | 'pie';
        maxPoints?: number;
    }

    interface DashboardTableOptions {
        columns?: string[];
        maxRows?: number;
    }

    // Context manager interface
    interface ContextManager {
        name: string;
        start(): void;
        end(): void;
    }

    // Express middleware types
    interface Request {
        method: string;
        path: string;
        route?: { path: string };
        get(name: string): string | undefined;
        ip: string;
    }

    interface Response {
        statusCode: number;
        send(data: any): void;
    }

    type NextFunction = () => void;
    type MiddlewareFunction = (req: Request, res: Response, next: NextFunction) => void;

    // Monitor class
    class Monitor {
        constructor(appName: string, options?: HawkConfig);
        
        counter(name: string, value?: number, tags?: Record<string, string>): number;
        gauge(name: string, value: number, tags?: Record<string, string>): number;
        timer(name: string): Timer;
        uptime(): number;
    }

    // Dashboard class
    class Dashboard {
        constructor(name: string, options?: any);
        
        addWidget(type: string, name: string, data?: any): Widget;
        addMetric(name: string, getValue: () => number, options?: DashboardMetricOptions): Widget;
        addChart(name: string, getDataPoints: () => any[], options?: DashboardChartOptions): Widget;
        addTable(name: string, getRows: () => any[], options?: DashboardTableOptions): Widget;
    }

    // Layer 0: Magic mode
    function auto(options?: AutoOptions): void;

    // Layer 1: Simple functions
    function log(message: string, options?: LogOptions): void;
    function metric(name: string, value: number, options?: MetricOptions): void;
    function config<T = any>(key: string, options?: ConfigOptions): T;
    function progress(name: string, current: number, total: number, options?: ProgressOptions): void;
    function event(name: string, data?: any, options?: EventOptions): void;

    // Layer 2: Middleware
    function middleware(options?: HawkConfig): MiddlewareFunction;

    // Decorators
    function monitor<T extends (...args: any[]) => any>(fn: T): T;
    function monitor(target: any, propertyKey: string, descriptor: PropertyDescriptor): PropertyDescriptor;
    
    function timed(name: string): (target: any, propertyKey: string, descriptor: PropertyDescriptor) => PropertyDescriptor;

    // Context management
    function context<T>(name: string, fn: () => T): T;
    function context<T>(name: string, fn: () => Promise<T>): Promise<T>;
    function context(name: string): ContextManager;

    // Utilities
    function configure(options: HawkConfig): void;
    function isEnabled(): boolean;
    function disable(): void;
    function enable(): void;
    function flush(): void;

    // Version
    const version: string;

    // Default export (for ES modules)
    const hawk: {
        auto: typeof auto;
        log: typeof log;
        metric: typeof metric;
        config: typeof config;
        progress: typeof progress;
        event: typeof event;
        Monitor: typeof Monitor;
        Dashboard: typeof Dashboard;
        middleware: typeof middleware;
        monitor: typeof monitor;
        timed: typeof timed;
        context: typeof context;
        configure: typeof configure;
        isEnabled: typeof isEnabled;
        disable: typeof disable;
        enable: typeof enable;
        flush: typeof flush;
        version: string;
    };

    export = hawk;
}

// Global augmentation for when used without module system
declare global {
    namespace NodeJS {
        interface Global {
            hawk?: typeof import('hawk-tui-client');
        }
    }
}