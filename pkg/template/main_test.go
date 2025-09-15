package template

import (
	"reflect"
	"testing"
)

func TestExec(t *testing.T) {
	type TestPayload struct {
		Name   string
		Age    int
		Active bool
		Items  []string
	}

	tests := []struct {
		name      string
		tplName   string
		raw       string
		payload   any
		want      *string
		wantError bool
	}{
		{
			name:      "simple text replacement",
			tplName:   "simple",
			raw:       "Hello {{.Name}}!",
			payload:   TestPayload{Name: "World"},
			want:      stringPtr("Hello World!"),
			wantError: false,
		},
		{
			name:      "multiple field access",
			tplName:   "multiple",
			raw:       "Name: {{.Name}}, Age: {{.Age}}, Active: {{.Active}}",
			payload:   TestPayload{Name: "John", Age: 30, Active: true},
			want:      stringPtr("Name: John, Age: 30, Active: true"),
			wantError: false,
		},
		{
			name:      "no template variables",
			tplName:   "static",
			raw:       "This is static text",
			payload:   TestPayload{},
			want:      stringPtr("This is static text"),
			wantError: false,
		},
		{
			name:      "empty template",
			tplName:   "empty",
			raw:       "",
			payload:   TestPayload{},
			want:      stringPtr(""),
			wantError: false,
		},
		{
			name:      "nil payload",
			tplName:   "nil_payload",
			raw:       "Static text only",
			payload:   nil,
			want:      stringPtr("Static text only"),
			wantError: false,
		},
		{
			name:      "range over slice",
			tplName:   "range",
			raw:       "Items: {{range .Items}}{{.}} {{end}}",
			payload:   TestPayload{Items: []string{"apple", "banana", "cherry"}},
			want:      stringPtr("Items: apple banana cherry "),
			wantError: false,
		},
		{
			name:      "conditional template",
			tplName:   "conditional",
			raw:       "{{if .Active}}User is active{{else}}User is inactive{{end}}",
			payload:   TestPayload{Active: true},
			want:      stringPtr("User is active"),
			wantError: false,
		},
		{
			name:      "conditional template - false condition",
			tplName:   "conditional_false",
			raw:       "{{if .Active}}User is active{{else}}User is inactive{{end}}",
			payload:   TestPayload{Active: false},
			want:      stringPtr("User is inactive"),
			wantError: false,
		},
		{
			name:      "join function with individual strings",
			tplName:   "join_strings",
			raw:       "{{join \"-\" \"hello\" \"world\" \"test\"}}",
			payload:   TestPayload{},
			want:      stringPtr("hello-world-test"),
			wantError: false,
		},
		{
			name:      "map payload",
			tplName:   "map_payload",
			raw:       "Key1: {{.key1}}, Key2: {{.key2}}",
			payload:   map[string]string{"key1": "value1", "key2": "value2"},
			want:      stringPtr("Key1: value1, Key2: value2"),
			wantError: false,
		},
		{
			name:    "nested struct access",
			tplName: "nested",
			raw:     "User: {{.User.Name}}, Email: {{.User.Email}}",
			payload: struct {
				User struct {
					Name  string
					Email string
				}
			}{
				User: struct {
					Name  string
					Email string
				}{Name: "John", Email: "john@example.com"},
			},
			want:      stringPtr("User: John, Email: john@example.com"),
			wantError: false,
		},
		{
			name:      "invalid template syntax",
			tplName:   "invalid",
			raw:       "Hello {{.Name",
			payload:   TestPayload{Name: "World"},
			want:      nil,
			wantError: true,
		},
		{
			name:      "accessing non-existent field",
			tplName:   "non_existent",
			raw:       "Hello {{.NonExistentField}}!",
			payload:   TestPayload{Name: "World"},
			want:      nil,
			wantError: true,
		},
		{
			name:      "whitespace handling",
			tplName:   "whitespace",
			raw:       "  {{.Name}}  ",
			payload:   TestPayload{Name: "Test"},
			want:      stringPtr("  Test  "),
			wantError: false,
		},
		{
			name:      "newlines in template",
			tplName:   "newlines",
			raw:       "Line 1: {{.Name}}\nLine 2: {{.Age}}",
			payload:   TestPayload{Name: "John", Age: 25},
			want:      stringPtr("Line 1: John\nLine 2: 25"),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec(tt.tplName, tt.raw, tt.payload)

			if tt.wantError {
				if err == nil {
					t.Errorf("Exec() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Exec() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExec_JoinFunction(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		payload   any
		want      string
		wantError bool
	}{
		{
			name:      "join with comma separator",
			template:  "{{join \", \" \"a\" \"b\" \"c\"}}",
			payload:   nil,
			want:      "a, b, c",
			wantError: false,
		},
		{
			name:      "join with dash separator",
			template:  "{{join \"-\" \"hello\" \"world\"}}",
			payload:   nil,
			want:      "hello-world",
			wantError: false,
		},
		{
			name:      "join with empty separator",
			template:  "{{join \"\" \"a\" \"b\" \"c\"}}",
			payload:   nil,
			want:      "abc",
			wantError: false,
		},
		{
			name:      "join with single string",
			template:  "{{join \", \" \"single\"}}",
			payload:   nil,
			want:      "single",
			wantError: false,
		},
		{
			name:      "join with no strings",
			template:  "{{join \", \"}}",
			payload:   nil,
			want:      "",
			wantError: false,
		},
		{
			name:      "join with space separator",
			template:  "{{join \" \" \"word1\" \"word2\" \"word3\"}}",
			payload:   nil,
			want:      "word1 word2 word3",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec("test", tt.template, tt.payload)

			if tt.wantError {
				if err == nil {
					t.Errorf("Exec() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Exec() unexpected error: %v", err)
				return
			}

			if got == nil {
				t.Errorf("Exec() returned nil result")
				return
			}

			if *got != tt.want {
				t.Errorf("Exec() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestExec_EdgeCases(t *testing.T) {
	t.Run("empty template name", func(t *testing.T) {
		got, err := Exec("", "Hello {{.Name}}", map[string]string{"Name": "World"})
		if err != nil {
			t.Errorf("Exec() with empty name unexpected error: %v", err)
		}
		if got == nil || *got != "Hello World" {
			t.Errorf("Exec() with empty name = %v, want %v", got, stringPtr("Hello World"))
		}
	})

	t.Run("very large template", func(t *testing.T) {
		// Create a large template
		largeTemplate := ""
		for i := 0; i < 1000; i++ {
			largeTemplate += "{{.Name}} "
		}

		got, err := Exec("large", largeTemplate, map[string]string{"Name": "Test"})
		if err != nil {
			t.Errorf("Exec() with large template unexpected error: %v", err)
		}

		if got == nil {
			t.Errorf("Exec() with large template returned nil")
		} else if len(*got) == 0 {
			t.Errorf("Exec() with large template returned empty string")
		}
	})

	t.Run("special characters in template", func(t *testing.T) {
		template := "Special chars: !@#$%^&*() {{.Name}} àáâãäå"
		payload := map[string]string{"Name": "Test"}
		want := "Special chars: !@#$%^&*() Test àáâãäå"

		got, err := Exec("special", template, payload)
		if err != nil {
			t.Errorf("Exec() with special chars unexpected error: %v", err)
		}

		if got == nil || *got != want {
			t.Errorf("Exec() with special chars = %v, want %v", got, stringPtr(want))
		}
	})
}

func TestExec_ReturnPointer(t *testing.T) {
	// Test that the function returns a pointer to the result string
	got, err := Exec("test", "Hello World", nil)
	if err != nil {
		t.Errorf("Exec() unexpected error: %v", err)
	}

	if got == nil {
		t.Errorf("Exec() returned nil pointer")
	}

	// Verify it's actually a pointer by modifying the original
	// and checking the result doesn't change
	original := *got
	*got = "Modified"

	// Call again to get a fresh result
	got2, err := Exec("test", "Hello World", nil)
	if err != nil {
		t.Errorf("Exec() second call unexpected error: %v", err)
	}

	if *got2 != "Hello World" {
		t.Errorf("Expected fresh result to be unchanged, got %v", *got2)
	}

	// Restore for cleanliness
	*got = original
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
