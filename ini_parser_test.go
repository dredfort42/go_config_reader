/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | ini_parser_test.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestINIParser_Comprehensive tests INI parsing with extensive coverage
func TestINIParser_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	t.Run("BasicKeyValuePairs", func(t *testing.T) {
		iniContent := `
key1=value1
key2=value2
key_with_underscore=value_with_underscore
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "value1", result["key1"])
		assert.Equal(t, "value2", result["key2"])
		assert.Equal(t, "value_with_underscore", result["key_with_underscore"])
	})

	t.Run("SectionsAndNesting", func(t *testing.T) {
		iniContent := `
global_key=global_value

[section1]
key1=section1_value1
key2=section1_value2

[section2]
key1=section2_value1
nested_key=nested_value
`
		result := c.parseINI(iniContent)

		// Global key
		assert.Equal(t, "global_value", result["global_key"])

		// Section1 values
		section1, ok := result["section1"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "section1_value1", section1["key1"])
		assert.Equal(t, "section1_value2", section1["key2"])

		// Section2 values
		section2, ok := result["section2"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "section2_value1", section2["key1"])
		assert.Equal(t, "nested_value", section2["nested_key"])
	})

	t.Run("CommentsAndEmptyLines", func(t *testing.T) {
		iniContent := `
# This is a comment
; This is also a comment
key1=value1

# Empty lines above and below

key2=value2 # Inline comment
key3=value3 ; Another inline comment
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "value1", result["key1"])
		assert.Equal(t, "value2", result["key2"])
		assert.Equal(t, "value3", result["key3"])
	})

	t.Run("QuotedStrings", func(t *testing.T) {
		iniContent := `
double_quoted="This is a quoted string"
single_quoted='This is also quoted'
quoted_with_spaces="  spaced value  "
empty_quotes=""
single_empty_quotes=''
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "This is a quoted string", result["double_quoted"])
		assert.Equal(t, "This is also quoted", result["single_quoted"])
		assert.Equal(t, "  spaced value  ", result["quoted_with_spaces"])
		assert.Equal(t, "", result["empty_quotes"])
		assert.Equal(t, "", result["single_empty_quotes"])
	})

	t.Run("EscapeSequences", func(t *testing.T) {
		iniContent := `
newline="Line1\nLine2"
tab="Column1\tColumn2"
carriage_return="Line1\rLine2"
backslash="Path\\to\\file"
quote="Say \"Hello\""
single_quote="Don\'t worry"
null_char="Null\0char"
unknown_escape="Unknown\xsequence"
no_escape=NoEscape
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "Line1\nLine2", result["newline"])
		assert.Equal(t, "Column1\tColumn2", result["tab"])
		assert.Equal(t, "Line1\rLine2", result["carriage_return"])
		assert.Equal(t, "Path\\to\\file", result["backslash"])
		assert.Equal(t, "Say \"Hello\"", result["quote"])
		assert.Equal(t, "Don't worry", result["single_quote"])
		assert.Equal(t, "Null\000char", result["null_char"])
		assert.Equal(t, "Unknown\\xsequence", result["unknown_escape"]) // Unknown escape kept as-is
		assert.Equal(t, "NoEscape", result["no_escape"])
	})

	t.Run("TypeConversion", func(t *testing.T) {
		iniContent := `
# Boolean values
bool_true=true
bool_false=false
bool_yes=yes
bool_no=no
bool_on=on
bool_off=off
bool_1=1
bool_0=0

# Integer values
int_positive=42
int_negative=-123
int_zero=0
int_large=9223372036854775807

# Float values
float_positive=3.14159
float_negative=-2.71828
float_zero=0.0
float_scientific=1.23e-4

# String values (should remain strings)
string_value=hello_world
mixed_value=123abc
`
		result := c.parseINI(iniContent)

		// Boolean values
		assert.Equal(t, true, result["bool_true"])
		assert.Equal(t, false, result["bool_false"])
		assert.Equal(t, true, result["bool_yes"])
		assert.Equal(t, false, result["bool_no"])
		assert.Equal(t, true, result["bool_on"])
		assert.Equal(t, false, result["bool_off"])
		assert.Equal(t, true, result["bool_1"])
		assert.Equal(t, false, result["bool_0"])

		// Integer values
		assert.Equal(t, 42, result["int_positive"])
		assert.Equal(t, -123, result["int_negative"])
		assert.Equal(t, false, result["int_zero"]) // "0" is parsed as boolean false
		assert.Equal(t, int64(9223372036854775807), result["int_large"])

		// Float values
		assert.Equal(t, 3.14159, result["float_positive"])
		assert.Equal(t, -2.71828, result["float_negative"])
		assert.Equal(t, 0.0, result["float_zero"])
		assert.Equal(t, 1.23e-4, result["float_scientific"])

		// String values
		assert.Equal(t, "hello_world", result["string_value"])
		assert.Equal(t, "123abc", result["mixed_value"])
	})

	t.Run("ArrayValues", func(t *testing.T) {
		iniContent := `
single_item=item1
comma_separated=item1,item2,item3
spaced_items=item1, item2 , item3
empty_items=item1,,item3
trailing_comma=item1,item2,
leading_comma=,item2,item3
`
		result := c.parseINI(iniContent)

		// Single item should remain as string
		assert.Equal(t, "item1", result["single_item"])

		// Multiple items should become array
		assert.Equal(t, []string{"item1", "item2", "item3"}, result["comma_separated"])
		assert.Equal(t, []string{"item1", "item2", "item3"}, result["spaced_items"])
		assert.Equal(t, []string{"item1", "item3"}, result["empty_items"]) // Empty items filtered out
		assert.Equal(t, []string{"item1", "item2"}, result["trailing_comma"])
		assert.Equal(t, []string{"item2", "item3"}, result["leading_comma"])
	})

	t.Run("MultilineValues", func(t *testing.T) {
		iniContent := `
multiline_value=Line 1 \
    Line 2 \
    Line 3
single_line=No continuation
empty_continuation=Value \

continuation_at_end=Start \
    End
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "Line 1 Line 2 Line 3", result["multiline_value"])
		assert.Equal(t, "No continuation", result["single_line"])
		assert.Equal(t, "Value", result["empty_continuation"])
		assert.Equal(t, "Start End", result["continuation_at_end"])
	})

	t.Run("EdgeCasesAndInvalidData", func(t *testing.T) {
		iniContent := `
# Keys in root context (before any sections)
key_without_value=
key_with_empty_value=""
valid_key_final=final_value

# Empty section name should be ignored
[]
key_after_empty_section=value1

# Invalid section names should be ignored
[section#with#invalid]
invalid_key1=value1

[section;with;semicolon]
invalid_key2=value2

[section=with=equals]
invalid_key3=value3

# Valid section after invalid ones
[valid_section]
valid_key=valid_value

# Invalid key names should be ignored
key#with#hash=ignored
key;with;semicolon=ignored
=value_without_key
`
		result := c.parseINI(iniContent)

		// Valid keys should be parsed
		assert.Equal(t, "value1", result["key_after_empty_section"])
		assert.Equal(t, "final_value", result["valid_key_final"])

		// Valid section should be parsed
		validSection, ok := result["valid_section"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "valid_value", validSection["valid_key"])

		// Empty values should be handled at root level
		assert.Equal(t, "", result["key_without_value"])
		assert.Equal(t, "", result["key_with_empty_value"])

		// Invalid sections should not exist
		assert.NotContains(t, result, "section#with#invalid")
		assert.NotContains(t, result, "section;with;semicolon")
		assert.NotContains(t, result, "section=with=equals")

		// Invalid keys should not exist
		assert.NotContains(t, result, "key#with#hash")
		assert.NotContains(t, result, "key;with;semicolon")
		// Note: key=with=equals=ignored creates a key "key" with value "with=equals=ignored"
	})

	t.Run("InlineCommentsInQuotes", func(t *testing.T) {
		iniContent := `
quoted_with_hash="Value with # hash"
quoted_with_semicolon='Value with ; semicolon'
unquoted_with_comment=Value # This is a comment
mixed_quotes="Start with hash # and 'single quotes'"
escaped_quotes="Escaped \" quote # not a comment"
`
		result := c.parseINI(iniContent)
		assert.Equal(t, "Value with # hash", result["quoted_with_hash"])
		assert.Equal(t, "Value with ; semicolon", result["quoted_with_semicolon"])
		assert.Equal(t, "Value", result["unquoted_with_comment"])
		assert.Equal(t, "Start with hash # and 'single quotes'", result["mixed_quotes"])
		assert.Equal(t, "Escaped \" quote # not a comment", result["escaped_quotes"])
	})

	t.Run("SectionOverwriting", func(t *testing.T) {
		iniContent := `
[section1]
key1=first_value
key2=another_value

[section1]
key1=overwritten_value
key3=new_value
`
		result := c.parseINI(iniContent)

		section1, ok := result["section1"].(map[string]any)
		require.True(t, ok)

		// In the current implementation, sections are merged, not replaced
		// key1 should be overwritten
		assert.Equal(t, "overwritten_value", section1["key1"])

		// key2 should still be present (sections merge)
		assert.Equal(t, "another_value", section1["key2"])

		// key3 should be present
		assert.Equal(t, "new_value", section1["key3"])
	})

	t.Run("WhitespaceHandling", func(t *testing.T) {
		iniContent := `
   spaced_key   =   spaced_value   
	tabbed_key	=	tabbed_value	
mixed_whitespace = 	value 	
empty_line_with_spaces=   

[  spaced_section  ]
  section_key  =  section_value  
`
		result := c.parseINI(iniContent)

		assert.Equal(t, "spaced_value", result["spaced_key"])
		assert.Equal(t, "tabbed_value", result["tabbed_key"])
		assert.Equal(t, "value", result["mixed_whitespace"])
		assert.Equal(t, "", result["empty_line_with_spaces"])

		spacedSection, ok := result["spaced_section"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "section_value", spacedSection["section_key"])
	})
}

// TestProcessEscapeSequences_Comprehensive tests escape sequence processing
func TestProcessEscapeSequences_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"NoEscapes", "simple string", "simple string"},
		{"NewlineEscape", "line1\\nline2", "line1\nline2"},
		{"TabEscape", "col1\\tcol2", "col1\tcol2"},
		{"CarriageReturnEscape", "line1\\rline2", "line1\rline2"},
		{"BackslashEscape", "path\\\\to\\\\file", "path\\to\\file"},
		{"DoubleQuoteEscape", "say \\\"hello\\\"", "say \"hello\""},
		{"SingleQuoteEscape", "don\\'t", "don't"},
		{"NullCharEscape", "null\\0char", "null\000char"},
		{"MultipleEscapes", "\\n\\t\\r\\\\", "\n\t\r\\"},
		{"UnknownEscape", "unknown\\xescape", "unknown\\xescape"},
		{"EscapeAtEnd", "value\\", "value\\"},
		{"EmptyString", "", ""},
		{"OnlyBackslash", "\\", "\\"},
		{"MixedContent", "Start\\nMiddle\\tEnd", "Start\nMiddle\tEnd"},
		{"BackslashWithoutEscape", "normal\\text", "normal\text"}, // \t is recognized as tab
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := c.processEscapeSequences(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestRemoveInlineComments_Comprehensive tests inline comment removal
func TestRemoveInlineComments_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"NoComments", "key=value", "key=value"},
		{"HashComment", "key=value # comment", "key=value"},
		{"SemicolonComment", "key=value ; comment", "key=value"},
		{"QuotedHash", "key=\"value # not comment\"", "key=\"value # not comment\""},
		{"QuotedSemicolon", "key='value ; not comment'", "key='value ; not comment'"},
		{"EscapedQuote", "key=\"value \\\" # comment\"", "key=\"value \\\" # comment\""},
		{"MixedQuotes", "key=\"value 'with' quotes # comment\"", "key=\"value 'with' quotes # comment\""},
		{"EmptyAfterComment", "key=value #", "key=value"},
		{"OnlyComment", "# just a comment", ""},
		{"CommentAtStart", "# comment", ""},
		{"UnmatchedQuote", "key=\"unmatched quote # comment", "key=\"unmatched quote # comment"},
		{"EscapeInQuotes", "key=\"value\\nwith\\tescape # comment\"", "key=\"value\\nwith\\tescape # comment\""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := c.removeInlineComments(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestINIParser_EdgeCases tests edge cases for INI parser
func TestINIParser_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with nil config for parseINI
	var nilConfig *Config
	result := nilConfig.parseINI("key=value")
	assert.Nil(t, result)

	// Test with empty content
	result = c.parseINI("")
	assert.Empty(t, result)

	// Test with only whitespace
	result = c.parseINI("   \n\t\r\n   ")
	assert.Empty(t, result)

	// Test with only comments
	result = c.parseINI("# comment1\n; comment2\n")
	assert.Empty(t, result)
}

// Benchmark Tests for ini_parser.go functions

func BenchmarkConfig_ParseINI(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	iniContent := `
# Configuration file
[server]
host = localhost
port = 8080
ssl = true

[database]
host = db.example.com
port = 5432
name = myapp
ssl = false

[features]
auth = true
api = true
web = false
list = item1,item2,item3

# Global settings
timeout = 30s
debug = false
max_connections = 100
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.parseINI(iniContent)
	}
}

func BenchmarkConfig_RemoveInlineComments(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	testLines := []string{
		"key=value # inline comment",
		"quoted=\"value with # inside quotes\"",
		"no_comment=just_value",
		"escaped=\"value with \\\" quote # comment\"",
		"semicolon=value ; another comment",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range testLines {
			c.removeInlineComments(line)
		}
	}
}

func BenchmarkConfig_ProcessINIValue(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	testValues := []string{
		"42",
		"3.14159",
		"true",
		"false",
		"simple_string",
		"\"quoted string\"",
		"item1,item2,item3",
		"5m30s",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, value := range testValues {
			c.processINIValue(value)
		}
	}
}

func BenchmarkConfig_ProcessEscapeSequences(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	testValues := []string{
		"simple string",
		"line1\\nline2\\tcolumn",
		"path\\\\to\\\\file",
		"quote\\\"inside",
		"multiple\\nescapes\\tand\\\\backslash",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, value := range testValues {
			c.processEscapeSequences(value)
		}
	}
}
