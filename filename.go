package gosan

import (
	"fmt"
	"runtime"
	"strings"
	"unicode"
)

const (
	Windows = "windows"
	Linux   = "linux"
	Darwin  = "darwin"
)

var ErrUnsupportedGOOS = fmt.Errorf("%s is not supported", runtime.GOOS)

type FilenameOptions struct {
	// Which environment should sanitization consider. Default=runtime.GOOS
	Environment string
	// String to use instead of illegal character. Default=""
	Replacement string
	// Instead of using the Replacement, attempts to replace illegal
	// characters with a visually similar, legal character. Default=false
	ReplaceWithVisuallySimilarRunes bool
	// String to prepend to reserved names. Default='RESERVED_'
	ReservedPrefix string
}

func appendDefaultOpts(incoming *FilenameOptions) *FilenameOptions {
	opts := &FilenameOptions{
		Environment:                     runtime.GOOS,
		Replacement:                     "",
		ReplaceWithVisuallySimilarRunes: false,
		ReservedPrefix:                  "RESERVED_",
	}

	if incoming == nil {
		return opts
	}

	if incoming.Environment != "" {
		opts.Environment = incoming.Environment
	}

	if incoming.Replacement != "" {
		opts.Replacement = incoming.Replacement
	}

	if incoming.ReplaceWithVisuallySimilarRunes {
		opts.ReplaceWithVisuallySimilarRunes = true
	}

	if incoming.ReservedPrefix != "" {
		opts.ReservedPrefix = incoming.ReservedPrefix
	}

	return opts
}

func sanitizeWindows(filename string, opts *FilenameOptions) string {
	sanitized := []rune{}
	replacement := []rune(opts.Replacement)

	// Cannot end with whitespace or "."
	filename = strings.TrimRightFunc(filename, func(r rune) bool {
		return unicode.IsSpace(r) || r == '.'
	})

	if filename == "" {
		return ""
	}

	// Illegal filename characters, printable characters
	illRunes := map[rune]int{
		'<':  1,
		'>':  1,
		':':  1,
		'"':  1,
		'/':  1,
		'\\': 1,
		'|':  1,
		'?':  1,
		'*':  1,
	}
	illRunesSimilar := map[rune]rune{
		'<':  '˂', // U+02C2
		'>':  '˃', // U+02C3
		':':  '꞉', // U+A789
		'"':  '＂', // U+FF02
		'/':  '∕', // U+2215
		'\\': '∖', // U+2216
		'|':  'ǀ', // U+01C0
		'?':  '？', // U+FF1F
		'*':  '∗', // U+2217
	}

	for _, existingRune := range filename {
		if opts.ReplaceWithVisuallySimilarRunes && illRunes[existingRune] == 1 {
			sanitized = append(sanitized, illRunesSimilar[existingRune])
			continue
		}

		if existingRune <= 31 || existingRune == 127 || !unicode.IsPrint(existingRune) || illRunes[existingRune] == 1 {
			if len(replacement) > 0 {
				sanitized = append(sanitized, replacement...)
			}
			continue
		}

		sanitized = append(sanitized, existingRune)
	}

	if len(sanitized) == 0 {
		return ""
	}

	// Reserved names
	illNames := map[string]int{
		"CON": 1, "PRN": 1, "AUX": 1, "NUL": 1,
		"COM1": 1, "COM2": 1, "COM3": 1, "COM4": 1, "COM5": 1, "COM6": 1, "COM7": 1, "COM8": 1, "COM9": 1,
		"LPT1": 1, "LPT2": 1, "LPT3": 1, "LPT4": 1, "LPT5": 1, "LPT6": 1, "LPT7": 1, "LPT8": 1, "LPT9": 1,
	}
	sanitizedStr := string(sanitized)
	parts := strings.Split(sanitizedStr, ".") // An attempt to consider reserved names with extensions. While this isn't "illegal", things get weird, so best to avoid.

	if len(parts) > 0 && illNames[strings.ToUpper(parts[0])] == 1 {
		parts[0] = opts.ReservedPrefix + parts[0]
		sanitizedStr = strings.Join(parts, ".")
	}

	return sanitizedStr
}

func sanitizeLinuxAndUnix(filename string, opts *FilenameOptions) string {
	if opts.ReservedPrefix == "" {
		opts.ReservedPrefix = "RESERVED_"
	}

	sanitized := []rune{}
	replacement := []rune(opts.Replacement)

	// Illegal filename characters
	for _, existingRune := range filename {
		if existingRune <= 31 || existingRune == 127 || existingRune == '/' || !unicode.IsPrint(existingRune) {
			if len(replacement) > 0 {
				sanitized = append(sanitized, replacement...)
			}
			continue
		}

		sanitized = append(sanitized, existingRune)
	}

	if len(sanitized) == 0 {
		return ""
	}

	sanitizedStr := string(sanitized)

	// Reserved names
	if sanitizedStr == "." || sanitizedStr == ".." {
		sanitizedStr = opts.ReservedPrefix + sanitizedStr
	}

	return sanitizedStr
}

// Sanitizes the given string under the assumption it represents a filename.
// Sanitation is OS-agnostic based on runtime.
func Filename(filename string, opts *FilenameOptions) (string, error) {
	newOpts := appendDefaultOpts(opts)

	if newOpts.Environment == Windows {
		return sanitizeWindows(filename, newOpts), nil
	}

	if newOpts.Environment == Linux || newOpts.Environment == Darwin {
		return sanitizeLinuxAndUnix(filename, newOpts), nil
	}

	return "", ErrUnsupportedGOOS
}
