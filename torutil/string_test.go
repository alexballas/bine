package torutil

import (
	"testing"
)

func TestPartitionString(t *testing.T) {
	assert := func(str string, ch byte, expectedA string, expectedB string, expectedOk bool) {
		t.Helper()
		a, b, ok := PartitionString(str, ch)
		if a != expectedA || b != expectedB || ok != expectedOk {
			t.Errorf("PartitionString(%q, %q) = (%q, %q, %v), want (%q, %q, %v)",
				str, ch, a, b, ok, expectedA, expectedB, expectedOk)
		}
	}
	assert("foo:bar", ':', "foo", "bar", true)
	assert(":bar", ':', "", "bar", true)
	assert("foo:", ':', "foo", "", true)
	assert("foo", ':', "foo", "", false)
	assert("foo:bar:baz", ':', "foo", "bar:baz", true)
}

func TestPartitionStringFromEnd(t *testing.T) {
	assert := func(str string, ch byte, expectedA string, expectedB string, expectedOk bool) {
		t.Helper()
		a, b, ok := PartitionStringFromEnd(str, ch)
		if a != expectedA || b != expectedB || ok != expectedOk {
			t.Errorf("PartitionStringFromEnd(%q, %q) = (%q, %q, %v), want (%q, %q, %v)",
				str, ch, a, b, ok, expectedA, expectedB, expectedOk)
		}
	}
	assert("foo:bar", ':', "foo", "bar", true)
	assert(":bar", ':', "", "bar", true)
	assert("foo:", ':', "foo", "", true)
	assert("foo", ':', "foo", "", false)
	assert("foo:bar:baz", ':', "foo:bar", "baz", true)
}

func TestEscapeSimpleQuotedStringIfNeeded(t *testing.T) {
	assert := func(str string, shouldBeDiff bool) {
		t.Helper()
		maybeEscaped := EscapeSimpleQuotedStringIfNeeded(str)
		if shouldBeDiff && maybeEscaped == str {
			t.Errorf("EscapeSimpleQuotedStringIfNeeded(%q) = %q, want it to be escaped", str, maybeEscaped)
		}
		if !shouldBeDiff && maybeEscaped != str {
			t.Errorf("EscapeSimpleQuotedStringIfNeeded(%q) = %q, want it unchanged", str, maybeEscaped)
		}
	}
	assert("foo", false)
	assert(" foo", true)
	assert("f\\oo", true)
	assert("fo\"o", true)
	assert("f\roo", true)
	assert("fo\no", true)
}

func TestEscapeSimpleQuotedString(t *testing.T) {
	if got := EscapeSimpleQuotedString("foo"); got != "\"foo\"" {
		t.Errorf("EscapeSimpleQuotedString(foo) = %q, want %q", got, "\"foo\"")
	}
}

func TestEscapeSimpleQuotedStringContents(t *testing.T) {
	assert := func(str string, expected string) {
		t.Helper()
		if got := EscapeSimpleQuotedStringContents(str); got != expected {
			t.Errorf("EscapeSimpleQuotedStringContents(%q) = %q, want %q", str, got, expected)
		}
	}
	assert("foo", "foo")
	assert("f\\oo", "f\\\\oo")
	assert("f\\noo", "f\\\\noo")
	assert("f\n o\ro", "f\\n o\\ro")
	assert("fo\r\\\"o", "fo\\r\\\\\\\"o")
}

func TestUnescapeSimpleQuotedStringIfNeeded(t *testing.T) {
	assert := func(str string, expectedStr string, expectedErr bool) {
		t.Helper()
		actualStr, actualErr := UnescapeSimpleQuotedStringIfNeeded(str)
		if actualStr != expectedStr || (actualErr != nil) != expectedErr {
			t.Errorf("UnescapeSimpleQuotedStringIfNeeded(%q) = (%q, %v), want (%q, err=%v)",
				str, actualStr, actualErr, expectedStr, expectedErr)
		}
	}
	assert("foo", "foo", false)
	assert("\"foo\"", "foo", false)
	assert("\"f\"oo\"", "", true)
}

func TestUnescapeSimpleQuotedString(t *testing.T) {
	assert := func(str string, expectedStr string, expectedErr bool) {
		t.Helper()
		actualStr, actualErr := UnescapeSimpleQuotedString(str)
		if actualStr != expectedStr || (actualErr != nil) != expectedErr {
			t.Errorf("UnescapeSimpleQuotedString(%q) = (%q, %v), want (%q, err=%v)",
				str, actualStr, actualErr, expectedStr, expectedErr)
		}
	}
	assert("foo", "", true)
	assert("\"foo\"", "foo", false)
	assert("\"f\"oo\"", "", true)
}

func TestUnescapeSimpleQuotedStringContents(t *testing.T) {
	assert := func(str string, expectedStr string, expectedErr bool) {
		t.Helper()
		actualStr, actualErr := UnescapeSimpleQuotedStringContents(str)
		if actualStr != expectedStr || (actualErr != nil) != expectedErr {
			t.Errorf("UnescapeSimpleQuotedStringContents(%q) = (%q, %v), want (%q, err=%v)",
				str, actualStr, actualErr, expectedStr, expectedErr)
		}
	}
	assert("foo", "foo", false)
	assert("f\\\\oo", "f\\oo", false)
	assert("f\\\\noo", "f\\noo", false)
	assert("f\\n o\\ro", "f\n o\ro", false)
	assert("fo\\r\\\\\\\"o", "fo\r\\\"o", false)
	assert("f\"oo", "", true)
	assert("f\roo", "", true)
	assert("f\noo", "", true)
	assert("f\\oo", "", true)
}
