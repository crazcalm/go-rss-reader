package cli

import (
	"testing"
	"log/slog"
	"os"
)

func TestMyConfigCliParseLogLevel(t *testing.T){
	var test_config Config

	//Needed to that flag.Parse sees the log-level flag
	os.Args = append(os.Args, "-log-level", "debug")

	// Needed to figure out which index to override
	args_index := len(os.Args) - 1
	
	tcs := []struct{
		name string
		log_level string
		expected slog.Level
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
		t.Run(tc.name, func(t * testing.T) {
			os.Args[args_index] = tc.log_level
			t.Logf("os.Args: %v", os.Args)

			test_config = NewConfig()

			err := test_config.CliParse()
			if err != nil {
				t.Fatalf("Was not expecting test_config.Parse() to return an error: %v", err)
			}

			if test_config.LogLevel() != tc.expected {
				t.Fatalf("Expected Log Level %v, but got %v", tc.expected, test_config.LogLevel())
			}
			
		})
	}
}

func TestNewConfig(t *testing.T) {
	test_config := NewConfig()

	if test_config == nil {
		t.Fatal("Returned config should not be equal to nil")
	}
}
