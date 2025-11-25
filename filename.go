package gosan

import (
	"runtime"
	"strings"
	"unicode"
)

var Simulation string

func sanitizeWindows(filename string, replacementStr string) string {
	sanitized := []rune{}
	replacement := []rune(replacementStr)

	// Cannot end with whitespace or "."
	// Note: We trim this BEFORE processing. If you want these to be replaced
	// instead of trimmed, this logic would need to move into the loop.
	filename = strings.TrimRightFunc(filename, func(r rune) bool {
		return unicode.IsSpace(r) || r == '.'
	})

	// Efficiency return - trimmed filename is empty
	if filename == "" {
		return ""
	}

	// Illegal filename characters, printable characters
	illRunes := map[rune]int{'<': 1, '>': 1, ':': 1, '"': 1, '/': 1, '\\': 1, '|': 1, '?': 1, '*': 1}

	for _, existingRune := range filename {
		if existingRune <= 31 || existingRune == 127 || !unicode.IsPrint(existingRune) || illRunes[existingRune] == 1 {
			// If a replacement was provided, append it
			if len(replacement) > 0 {
				sanitized = append(sanitized, replacement...)
			}
			continue
		}

		sanitized = append(sanitized, existingRune)
	}

	// Efficiency return - sanitized filename is empty
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
	parts := strings.Split(sanitizedStr, ".") // An attempt to consider reserved names with extensions.

	if len(parts) > 0 && illNames[strings.ToUpper(parts[0])] == 1 {
		// Reserved names are tricky. If the filename is "CON.txt",
		// replacing "CON" with "_" results in "_.txt", which is valid.
		// If no replacement is provided, it stays as standard behavior (removes the name).
		if replacementStr != "" {
			parts[0] = replacementStr
		} else {
			parts[0] = ""
		}
		sanitizedStr = strings.Join(parts, ".")
	}

	return sanitizedStr
}

func sanitizeLinuxAndUnix(filename string, replacementStr string) string {
	sanitized := []rune{}
	replacement := []rune(replacementStr)

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

	// Efficiency return - sanitized filename is empty
	if len(sanitized) == 0 {
		return ""
	}

	sanitizedStr := string(sanitized)

	// Reserved names
	if sanitizedStr == "." || sanitizedStr == ".." {
		return ""
	}

	return sanitizedStr
}

// Filename sanitizes the given string under the assumption it represents a filename.
// Sanitation is OS-agnostic based on runtime or simulation. The Simulation variable
// is exposed via gosan.Simulation.
func Filename(filename string, replacement string) string {
	// Check for simulation first
	if Simulation == "windows" {
		return sanitizeWindows(filename, replacement)
	}
	if Simulation != "" {
		return sanitizeLinuxAndUnix(filename, replacement)
	}

	// Then check for runtime OS
	if runtime.GOOS == "windows" {
		return sanitizeWindows(filename, replacement)
	}

	return sanitizeLinuxAndUnix(filename, replacement)

}
