package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hawk-tui/hawk-tui/internal/tui"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	// Parse command line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		appName     = flag.String("app", "hawk-tui", "Application name")
		logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
		configFile  = flag.String("config", "", "Configuration file path")
		debug       = flag.Bool("debug", false, "Enable debug mode")
	)

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("Hawk TUI v%s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", buildDate)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Create configuration
	config := tui.Config{
		AppName:    *appName,
		LogLevel:   *logLevel,
		ConfigFile: *configFile,
		Debug:      *debug,
	}

	// Create TUI model
	model, err := tui.NewModel(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TUI: %v\n", err)
		os.Exit(1)
	}

	// Create and start the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Hawk TUI - Universal Terminal UI Framework")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  hawk [OPTIONS]")
	fmt.Println("  your-app | hawk [OPTIONS]")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Hawk TUI is a universal terminal UI framework that displays structured")
	fmt.Println("  data from any application via JSON-RPC 2.0 protocol over stdin/stdout.")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --version          Show version information")
	fmt.Println("  --help             Show this help message")
	fmt.Println("  --app NAME         Application name (default: hawk-tui)")
	fmt.Println("  --log-level LEVEL  Log level: debug, info, warn, error (default: info)")
	fmt.Println("  --config FILE      Configuration file path")
	fmt.Println("  --debug            Enable debug mode")
	fmt.Println()
	fmt.Println("JSON-RPC METHODS:")
	fmt.Println("  hawk.log           Send log messages")
	fmt.Println("  hawk.metric        Send metrics data")
	fmt.Println("  hawk.config        Send configuration updates")
	fmt.Println("  hawk.progress      Send progress updates")
	fmt.Println("  hawk.dashboard     Send dashboard widgets")
	fmt.Println("  hawk.event         Send application events")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Run a demo application through Hawk TUI")
	fmt.Println("  python examples/python/demo.py | hawk")
	fmt.Println()
	fmt.Println("  # Run with custom application name")
	fmt.Println("  node examples/nodejs/demo.js | hawk --app my-app")
	fmt.Println()
	fmt.Println("  # Run with debug logging")
	fmt.Println("  your-app | hawk --debug --log-level debug")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/hawk-tui/hawk-tui")
}
