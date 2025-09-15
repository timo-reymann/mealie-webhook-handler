package validation

import (
	"fmt"
	"strings"
	"testing"
)

func TestRequireKey(t *testing.T) {
	tests := []struct {
		name      string
		config    map[string]string
		key       string
		wantError bool
		errorMsg  string
	}{
		{
			name: "key exists - should pass",
			config: map[string]string{
				"database_url": "localhost:5432",
				"api_key":      "secret123",
			},
			key:       "database_url",
			wantError: false,
		},
		{
			name: "key missing - should fail",
			config: map[string]string{
				"database_url": "localhost:5432",
			},
			key:       "api_key",
			wantError: true,
			errorMsg:  "missing required config key 'api_key'",
		},
		{
			name:      "empty config - should fail",
			config:    map[string]string{},
			key:       "required_key",
			wantError: true,
			errorMsg:  "missing required config key 'required_key'",
		},
		{
			name:      "nil config - should fail",
			config:    nil,
			key:       "required_key",
			wantError: true,
			errorMsg:  "missing required config key 'required_key'",
		},
		{
			name: "key exists with empty value - should pass",
			config: map[string]string{
				"empty_key": "",
			},
			key:       "empty_key",
			wantError: false,
		},
		{
			name: "special characters in key name",
			config: map[string]string{
				"key-with-dashes":      "value1",
				"key_with_underscores": "value2",
			},
			key:       "key-with-dashes",
			wantError: false,
		},
		{
			name: "missing special characters key",
			config: map[string]string{
				"normal_key": "value",
			},
			key:       "key-with-dashes",
			wantError: true,
			errorMsg:  "missing required config key 'key-with-dashes'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := RequireKey(tt.key)
			err := check(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("RequireKey() expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("RequireKey() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("RequireKey() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFailOnFirst(t *testing.T) {
	// Helper function to create a check that always passes
	alwaysPass := func(name string) Check {
		return func(config map[string]string) error {
			return nil
		}
	}

	// Helper function to create a check that always fails
	alwaysFail := func(name string) Check {
		return func(config map[string]string) error {
			return fmt.Errorf("check %s failed", name)
		}
	}

	tests := []struct {
		name      string
		checks    []Check
		config    map[string]string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "no checks - should pass",
			checks:    []Check{},
			config:    map[string]string{},
			wantError: false,
		},
		{
			name: "single passing check - should pass",
			checks: []Check{
				alwaysPass("check1"),
			},
			config:    map[string]string{},
			wantError: false,
		},
		{
			name: "single failing check - should fail",
			checks: []Check{
				alwaysFail("check1"),
			},
			config:    map[string]string{},
			wantError: true,
			errorMsg:  "check check1 failed",
		},
		{
			name: "multiple passing checks - should pass",
			checks: []Check{
				alwaysPass("check1"),
				alwaysPass("check2"),
				alwaysPass("check3"),
			},
			config:    map[string]string{},
			wantError: false,
		},
		{
			name: "first check fails - should fail immediately",
			checks: []Check{
				alwaysFail("check1"),
				alwaysPass("check2"),
				alwaysPass("check3"),
			},
			config:    map[string]string{},
			wantError: true,
			errorMsg:  "check check1 failed",
		},
		{
			name: "middle check fails - should fail on middle check",
			checks: []Check{
				alwaysPass("check1"),
				alwaysFail("check2"),
				alwaysPass("check3"),
			},
			config:    map[string]string{},
			wantError: true,
			errorMsg:  "check check2 failed",
		},
		{
			name: "last check fails - should fail on last check",
			checks: []Check{
				alwaysPass("check1"),
				alwaysPass("check2"),
				alwaysFail("check3"),
			},
			config:    map[string]string{},
			wantError: true,
			errorMsg:  "check check3 failed",
		},
		{
			name: "multiple RequireKey checks - all pass",
			checks: []Check{
				RequireKey("key1"),
				RequireKey("key2"),
				RequireKey("key3"),
			},
			config: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			wantError: false,
		},
		{
			name: "multiple RequireKey checks - first missing",
			checks: []Check{
				RequireKey("missing_key"),
				RequireKey("key2"),
				RequireKey("key3"),
			},
			config: map[string]string{
				"key2": "value2",
				"key3": "value3",
			},
			wantError: true,
			errorMsg:  "missing required config key 'missing_key'",
		},
		{
			name: "multiple RequireKey checks - middle missing",
			checks: []Check{
				RequireKey("key1"),
				RequireKey("missing_key"),
				RequireKey("key3"),
			},
			config: map[string]string{
				"key1": "value1",
				"key3": "value3",
			},
			wantError: true,
			errorMsg:  "missing required config key 'missing_key'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := FailOnFirst(tt.checks...)
			err := check(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("FailOnFirst() expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("FailOnFirst() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("FailOnFirst() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFailOnFirst_StopsOnFirstError(t *testing.T) {
	// This test specifically verifies that FailOnFirst stops on the first error
	// and doesn't execute subsequent checks
	var executedChecks []string

	createTrackingCheck := func(name string, shouldFail bool) Check {
		return func(config map[string]string) error {
			executedChecks = append(executedChecks, name)
			if shouldFail {
				return fmt.Errorf("check %s failed", name)
			}
			return nil
		}
	}

	t.Run("stops after first failure", func(t *testing.T) {
		executedChecks = []string{} // Reset tracking

		checks := []Check{
			createTrackingCheck("check1", false), // passes
			createTrackingCheck("check2", true),  // fails
			createTrackingCheck("check3", false), // should not be executed
		}

		failOnFirstCheck := FailOnFirst(checks...)
		err := failOnFirstCheck(map[string]string{})

		if err == nil {
			t.Errorf("FailOnFirst() expected error but got none")
		}

		expectedExecuted := []string{"check1", "check2"}
		if len(executedChecks) != len(expectedExecuted) {
			t.Errorf("Expected %d checks to be executed, got %d", len(expectedExecuted), len(executedChecks))
		}

		for i, expected := range expectedExecuted {
			if i >= len(executedChecks) || executedChecks[i] != expected {
				t.Errorf("Expected check %s at position %d, got %s", expected, i, executedChecks[i])
			}
		}

		if strings.Contains(strings.Join(executedChecks, ","), "check3") {
			t.Errorf("check3 should not have been executed after check2 failed")
		}
	})
}

func TestIntegration_RealWorldScenario(t *testing.T) {
	// Test based on the actual usage in the codebase
	t.Run("github pr validation scenario", func(t *testing.T) {
		validConfig := map[string]string{
			"title": "My PR Title",
			"body":  "My PR Body",
		}

		invalidConfig := map[string]string{
			"title": "My PR Title",
			// missing "body"
		}

		// Test valid configuration
		validator := FailOnFirst(
			RequireKey("title"),
			RequireKey("body"),
		)

		err := validator(validConfig)
		if err != nil {
			t.Errorf("Valid config should pass validation, got error: %v", err)
		}

		// Test invalid configuration
		err = validator(invalidConfig)
		if err == nil {
			t.Errorf("Invalid config should fail validation")
		}

		expectedError := "missing required config key 'body'"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
}
