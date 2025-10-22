package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hawk-tui/hawk-tui/pkg/hawktui/styles"
)

// Level represents a log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSuccess
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelSuccess:
		return "SUCCESS"
	default:
		return "UNKNOWN"
	}
}

// Entry represents a log entry
type Entry struct {
	Time    time.Time
	Level   Level
	Message string
	Fields  map[string]interface{}
}

// Logger provides structured, styled logging
type Logger struct {
	theme      *styles.Theme
	minLevel   Level
	output     io.Writer
	entries    []Entry
	maxEntries int
	mu         sync.RWMutex
	showTime   bool
	showLevel  bool
	colorize   bool
}

// New creates a new logger
func New() *Logger {
	return &Logger{
		theme:      styles.DefaultTheme(),
		minLevel:   LevelInfo,
		output:     os.Stdout,
		entries:    make([]Entry, 0),
		maxEntries: 1000,
		showTime:   true,
		showLevel:  true,
		colorize:   true,
	}
}

// NewWithTheme creates a new logger with a specific theme
func NewWithTheme(theme *styles.Theme) *Logger {
	return &Logger{
		theme:      theme,
		minLevel:   LevelInfo,
		output:     os.Stdout,
		entries:    make([]Entry, 0),
		maxEntries: 1000,
		showTime:   true,
		showLevel:  true,
		colorize:   true,
	}
}

// SetTheme sets the theme
func (l *Logger) SetTheme(theme *styles.Theme) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.theme = theme
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetMaxEntries sets the maximum number of entries to keep
func (l *Logger) SetMaxEntries(max int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.maxEntries = max
}

// SetShowTime sets whether to show timestamps
func (l *Logger) SetShowTime(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showTime = show
}

// SetShowLevel sets whether to show log levels
func (l *Logger) SetShowLevel(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showLevel = show
}

// SetColorize sets whether to colorize output
func (l *Logger) SetColorize(colorize bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colorize = colorize
}

// log is the internal logging function
func (l *Logger) log(level Level, message string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.minLevel {
		return
	}

	entry := Entry{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Fields:  fields,
	}

	// Add to entries
	l.entries = append(l.entries, entry)
	if len(l.entries) > l.maxEntries {
		l.entries = l.entries[len(l.entries)-l.maxEntries:]
	}

	// Format and write
	formatted := l.format(entry)
	if l.output != nil {
		fmt.Fprintln(l.output, formatted)
	}
}

// format formats a log entry
func (l *Logger) format(entry Entry) string {
	var parts []string

	// Timestamp
	if l.showTime {
		timeStyle := lipgloss.NewStyle().
			Foreground(l.theme.TextMuted)
		timestamp := entry.Time.Format("2006-01-02 15:04:05")
		if l.colorize {
			parts = append(parts, timeStyle.Render(timestamp))
		} else {
			parts = append(parts, timestamp)
		}
	}

	// Level
	if l.showLevel {
		levelStr := fmt.Sprintf("%-7s", entry.Level.String())
		if l.colorize {
			levelStyle := l.getLevelStyle(entry.Level)
			parts = append(parts, levelStyle.Render(levelStr))
		} else {
			parts = append(parts, levelStr)
		}
	}

	// Message
	if l.colorize {
		messageStyle := lipgloss.NewStyle().
			Foreground(l.theme.TextPrimary)
		parts = append(parts, messageStyle.Render(entry.Message))
	} else {
		parts = append(parts, entry.Message)
	}

	// Fields
	if len(entry.Fields) > 0 {
		var fieldStrs []string
		for k, v := range entry.Fields {
			fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%v", k, v))
		}
		fieldsStr := strings.Join(fieldStrs, " ")
		if l.colorize {
			fieldStyle := lipgloss.NewStyle().
				Foreground(l.theme.TextSecondary).
				Italic(true)
			parts = append(parts, fieldStyle.Render(fieldsStr))
		} else {
			parts = append(parts, fieldsStr)
		}
	}

	return strings.Join(parts, " ")
}

// getLevelStyle returns the style for a log level
func (l *Logger) getLevelStyle(level Level) lipgloss.Style {
	baseStyle := lipgloss.NewStyle().Bold(true)

	switch level {
	case LevelDebug:
		return baseStyle.Foreground(l.theme.Debug)
	case LevelInfo:
		return baseStyle.Foreground(l.theme.Info)
	case LevelWarn:
		return baseStyle.Foreground(l.theme.Warning)
	case LevelError:
		return baseStyle.Foreground(l.theme.Error)
	case LevelSuccess:
		return baseStyle.Foreground(l.theme.Success)
	default:
		return baseStyle.Foreground(l.theme.TextMuted)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelDebug, message, f)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelInfo, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelWarn, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelError, message, f)
}

// Success logs a success message
func (l *Logger) Success(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelSuccess, message, f)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, args...), nil)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, args...), nil)
}

// Successf logs a formatted success message
func (l *Logger) Successf(format string, args ...interface{}) {
	l.log(LevelSuccess, fmt.Sprintf(format, args...), nil)
}

// GetEntries returns all log entries
func (l *Logger) GetEntries() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	entries := make([]Entry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

// Clear clears all log entries
func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make([]Entry, 0)
}
