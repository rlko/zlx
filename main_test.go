package main

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	// Test empty config
	config = Config{}
	errors := validateConfig()
	if len(errors) != 2 {
		t.Errorf("Expected 2 validation errors for empty config, got %d", len(errors))
	}

	// Test valid config
	config = Config{
		ServerName: "example.com",
		Token:      "valid-token",
	}
	errors = validateConfig()
	if len(errors) != 0 {
		t.Errorf("Expected no validation errors for valid config, got %d", len(errors))
	}
}
