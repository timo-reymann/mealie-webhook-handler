package configuration // Replace with your actual package name

import (
	"os"

	"testing"
)

func TestParse(t *testing.T) {
	content, err := os.ReadFile("testdata/valid.toml")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	configuration, err := ParseConfiguration(content)
	if err != nil {
		t.Fatalf("Failed to parse configuration: %v", err)
	}

	if configuration == nil {
		t.Fatalf("Configuration should not be nil")
	}

	// Test for identifierA
	identifierA, exists := configuration.Webhooks["identifierA"]
	if !exists {
		t.Errorf("Webhook 'identifierA' should exist")
	}

	if identifierA.TemplatePath != "file/to/template.gotpl" {
		t.Errorf("Expected TemplatePath 'file/to/template.gotpl' for identifierA, got %s", identifierA.TemplatePath)
	}

	if identifierA.Output != "github_pr" {
		t.Errorf("Expected Output.Name 'github_pr' for identifierA, got %s", identifierA.Output)
	}

	if identifierA.OutputOptions["foo"] != "bar" {
		t.Errorf("Expected Output.Options['foo'] 'bar' for identifierA, got %s", identifierA.OutputOptions["foo"])
	}

	// Test for identifierB
	identifierB, exists := configuration.Webhooks["identifierB"]
	if !exists {
		t.Errorf("Webhook 'identifierB' should exist")
	}

	if identifierB.TemplatePath != "file/to/template.gotpl" {
		t.Errorf("Expected TemplatePath 'file/to/template.gotpl' for identifierB, got %s", identifierB.TemplatePath)
	}

	if identifierB.Output != "github_pr" {
		t.Errorf("Expected Output.Name 'github_pr' for identifierB, got %s", identifierB.Output)
	}

	if identifierB.OutputOptions["foo"] != "bar" {
		t.Errorf("Expected Output.Options['foo'] 'bar' for identifierB, got %s", identifierB.OutputOptions["foo"])
	}
}
