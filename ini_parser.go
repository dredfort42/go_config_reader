/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | ini_parser.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"strconv"
	"strings"
)

// parseINI parses INI format content according to INI file structure.
// Supports sections, nested structure, multiple comment styles, quoted strings,
// escape sequences, multi-line values, and proper error handling.
func (c *Config) parseINI(content string) map[string]any {
	if c == nil {
		return nil
	}

	result := make(map[string]any)
	lines := strings.Split(content, "\n")

	var currentSection string

	var currentMap map[string]any = result

	for i, line := range lines {
		// Handle multi-line values (lines ending with backslash)
		for strings.HasSuffix(strings.TrimSpace(line), "\\") && i+1 < len(lines) {
			// Remove backslash and any trailing whitespace from current line
			line = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), "\\"))

			i++
			if i < len(lines) {
				nextLine := strings.TrimSpace(lines[i])
				if nextLine != "" {
					// Only add space if current line is not empty
					if line != "" {
						line += " " + nextLine
					} else {
						line = nextLine
					}
				}
			}
		}

		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Handle section headers [section_name]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimSpace(line[1 : len(line)-1])

			// Validate section name
			if sectionName == "" {
				// Empty section name - ignore but don't change context
				continue
			}

			// Check for invalid characters in section name
			if strings.ContainsAny(sectionName, "[]#;=") {
				// For invalid sections, invalidate current context so keys are ignored
				currentMap = nil

				continue
			}

			currentSection = sectionName

			// Create nested map for section if it doesn't exist
			if _, exists := result[currentSection]; !exists {
				result[currentSection] = make(map[string]any)
			}

			// Set current map to the section map
			if sectionMap, ok := result[currentSection].(map[string]any); ok {
				currentMap = sectionMap
			} else {
				// If section exists but is not a map, replace it with a map
				result[currentSection] = make(map[string]any)
				if m, ok := result[currentSection].(map[string]any); ok {
					currentMap = m
				}
			}

			continue
		}

		// Remove inline comments (but not if inside quotes)
		processedLine := c.removeInlineComments(line)
		if processedLine == "" {
			continue
		}

		// Parse key-value pairs
		key, value, found := strings.Cut(processedLine, "=")
		if !found {
			continue // Skip lines without = separator
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Validate key name
		if key == "" || strings.ContainsAny(key, "[]#;=") {
			continue // Skip invalid keys
		}

		// Process the value (handle quotes and escape sequences)
		processedValue := c.processINIValue(value)

		// Store in current map (either root or current section) if valid context
		if currentMap != nil {
			currentMap[key] = processedValue
		}
	}

	return result
}

// removeInlineComments removes inline comments while preserving quoted strings.
func (c *Config) removeInlineComments(line string) string {
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		char := line[i]

		// Handle escape sequences
		if char == '\\' && i+1 < len(line) {
			i++ // Skip next character

			continue
		}

		// Handle quotes
		if (char == '"' || char == '\'') && !inQuotes {
			inQuotes = true
			quoteChar = char

			continue
		} else if char == quoteChar && inQuotes {
			inQuotes = false
			quoteChar = 0

			continue
		}

		// Handle comments (only if not in quotes)
		if !inQuotes && (char == '#' || char == ';') {
			return strings.TrimSpace(line[:i])
		}
	}

	return line
}

// processINIValue processes INI values, handling quotes, escape sequences, and type conversion.
func (c *Config) processINIValue(value string) any {
	if value == "" {
		return ""
	}

	// Handle quoted strings
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			// Remove quotes and process escape sequences
			unquoted := value[1 : len(value)-1]

			return c.processEscapeSequences(unquoted)
		}
	}

	// Handle boolean values
	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "yes", "on", "1":
		return true
	case "false", "no", "off", "0":
		return false
	}

	// Try to parse as integer
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		// Return int for smaller values, int64 for larger ones
		if intVal >= int64(^uint(0)>>1) || intVal <= -int64(^uint(0)>>1)-1 {
			return intVal
		}

		return int(intVal)
	}

	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Handle comma-separated lists
	if strings.Contains(value, ",") {
		parts := strings.Split(value, ",")

		var result []string

		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, c.processEscapeSequences(trimmed))
			}
		}

		if len(result) > 1 {
			return result
		}
	}

	// Return as string with escape sequences processed
	return c.processEscapeSequences(value)
}

// processEscapeSequences processes escape sequences in INI values.
func (c *Config) processEscapeSequences(value string) string {
	if !strings.Contains(value, "\\") {
		return value
	}

	var result strings.Builder

	result.Grow(len(value))

	for i := 0; i < len(value); i++ {
		if value[i] == '\\' && i+1 < len(value) {
			switch value[i+1] {
			case 'n':
				result.WriteByte('\n')

				i++
			case 't':
				result.WriteByte('\t')

				i++
			case 'r':
				result.WriteByte('\r')

				i++
			case '\\':
				result.WriteByte('\\')

				i++
			case '"':
				result.WriteByte('"')

				i++
			case '\'':
				result.WriteByte('\'')

				i++
			case '0':
				result.WriteByte('\000')

				i++
			default:
				// Unknown escape sequence, keep as is
				result.WriteByte(value[i])
			}
		} else {
			result.WriteByte(value[i])
		}
	}

	return result.String()
}
