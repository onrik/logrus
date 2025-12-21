package gorm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected interface{}
	}{
		{
			name:     "Empty input",
			input:    []interface{}{},
			expected: nil,
		},
		{
			name: "Formatted SQL",
			input: []interface{}{
				"sql",
				nil,
				"Test API",
				"SELECT * FROM `test_data` WHERE `id` = $1",
				[]interface{}{10},
			},
			expected: "SELECT * FROM `test_data` WHERE `id` = '10' [Test API]",
		},
		{
			name: "Formatted SQL with multiple placeholders",
			input: []interface{}{
				"sql",
				nil,
				"Test API",
				"SELECT * FROM `test_data` WHERE `id` = $1 AND `name` = $2",
				[]interface{}{42, "John"},
			},
			expected: "SELECT * FROM `test_data` WHERE `id` = '42' AND `name` = 'John' [Test API]",
		},
		{
			name: "Formatted SQL with NULL value",
			input: []interface{}{
				"sql",
				nil,
				"Test API",
				"SELECT * FROM `test_data` WHERE `id` = $1",
				[]interface{}{nil},
			},
			expected: "SELECT * FROM `test_data` WHERE `id` = NULL [Test API]",
		},
		{
			name: "Non-SQL case",
			input: []interface{}{
				"other",
				nil,
				"Test API",
				"This is a message",
				[]interface{}{},
			},
			expected: "Test API",
		},
		{
			name: "Invalid input (less than 5 values)",
			input: []interface{}{
				"sql",
				nil,
				"Test API",
				"SELECT * FROM table WHERE id = $1",
			},
			expected: string("Test API"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Formatter(test.input...)
			require.Equal(t, test.expected, actual)
		})
	}
}
func TestFormatSource(t *testing.T) {
	tests := []struct {
		source         string
		expectedSource string
	}{
		{
			source:         "github.com/user1/repo",
			expectedSource: "user1/repo",
		},
		{
			source:         "github.com/user2/repo",
			expectedSource: "user2/repo",
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			result := formatSource(test.source)
			require.Equal(t, test.expectedSource, result)
		})
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "Hello, World!",
			expected: true,
		},
		{
			input:    "Some non-printable \x1b[31mtext\x1b[0m",
			expected: false,
		},
		{
			input:    "123",
			expected: true,
		},
		{
			input:    " ",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := isPrintable(test.input)
			require.Equal(t, test.expected, result)
		})
	}
}
