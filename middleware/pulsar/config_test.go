package pulsar

import (
	"testing"

	"github.com/romberli/go-util/constant"
)

func TestConfig_getURLString(t *testing.T) {
	tests := []struct {
		name        string
		urls        []string
		expected    string
		expectError bool
	}{
		{
			name:        "single url without prefix",
			urls:        []string{"192.168.1.1:6650"},
			expected:    PulsarSchemePrefix + "192.168.1.1:6650",
			expectError: false,
		},
		{
			name:        "single url with pulsar prefix",
			urls:        []string{"pulsar://192.168.1.1:6650"},
			expected:    "pulsar://192.168.1.1:6650",
			expectError: false,
		},
		{
			name:        "single url with pulsar+ssl prefix",
			urls:        []string{"pulsar+ssl://192.168.1.1:6650"},
			expected:    "pulsar+ssl://192.168.1.1:6650",
			expectError: false,
		},
		{
			name:        "multiple urls without prefix",
			urls:        []string{"192.168.1.1:6650", "192.168.1.2:6650", "192.168.1.3:6650"},
			expected:    PulsarSchemePrefix + "192.168.1.1:6650" + constant.CommaString + PulsarSchemePrefix + "192.168.1.2:6650" + constant.CommaString + PulsarSchemePrefix + "192.168.1.3:6650",
			expectError: false,
		},
		{
			name:        "urls with whitespace",
			urls:        []string{"  192.168.1.1:6650  ", "  192.168.1.2:6650  "},
			expected:    PulsarSchemePrefix + "192.168.1.1:6650" + constant.CommaString + PulsarSchemePrefix + "192.168.1.2:6650",
			expectError: false,
		},
		{
			name:        "urls with empty strings skipped",
			urls:        []string{"192.168.1.1:6650", "", "  ", "192.168.1.2:6650"},
			expected:    PulsarSchemePrefix + "192.168.1.1:6650" + constant.CommaString + PulsarSchemePrefix + "192.168.1.2:6650",
			expectError: false,
		},
		{
			name:        "all empty urls",
			urls:        []string{"", "  ", ""},
			expected:    constant.EmptyString,
			expectError: true,
		},
		{
			name:        "nil urls",
			urls:        nil,
			expected:    constant.EmptyString,
			expectError: true,
		},
		{
			name:        "empty urls slice",
			urls:        []string{},
			expected:    constant.EmptyString,
			expectError: true,
		},
		{
			name:        "mixed prefix urls",
			urls:        []string{"pulsar://192.168.1.1:6650", "192.168.1.2:6650", "pulsar+ssl://192.168.1.3:6650"},
			expected:    "pulsar://192.168.1.1:6650" + constant.CommaString + "pulsar://192.168.1.2:6650" + constant.CommaString + "pulsar+ssl://192.168.1.3:6650",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.urls, constant.EmptyString, constant.EmptyString)
			result, err := config.getURLString()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Fatalf("result mismatch: expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestConfig_Clone(t *testing.T) {
	original := NewConfig([]string{"192.168.1.1:6650", "192.168.1.2:6650"}, "token", "topic")
	cloned := original.Clone()

	if original == cloned {
		t.Fatal("clone should not return same pointer")
	}

	if len(original.URLs) != len(cloned.URLs) {
		t.Fatalf("urls length mismatch: original %d, cloned %d", len(original.URLs), len(cloned.URLs))
	}

	for i := range original.URLs {
		if original.URLs[i] != cloned.URLs[i] {
			t.Fatalf("url mismatch at index %d: original %s, cloned %s", i, original.URLs[i], cloned.URLs[i])
		}
	}

	if original.Token != cloned.Token {
		t.Fatalf("token mismatch: original %s, cloned %s", original.Token, cloned.Token)
	}

	if original.Topic != cloned.Topic {
		t.Fatalf("topic mismatch: original %s, cloned %s", original.Topic, cloned.Topic)
	}

	original.URLs[0] = "modified"
	if cloned.URLs[0] == "modified" {
		t.Fatal("clone should have independent urls slice")
	}
}