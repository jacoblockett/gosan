package gosan

import (
	"testing"
)

func TestSanitizeWindowsFilename(t *testing.T) {
	Simulation = "windows"
	defer func() {
		Simulation = ""
	}()
	tests := []struct {
		input    string
		expected string
	}{
		{"   ", ""},
		{"example", "example"},
		{"abc.txt", "abc.txt"},
		{"<>:\"/\\|?*abc.txt", "abc.txt"},
		{"abc\x1f.txt", "abc.txt"},
		{"abc\x7f.txt", "abc.txt"},
		{"NUL", ""},
		{"NUL.txt", ".txt"},
		{"NUL.tar.gz", ".tar.gz"},
	}

	for _, test := range tests {
		result := Filename(test.input, "")
		if result != test.expected {
			t.Errorf("Sanitize(\"%s\") = \"%s\"; expected \"%s\"", test.input, result, test.expected)
		}
	}
}

func TestSanitizeLinuxAndUnixFilename(t *testing.T) {
	Simulation = "not windows"
	defer func() {
		Simulation = ""
	}()
	tests := []struct {
		input    string
		expected string
	}{
		{"   ", "   "},
		{"example", "example"},
		{"abc.txt", "abc.txt"},
		{"abc\x1f.txt", "abc.txt"},
		{"abc\x7f.txt", "abc.txt"},
		{".", ""},
		{"..", ""},
		{"/..", ""},
	}

	for _, test := range tests {
		result := Filename(test.input, "")
		if result != test.expected {
			t.Errorf("Sanitize(\"%s\") = \"%s\"; expected \"%s\"", test.input, result, test.expected)
		}
	}
}
