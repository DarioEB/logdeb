package logdbar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type LoggerData struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

const testDir = "test_logs"

func cleanup() {
	os.RemoveAll(testDir)
}

func TestNew(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:          testDir,
		InfoEnabled:  true,
		ErrorEnabled: true,
		WarnEnabled:  true,
		DebugEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	expectedFiles := []string{"info.log", "error.log", "warn.log", "debug.log"}
	for _, f := range expectedFiles {
		path := filepath.Join(testDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("expected file %s to exist", path)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Dir != "logs" {
		t.Errorf("expected Dir to be 'logs', got %s", cfg.Dir)
	}
	if !cfg.InfoEnabled || !cfg.ErrorEnabled || !cfg.WarnEnabled || !cfg.DebugEnabled {
		t.Error("expected all levels to be enabled by default")
	}
}

func TestPartialLevels(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:          testDir,
		InfoEnabled:  true,
		ErrorEnabled: false,
		WarnEnabled:  false,
		DebugEnabled: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	if _, err := os.Stat(filepath.Join(testDir, "info.log")); os.IsNotExist(err) {
		t.Fatal("info.log should exist")
	}

	if _, err := os.Stat(filepath.Join(testDir, "error.log")); !os.IsNotExist(err) {
		t.Fatal("error.log should not exist")
	}
}

func TestInfoLogger(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:         testDir,
		InfoEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := "testing logger info"
	logger.Info(data)

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(testDir, "info.log"))
	if err != nil {
		t.Fatal(err)
	}

	var logData LoggerData
	json.Unmarshal(content, &logData)

	if logData.Msg != data {
		t.Fatalf("expected message %q, got %q", data, logData.Msg)
	}
	if logData.Level != "INFO" {
		t.Fatalf("expected level INFO, got %s", logData.Level)
	}

	testTime, err := time.Parse(time.RFC3339, logData.Time)
	if err != nil {
		t.Fatal("invalid time format")
	}
	if time.Now().Before(testTime) {
		t.Fatal("log time is in the future")
	}
}

func TestErrorLogger(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:          testDir,
		ErrorEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := "testing logger error"
	logger.Error(data)

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(testDir, "error.log"))
	if err != nil {
		t.Fatal(err)
	}

	var logData LoggerData
	json.Unmarshal(content, &logData)

	if logData.Msg != data {
		t.Fatalf("expected message %q, got %q", data, logData.Msg)
	}
	if logData.Level != "ERROR" {
		t.Fatalf("expected level ERROR, got %s", logData.Level)
	}
}

func TestWarnLogger(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:         testDir,
		WarnEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := "testing logger warn"
	logger.Warn(data)

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(testDir, "warn.log"))
	if err != nil {
		t.Fatal(err)
	}

	var logData LoggerData
	json.Unmarshal(content, &logData)

	if logData.Msg != data {
		t.Fatalf("expected message %q, got %q", data, logData.Msg)
	}
	if logData.Level != "WARN" {
		t.Fatalf("expected level WARN, got %s", logData.Level)
	}
}

func TestDebugLogger(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:          testDir,
		DebugEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	data := "testing logger debug"
	logger.Debug(data)

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(testDir, "debug.log"))
	if err != nil {
		t.Fatal(err)
	}

	var logData LoggerData
	json.Unmarshal(content, &logData)

	if logData.Msg != data {
		t.Fatalf("expected message %q, got %q", data, logData.Msg)
	}
	if logData.Level != "DEBUG" {
		t.Fatalf("expected level DEBUG, got %s", logData.Level)
	}
}

func TestLoggerWithAttrs(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:         testDir,
		InfoEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("user login", "user_id", 123, "ip", "192.168.1.1")

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(testDir, "info.log"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "user_id") {
		t.Fatal("expected attrs to be in log output")
	}
	if !strings.Contains(string(content), "123") {
		t.Fatal("expected user_id value in log output")
	}
}

func TestDisabledLevelNoFile(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:          testDir,
		InfoEnabled:  true,
		ErrorEnabled: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Should not panic when calling disabled level
	logger.Error("this should not crash")

	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveLogFiles(t *testing.T) {
	cleanup()
	defer cleanup()

	logger, err := New(Config{
		Dir:         testDir,
		InfoEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("test")

	if err := logger.RemoveLogFiles(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(testDir, "info.log")); !os.IsNotExist(err) {
		t.Fatal("info.log should have been removed")
	}
}
