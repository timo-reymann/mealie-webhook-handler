package configuration

import (
	"context"
	"reflect"
	"testing"

	"go.deepl.dev/mealie-webhook-handler/pkg/output"
)

func TestParseConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		config    []byte
		want      *Config
		wantError bool
	}{
		{
			name: "valid configuration",
			config: []byte(`
[webhook.test_webhook]
template_path = "/path/to/template.gotpl"
output = "github_pr"

[webhook.test_webhook.output_options]
title = "Test PR"
body = "Test body"

[mealie]
api_url = "https://mealie.example.com"
`),
			want: &Config{
				Webhooks: map[string]WebhookConfig{
					"test_webhook": {
						TemplatePath: "/path/to/template.gotpl",
						Output:       "github_pr",
						OutputOptions: map[string]string{
							"title": "Test PR",
							"body":  "Test body",
						},
					},
				},
				Mealie: Mealie{
					ApiUrl: "https://mealie.example.com",
				},
			},
			wantError: false,
		},
		{
			name: "multiple webhooks",
			config: []byte(`
[webhook.webhook_one]
template_path = "/path/to/template1.gotpl"
output = "github_pr"

[webhook.webhook_one.output_options]
title = "First PR"

[webhook.webhook_two]
template_path = "/path/to/template2.gotpl"
output = "github_pr"

[webhook.webhook_two.output_options]
title = "Second PR"

[mealie]
api_url = "https://mealie.example.com"
`),
			want: &Config{
				Webhooks: map[string]WebhookConfig{
					"webhook_one": {
						TemplatePath: "/path/to/template1.gotpl",
						Output:       "github_pr",
						OutputOptions: map[string]string{
							"title": "First PR",
						},
					},
					"webhook_two": {
						TemplatePath: "/path/to/template2.gotpl",
						Output:       "github_pr",
						OutputOptions: map[string]string{
							"title": "Second PR",
						},
					},
				},
				Mealie: Mealie{
					ApiUrl: "https://mealie.example.com",
				},
			},
			wantError: false,
		},

		{
			name:      "empty configuration",
			config:    []byte(``),
			want:      &Config{},
			wantError: false,
		},
		{
			name:      "invalid TOML syntax",
			config:    []byte(`[invalid toml syntax`),
			want:      nil,
			wantError: true,
		},
		{
			name: "invalid TOML structure",
			config: []byte(`
[webhook]
invalid = "structure"
`),
			want:      nil,
			wantError: true,
		},
		{
			name: "complex output options",
			config: []byte(`
[webhook.complex]
template_path = "/path/to/template.gotpl"
output = "github_pr"

[webhook.complex.output_options]
title = "Complex PR"
body = "Multi-line\nbody content"
source_branch = "feature/new-recipe"
target_branch = "main"
repo_slug = "owner/repo"
`),
			want: &Config{
				Webhooks: map[string]WebhookConfig{
					"complex": {
						TemplatePath: "/path/to/template.gotpl",
						Output:       "github_pr",
						OutputOptions: map[string]string{
							"title":         "Complex PR",
							"body":          "Multi-line\nbody content",
							"source_branch": "feature/new-recipe",
							"target_branch": "main",
							"repo_slug":     "owner/repo",
						},
					},
				},
				Mealie: Mealie{},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfiguration(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseConfiguration() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseConfiguration() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfiguration() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseConfiguration_NilInput(t *testing.T) {
	got, err := ParseConfiguration(nil)
	if err != nil {
		t.Errorf("ParseConfiguration(nil) unexpected error: %v", err)
	}

	expected := &Config{}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ParseConfiguration(nil) = %+v, want %+v", got, expected)
	}
}

// Mock output for testing
type mockOutput struct {
	name             string
	initError        error
	validateError    error
	initCalled       bool
	validateCalled   bool
	validatedOptions map[string]string
}

func (m *mockOutput) Name() string {
	return m.name
}

func (m *mockOutput) Init() error {
	m.initCalled = true
	return m.initError
}

func (m *mockOutput) ValidateOptions(options map[string]string) error {
	m.validateCalled = true
	m.validatedOptions = options
	return m.validateError
}

func (m *mockOutput) Output(ctx context.Context, templatedRecipe string, image []byte, config map[string]string) error {
	return nil
}

func TestConfig_Init(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		setupMock func() *mockOutput
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty config - should pass",
			config: &Config{
				Webhooks: map[string]WebhookConfig{},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{name: "test_output"}
			},
			wantError: false,
		},
		{
			name: "valid single webhook - should pass",
			config: &Config{
				Webhooks: map[string]WebhookConfig{
					"test_webhook": {
						Output: "test_output",
						OutputOptions: map[string]string{
							"title": "Test Title",
						},
					},
				},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{
					name:          "test_output",
					initError:     nil,
					validateError: nil,
				}
			},
			wantError: false,
		},
		{
			name: "webhook with unknown output - should fail",
			config: &Config{
				Webhooks: map[string]WebhookConfig{
					"test_webhook": {
						Output: "unknown_output",
					},
				},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{name: "test_output"}
			},
			wantError: true,
			errorMsg:  "webhook test_webhook has no valid output",
		},
		{
			name: "output init fails - should fail",
			config: &Config{
				Webhooks: map[string]WebhookConfig{
					"test_webhook": {
						Output: "test_output",
					},
				},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{
					name:      "test_output",
					initError: &testError{msg: "init failed"},
				}
			},
			wantError: true,
			errorMsg:  "webhook test_webhook could not be initialized as output init faile with init failed",
		},
		{
			name: "output validation fails - should fail",
			config: &Config{
				Webhooks: map[string]WebhookConfig{
					"test_webhook": {
						Output: "test_output",
						OutputOptions: map[string]string{
							"invalid": "option",
						},
					},
				},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{
					name:          "test_output",
					validateError: &testError{msg: "validation failed"},
				}
			},
			wantError: true,
			errorMsg:  "webhook test_webhook has invalid output options validation failed",
		},
		{
			name: "multiple valid webhooks - should pass",
			config: &Config{
				Webhooks: map[string]WebhookConfig{
					"webhook_one": {
						Output: "test_output",
						OutputOptions: map[string]string{
							"title": "First",
						},
					},
					"webhook_two": {
						Output: "test_output",
						OutputOptions: map[string]string{
							"title": "Second",
						},
					},
				},
			},
			setupMock: func() *mockOutput {
				return &mockOutput{
					name: "test_output",
				}
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock output
			mock := tt.setupMock()
			output.Register(mock)

			// Run the test
			err := tt.config.Init()

			if tt.wantError {
				if err == nil {
					t.Errorf("Config.Init() expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Config.Init() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Config.Init() unexpected error: %v", err)
				}
			}

			// Verify mock interactions for successful cases
			if !tt.wantError && len(tt.config.Webhooks) > 0 {
				if !mock.initCalled {
					t.Errorf("Expected Init() to be called on mock output")
				}
				if !mock.validateCalled {
					t.Errorf("Expected ValidateOptions() to be called on mock output")
				}
			}
		})
	}
}

func TestConfig_Init_CallsCorrectMethods(t *testing.T) {

	config := &Config{
		Webhooks: map[string]WebhookConfig{
			"test_webhook": {
				Output: "test_output",
				OutputOptions: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}

	err := config.Init()
	if err != nil {
		t.Errorf("Config.Init() unexpected error: %v", err)
	}

}

// Helper error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
