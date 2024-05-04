package cli

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestMyConfigSetUrlFile(t *testing.T) {
	t.Parallel()

	config := NewConfigWithValues("db_path", "url_file", "log_level")

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to get current working directory")
	}
	good_url_file := filepath.Join(dir, "cli_test.go")

	tcs := []struct {
		name     string
		url_file string
		expected error
	}{
		{"Empty Path", "", errors.New("Must provide path to url file")},
		{"File Does Not Exist", "Not_a_File", errors.New("The urls file was not found: stat Not_a_File: no such file or directory")},
		{"No Errors", good_url_file, nil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := config.SetUrlFile(tc.url_file)

			if err != nil {
				if err.Error() != tc.expected.Error() {
					t.Fatalf("Expected %v, but got %v", tc.expected, err)
				}
			} else {
				if config.GetUrlFile() != good_url_file {
					t.Fatalf("Expected %s, but got %s", good_url_file, config.GetUrlFile())
				}
			}
		})

	}
}

func TestMyConfigLogLevel(t *testing.T) {
	t.Parallel()

	var test_config Config

	tcs := []struct {
		name      string
		log_level string
		expected  slog.Level
	}{
		{"Debug Lowercase", "debug", slog.LevelDebug},
		{"Debug Uppercase", "DEBUG", slog.LevelDebug},
		{"Debug MixCase", "DeBuG", slog.LevelDebug},
		{"Info Lowercase", "info", slog.LevelInfo},
		{"Info Uppercase", "INFO", slog.LevelInfo},
		{"Info MixCase", "InFo", slog.LevelInfo},
		{"Warn Lowercase", "warn", slog.LevelWarn},
		{"Warn Uppercase", "WARN", slog.LevelWarn},
		{"Warn MixCase", "WaRn", slog.LevelWarn},
		{"Error Lowercase", "error", slog.LevelError},
		{"Error Uppercase", "ERROR", slog.LevelError},
		{"Error MixCase", "ErRoR", slog.LevelError},
		{"Invalid Level", "Not_a_Level", slog.LevelInfo},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			test_config = NewConfigWithValues("db_path", "url_files", tc.log_level)

			if test_config.LogLevel() != tc.expected {
				t.Fatalf("Expected Log Level %v, but got %v", tc.expected, test_config.LogLevel())
			}

		})
	}
}

func TestNewConfigWithValues(t *testing.T) {
	t.Parallel()

	db_path := "db_path"
	url_file := "url_file"
	log_level := "warn"

	config := NewConfigWithValues(db_path, url_file, log_level)

	if config.GetDBPath() != db_path {
		t.Fatalf("Expected %s, but got %s", db_path, config.GetDBPath())
	}

	if config.GetUrlFile() != url_file {
		t.Fatalf("Expected %s, but got %s", url_file, config.GetUrlFile())
	}

	if config.LogLevel() != slog.LevelWarn {
		t.Fatalf("Expected %v, but got %v", config.LogLevel(), slog.LevelWarn)
	}
}

func TestNewConfig(t *testing.T) {
	test_config := NewConfig()

	if test_config == nil {
		t.Fatal("Returned config should not be equal to nil")
	}
}
