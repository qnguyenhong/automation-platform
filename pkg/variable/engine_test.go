package variable

import (
	"fmt"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantErr  bool
	}{
		{
			name:    "literal only",
			input:   "hello world",
			wantLen: 1,
		},
		{
			name:    "simple variable",
			input:   "Hello {{name}}",
			wantLen: 2,
		},
		{
			name:    "function call",
			input:   "{{$random_int(1,10)}}",
			wantLen: 1,
		},
		{
			name:    "mixed",
			input:   "User {{name}} has {{$random_int(1,100)}} items",
			wantLen: 5,
		},
		{
			name:    "uuid",
			input:   "{{$uuid}}",
			wantLen: 1,
		},
		{
			name:    "unclosed brace",
			input:   "{{name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTemplate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && len(got.Parts) != tt.wantLen {
				t.Errorf("ParseTemplate(%q) got %d parts, want %d", tt.input, len(got.Parts), tt.wantLen)
			}
		})
	}
}

func TestEngineResolve(t *testing.T) {
	engine := NewEngine(nil, map[string]string{
		"API_KEY": "test-key-123",
		"BASE_URL": "https://api.example.com",
	})

	tests := []struct {
		name     string
		input    string
		checkFn  func(string) bool
	}{
		{
			name:  "literal",
			input: "hello",
			checkFn: func(s string) bool { return s == "hello" },
		},
		{
			name:  "env variable",
			input: "Key: {{$env.API_KEY}}",
			checkFn: func(s string) bool { return s == "Key: test-key-123" },
		},
		{
			name:  "uuid generates valid format",
			input: "{{$uuid}}",
			checkFn: func(s string) bool { return len(s) == 36 && s[8] == '-' },
		},
		{
			name:  "random email",
			input: "{{$random_email}}",
			checkFn: func(s string) bool {
				for i, c := range s {
					if c == '@' && i > 0 {
						return true
					}
				}
				return false
			},
		},
		{
			name:  "random int in range",
			input: "{{$random_int(1,10)}}",
			checkFn: func(s string) bool {
				var n int
				_, err := fmt.Sscanf(s, "%d", &n)
				return err == nil && n >= 1 && n <= 10
			},
		},
		{
			name:  "random string length",
			input: "{{$random_string(16)}}",
			checkFn: func(s string) bool { return len(s) == 16 },
		},
		{
			name:  "timestamp",
			input: "{{$timestamp}}",
			checkFn: func(s string) bool { return len(s) == 10 },
		},
		{
			name:  "iso date",
			input: "{{$iso_date}}",
			checkFn: func(s string) bool { return len(s) > 10 && s[4] == '-' },
		},
		{
			name:  "random name",
			input: "{{$random_name}}",
			checkFn: func(s string) bool { return len(s) > 3 },
		},
		{
			name:  "captured variable",
			input: "User ID: {{user_id}}",
			checkFn: func(s string) bool { return s == "User ID: 42" },
		},
		{
			name:  "env base url",
			input: "{{$env.BASE_URL}}/users",
			checkFn: func(s string) bool { return s == "https://api.example.com/users" },
		},
	}

	// Set captured variable
	engine.Capture("user_id", "42")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Resolve(tt.input)
			if err != nil {
				t.Errorf("Resolve(%q) error = %v", tt.input, err)
				return
			}
			if !tt.checkFn(got) {
				t.Errorf("Resolve(%q) = %q, check failed", tt.input, got)
			}
		})
	}
}

func TestEngineResolveMap(t *testing.T) {
	engine := NewEngine(nil, map[string]string{
		"TOKEN": "abc123",
	})

	input := map[string]any{
		"url":     "{{$random_url}}",
		"headers": map[string]any{
			"Authorization": "Bearer {{$env.TOKEN}}",
		},
		"body": map[string]any{
			"name":  "{{$random_name}}",
			"email": "{{$random_email}}",
		},
	}

	result, err := engine.ResolveMap(input)
	if err != nil {
		t.Fatalf("ResolveMap error: %v", err)
	}

	// Check URL was resolved (no {{ remaining)
	url := result["url"].(string)
	if url == "{{$random_url}}" {
		t.Error("URL was not resolved")
	}

	// Check header was resolved
	headers := result["headers"].(map[string]any)
	auth := headers["Authorization"].(string)
	if auth != "Bearer abc123" {
		t.Errorf("Authorization = %q, want %q", auth, "Bearer abc123")
	}

	// Check body was resolved
	body := result["body"].(map[string]any)
	name := body["name"].(string)
	if name == "{{$random_name}}" {
		t.Error("Name was not resolved")
	}
}

func TestEngineSequence(t *testing.T) {
	engine := NewEngine(nil, nil)

	v1, _ := engine.Resolve("{{$seq(counter)}}")
	v2, _ := engine.Resolve("{{$seq(counter)}}")
	v3, _ := engine.Resolve("{{$seq(counter)}}")

	if v1 != "1" || v2 != "2" || v3 != "3" {
		t.Errorf("Sequence = %s, %s, %s, want 1, 2, 3", v1, v2, v3)
	}
}
