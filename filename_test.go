package gosan

import (
	"testing"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		input    string
		opts     *FilenameOptions
		expected string
	}{
		// nil opts test (will use runtime.GOOS)
		{"abc.txt", nil, "abc.txt"},
		// Windows
		{"   ", &FilenameOptions{Environment: Windows}, ""},
		{"example", &FilenameOptions{Environment: Windows}, "example"},
		{"abc.txt", &FilenameOptions{Environment: Windows}, "abc.txt"},
		{"<>:\"/\\|?*abc.txt", &FilenameOptions{Environment: Windows}, "abc.txt"},
		{"abc\x1f.txt", &FilenameOptions{Environment: Windows}, "abc.txt"},
		{"abc\x7f.txt", &FilenameOptions{Environment: Windows}, "abc.txt"},
		{"abc\x7f.txt", &FilenameOptions{Environment: Windows, Replacement: "_x_"}, "abc_x_.txt"},
		{"NUL", &FilenameOptions{Environment: Windows}, "RESERVED_NUL"},
		{"nul", &FilenameOptions{Environment: Windows}, "RESERVED_nul"},
		{"NUL.txt", &FilenameOptions{Environment: Windows}, "RESERVED_NUL.txt"},
		{"NUL.tar.gz", &FilenameOptions{Environment: Windows}, "RESERVED_NUL.tar.gz"},
		{"file. ", &FilenameOptions{Environment: Windows}, "file"},
		{"file .", &FilenameOptions{Environment: Windows}, "file"},
		{".  . .. .", &FilenameOptions{Environment: Windows}, ""},
		{"NUL", &FilenameOptions{Environment: Windows, ReservedPrefix: "r_"}, "r_NUL"},
		{"<>:\"/\\|?*abc.txt", &FilenameOptions{Environment: Windows, ReplaceWithVisuallySimilarRunes: true}, "˂˃꞉＂∕∖ǀ？∗abc.txt"},
		// Linux/Darwin
		{"   ", &FilenameOptions{Environment: Linux}, "   "},
		{"example", &FilenameOptions{Environment: Linux}, "example"},
		{"abc.txt", &FilenameOptions{Environment: Linux}, "abc.txt"},
		{"abc\x1f.txt", &FilenameOptions{Environment: Linux}, "abc.txt"},
		{"abc\x7f.txt", &FilenameOptions{Environment: Linux}, "abc.txt"},
		{"abc\x7f.txt", &FilenameOptions{Environment: Linux, Replacement: "_x_"}, "abc_x_.txt"},
		{"/..", &FilenameOptions{Environment: Linux, Replacement: "x"}, "x.."},
		{".", &FilenameOptions{Environment: Linux}, "RESERVED_."},
		{"..", &FilenameOptions{Environment: Linux}, "RESERVED_.."},
		{"/..", &FilenameOptions{Environment: Linux}, "RESERVED_.."},
		{".", &FilenameOptions{Environment: Linux, ReservedPrefix: "r_"}, "r_."},
	}

	for _, test := range tests {
		result, _ := Filename(test.input, test.opts)
		if result != test.expected {
			t.Errorf("[%s] Sanitize(\"%s\") = \"%s\"; expected \"%s\"", test.opts.Environment, test.input, result, test.expected)
		}
	}
}
