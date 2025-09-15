package mealie

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchRecipeImage(t *testing.T) {
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		recipeId      string
		version       string
		want          []byte
		wantError     bool
		errorContains string
	}{
		{
			name: "successful image fetch",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					expectedPath := "/media/recipes/recipe123/images/original.webp"
					expectedVersion := "v1"

					if r.URL.Path != expectedPath {
						t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
					}

					if r.URL.Query().Get("version") != expectedVersion {
						t.Errorf("Expected version %s, got %s", expectedVersion, r.URL.Query().Get("version"))
					}

					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("fake-webp-content"))
				}))
			},
			recipeId:  "recipe123",
			version:   "v1",
			want:      []byte("fake-webp-content"),
			wantError: false,
		},
		{
			name: "json response returns nil",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"error": "image not found"}`))
				}))
			},
			recipeId:  "recipe456",
			version:   "v2",
			want:      nil,
			wantError: false,
		},
		{
			name: "server error response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("server error"))
				}))
			},
			recipeId:  "recipe789",
			version:   "v3",
			want:      []byte("server error"),
			wantError: false,
		},
		{
			name: "empty image response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusOK)
					// Empty response body
				}))
			},
			recipeId:  "empty-recipe",
			version:   "v1",
			want:      []byte(""),
			wantError: false,
		},
		{
			name: "large image response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusOK)

					// Generate a large response (1MB of 'A' characters)
					largeContent := bytes.Repeat([]byte("A"), 1024*1024)
					w.Write(largeContent)
				}))
			},
			recipeId:  "large-recipe",
			version:   "v1",
			want:      bytes.Repeat([]byte("A"), 1024*1024),
			wantError: false,
		},
		{
			name: "different content types",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/jpeg")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("fake-jpeg-content"))
				}))
			},
			recipeId:  "jpeg-recipe",
			version:   "v1",
			want:      []byte("fake-jpeg-content"),
			wantError: false,
		},
		{
			name: "special characters in recipe id",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					expectedPath := "/media/recipes/recipe-with-dashes_and_underscores/images/original.webp"
					if r.URL.Path != expectedPath {
						t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
					}

					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("special-id-content"))
				}))
			},
			recipeId:  "recipe-with-dashes_and_underscores",
			version:   "v1",
			want:      []byte("special-id-content"),
			wantError: false,
		},
		{
			name: "empty version parameter",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Query().Get("version") != "" {
						t.Errorf("Expected empty version, got %s", r.URL.Query().Get("version"))
					}

					w.Header().Set("Content-Type", "image/webp")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("empty-version-content"))
				}))
			},
			recipeId:  "recipe123",
			version:   "",
			want:      []byte("empty-version-content"),
			wantError: false,
		},
		{
			name: "json content-type case insensitive",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "Application/JSON")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"message": "not found"}`))
				}))
			},
			recipeId:  "case-test",
			version:   "v1",
			want:      []byte(`{"message": "not found"}`), // Should not return nil because case doesn't match exactly
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			got, err := FetchRecipeImage(server.URL, tt.recipeId, tt.version)

			if tt.wantError {
				if err == nil {
					t.Errorf("FetchRecipeImage() expected error but got none")
				}
				if tt.errorContains != "" && err != nil {
					if !bytes.Contains([]byte(err.Error()), []byte(tt.errorContains)) {
						t.Errorf("FetchRecipeImage() error = %v, want to contain %v", err, tt.errorContains)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FetchRecipeImage() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchRecipeImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchRecipeImage_NetworkErrors(t *testing.T) {
	tests := []struct {
		name      string
		apiUrl    string
		recipeId  string
		version   string
		wantError bool
	}{
		{
			name:      "invalid url",
			apiUrl:    "invalid-url",
			recipeId:  "recipe123",
			version:   "v1",
			wantError: true,
		},
		{
			name:      "unreachable server",
			apiUrl:    "http://localhost:99999",
			recipeId:  "recipe123",
			version:   "v1",
			wantError: true,
		},
		{
			name:      "malformed url with spaces",
			apiUrl:    "http://example .com",
			recipeId:  "recipe123",
			version:   "v1",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchRecipeImage(tt.apiUrl, tt.recipeId, tt.version)

			if tt.wantError {
				if err == nil {
					t.Errorf("FetchRecipeImage() expected error but got none")
				}
				if got != nil {
					t.Errorf("FetchRecipeImage() expected nil result on error, got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("FetchRecipeImage() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFetchRecipeImage_URLConstruction(t *testing.T) {
	// Test that URL is constructed correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the full URL path and query parameters
		expectedPath := "/media/recipes/test-recipe-id/images/original.webp"
		expectedVersion := "test-version"

		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.URL.Query().Get("version") != expectedVersion {
			t.Errorf("Expected version %s, got %s", expectedVersion, r.URL.Query().Get("version"))
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "image/webp")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("url-test-content"))
	}))
	defer server.Close()

	got, err := FetchRecipeImage(server.URL, "test-recipe-id", "test-version")
	if err != nil {
		t.Errorf("FetchRecipeImage() unexpected error: %v", err)
	}

	expected := []byte("url-test-content")
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("FetchRecipeImage() = %v, want %v", got, expected)
	}
}

func TestFetchRecipeImage_ContentTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		content     []byte
		expectNil   bool
	}{
		{
			name:        "exact json match",
			contentType: "application/json",
			content:     []byte(`{"error": "not found"}`),
			expectNil:   true,
		},
		{
			name:        "json with charset",
			contentType: "application/json; charset=utf-8",
			content:     []byte(`{"error": "not found"}`),
			expectNil:   false, // Should not match exactly
		},
		{
			name:        "text/json",
			contentType: "text/json",
			content:     []byte(`{"error": "not found"}`),
			expectNil:   false,
		},
		{
			name:        "empty content type",
			contentType: "",
			content:     []byte("some-content"),
			expectNil:   false,
		},
		{
			name:        "uppercase json",
			contentType: "APPLICATION/JSON",
			content:     []byte(`{"error": "not found"}`),
			expectNil:   false, // Case sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(tt.content)
			}))
			defer server.Close()

			got, err := FetchRecipeImage(server.URL, "test-recipe", "v1")
			if err != nil {
				t.Errorf("FetchRecipeImage() unexpected error: %v", err)
			}

			if tt.expectNil {
				if got != nil {
					t.Errorf("FetchRecipeImage() expected nil, got %v", got)
				}
			} else {
				if !reflect.DeepEqual(got, tt.content) {
					t.Errorf("FetchRecipeImage() = %v, want %v", got, tt.content)
				}
			}
		})
	}
}

func TestFetchRecipeImage_ResponseBodyReadError(t *testing.T) {
	// This test simulates a scenario where the response body reading fails
	// We can't easily simulate io.ReadAll errors with httptest, but we can test
	// the case where server closes connection prematurely

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.WriteHeader(http.StatusOK)

		// Write some content then close connection abruptly
		w.Write([]byte("partial"))

		// Force connection close
		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, err := hijacker.Hijack()
			if err == nil {
				conn.Close()
			}
		}
	}))
	defer server.Close()

	// This should still work as the content was written before connection closed
	got, err := FetchRecipeImage(server.URL, "test-recipe", "v1")
	if err != nil {
		// Connection errors are acceptable in this test
		t.Logf("Expected connection error occurred: %v", err)
		return
	}

	// If no error, verify we got the partial content
	expected := []byte("partial")
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("FetchRecipeImage() = %v, want %v", got, expected)
	}
}
