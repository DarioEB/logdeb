// Package logdbar provides a simple multi-level logger that writes to both
// console and JSON files. Each log level (Info, Error, Warn, Debug) can be
// independently enabled and writes to its own file.
package logdbar

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// Config holds the configuration for creating a new Logdbar instance.
type Config struct {
	// Dir is the directory where log files will be created.
	// If it doesn't exist, it will be created with 0755 permissions.
	Dir string

	// Enable specific log levels. Each enabled level will create
	// a corresponding file in Dir (info.log, error.log, etc.)
	InfoEnabled  bool
	ErrorEnabled bool
	WarnEnabled  bool
	DebugEnabled bool
}

// DefaultConfig returns a Config with all log levels enabled
// and "logs" as the default directory.
func DefaultConfig() Config {
	return Config{
		Dir:          "logs",
		InfoEnabled:  true,
		ErrorEnabled: true,
		WarnEnabled:  true,
		DebugEnabled: true,
	}
}

// LogInfo holds the file handle and logger for a specific log level.
type LogInfo struct {
	Path   string
	File   *os.File
	Logger *slog.Logger
}

// Logdbar is a multi-level logger that writes to console and JSON files.
type Logdbar struct {
	info    *LogInfo
	err     *LogInfo
	console *LogInfo
	warn    *LogInfo
	debug   *LogInfo
}

// New creates a new Logdbar instance based on the provided configuration.
// Returns an error if the directory cannot be created or files cannot be opened.
func New(cfg Config) (*Logdbar, error) {
	if cfg.Dir == "" {
		cfg.Dir = "logs"
	}

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating logs directory: %w", err)
	}

	logger := &Logdbar{
		console: &LogInfo{Logger: slog.New(slog.Default().Handler())},
	}

	type levelConfig struct {
		enabled  bool
		filename string
		level    slog.Level
		dest     **LogInfo
	}

	levels := []levelConfig{
		{cfg.InfoEnabled, "info.log", slog.LevelInfo, &logger.info},
		{cfg.ErrorEnabled, "error.log", slog.LevelError, &logger.err},
		{cfg.WarnEnabled, "warn.log", slog.LevelWarn, &logger.warn},
		{cfg.DebugEnabled, "debug.log", slog.LevelDebug, &logger.debug},
	}

	for _, lc := range levels {
		if !lc.enabled {
			continue
		}

		path := filepath.Join(cfg.Dir, lc.filename)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Close()
			return nil, fmt.Errorf("error opening %s: %w", path, err)
		}

		*lc.dest = &LogInfo{
			Path:   path,
			File:   file,
			Logger: slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{Level: lc.level})),
		}
	}

	return logger, nil
}

// Close closes all open log file handles.
// It's recommended to defer Close() after creating a new Logdbar.
func (l *Logdbar) Close() error {
	var errs []error
	logInfos := []*LogInfo{l.info, l.err, l.warn, l.debug}

	for _, li := range logInfos {
		if li != nil && li.File != nil {
			if err := li.File.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing log files: %v", errs)
	}
	return nil
}

// RemoveLogFiles closes all file handles and removes the log files.
// Use with caution as this permanently deletes the log files.
func (l *Logdbar) RemoveLogFiles() error {
	if err := l.Close(); err != nil {
		return err
	}

	var errs []error
	logInfos := []*LogInfo{l.info, l.err, l.warn, l.debug}

	for _, li := range logInfos {
		if li != nil && li.Path != "" {
			if err := os.Remove(li.Path); err != nil && !os.IsNotExist(err) {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors removing log files: %v", errs)
	}
	return nil
}

// Info logs a message at INFO level to console and info.log file.
// Optional attrs can be provided as key-value pairs.
func (l *Logdbar) Info(message string, attrs ...any) {
	l.console.Logger.Info(message, attrs...)
	if l.info != nil {
		l.info.Logger.Info(message, attrs...)
	}
}

// Error logs a message at ERROR level to console and error.log file.
// Optional attrs can be provided as key-value pairs.
func (l *Logdbar) Error(message string, attrs ...any) {
	l.console.Logger.Error(message, attrs...)
	if l.err != nil {
		l.err.Logger.Error(message, attrs...)
	}
}

// Debug logs a message at DEBUG level to console and debug.log file.
// Optional attrs can be provided as key-value pairs.
// Note: Console output depends on the default slog handler's minimum level.
func (l *Logdbar) Debug(message string, attrs ...any) {
	l.console.Logger.Debug(message, attrs...)
	if l.debug != nil {
		l.debug.Logger.Debug(message, attrs...)
	}
}

// Warn logs a message at WARN level to console and warn.log file.
// Optional attrs can be provided as key-value pairs.
func (l *Logdbar) Warn(message string, attrs ...any) {
	l.console.Logger.Warn(message, attrs...)
	if l.warn != nil {
		l.warn.Logger.Warn(message, attrs...)
	}
}
