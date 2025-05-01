package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalPort := os.Getenv("PORT")
	originalLogLevel := os.Getenv("LOG_LEVEL")

	// Clean up after test
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("LOG_LEVEL", originalLogLevel)
	}()

	tests := []struct {
		name             string
		portEnv          string
		logLevelEnv      string
		expectedPort     int
		expectedLogLevel string
	}{
		{
			name:             "default values",
			portEnv:          "",
			logLevelEnv:      "",
			expectedPort:     8080,
			expectedLogLevel: "info",
		},
		{
			name:             "custom values",
			portEnv:          "3000",
			logLevelEnv:      "debug",
			expectedPort:     3000,
			expectedLogLevel: "debug",
		},
		{
			name:             "invalid port",
			portEnv:          "invalid",
			logLevelEnv:      "warn",
			expectedPort:     8080, // should fallback to default
			expectedLogLevel: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("PORT", tt.portEnv)
			os.Setenv("LOG_LEVEL", tt.logLevelEnv)

			// Load configuration
			cfg := Load()

			// Verify port
			if cfg.Port != tt.expectedPort {
				t.Errorf("Port = %v, want %v", cfg.Port, tt.expectedPort)
			}

			// Verify log level
			if cfg.LogLevel != tt.expectedLogLevel {
				t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, tt.expectedLogLevel)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Save original environment variable
	originalValue := os.Getenv("TEST_KEY")

	// Clean up after test
	defer func() {
		os.Setenv("TEST_KEY", originalValue)
	}()

	tests := []struct {
		name          string
		key           string
		defaultValue  string
		envValue      string
		expectedValue string
	}{
		{
			name:          "environment variable set",
			key:           "TEST_KEY",
			defaultValue:  "default",
			envValue:      "custom",
			expectedValue: "custom",
		},
		{
			name:          "environment variable not set",
			key:           "TEST_KEY",
			defaultValue:  "default",
			envValue:      "",
			expectedValue: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv(tt.key, tt.envValue)

			// Get value
			value := getEnv(tt.key, tt.defaultValue)

			// Verify value
			if value != tt.expectedValue {
				t.Errorf("getEnv(%q, %q) = %v, want %v", tt.key, tt.defaultValue, value, tt.expectedValue)
			}
		})
	}
}
