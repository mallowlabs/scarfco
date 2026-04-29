package output

import (
	"strings"
	"testing"
)

func TestToCheckstyle_Basic(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/Foo.java")
	f.AddError("MyRule", "warning", "something bad", 42)

	got := toChekstyle(r)

	if !strings.Contains(got, `<file name="/src/Foo.java">`) {
		t.Errorf("expected file element, got:\n%s", got)
	}
	if !strings.Contains(got, `line="42"`) {
		t.Errorf("expected line=42, got:\n%s", got)
	}
	if !strings.Contains(got, `severity="warning"`) {
		t.Errorf("expected severity=warning, got:\n%s", got)
	}
	if !strings.Contains(got, `message="something bad"`) {
		t.Errorf("expected message, got:\n%s", got)
	}
	if !strings.Contains(got, `source="MyRule"`) {
		t.Errorf("expected source=MyRule, got:\n%s", got)
	}
}

func TestToCheckstyle_XmlEscape(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/<Bad>&File.java")
	f.AddError(`a"b'c`, "error", `msg with <>&"'`, 1)

	got := toChekstyle(r)

	if strings.Contains(got, `<Bad>`) {
		t.Errorf("< and > should be escaped in file name, got:\n%s", got)
	}
	if strings.Contains(got, `<>&"'`) {
		t.Errorf("special chars should be escaped in message, got:\n%s", got)
	}
	if !strings.Contains(got, `&lt;`) {
		t.Errorf("expected &lt; in output, got:\n%s", got)
	}
}

func TestConvert_Checkstyle(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/A.java")
	f.AddError("Rule", "error", "msg", 10)

	got, err := Convert(r, "checkstyle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, "<?xml") {
		t.Errorf("expected XML declaration, got:\n%s", got)
	}
	if !strings.Contains(got, `<checkstyle version="5.0">`) {
		t.Errorf("expected checkstyle root element, got:\n%s", got)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	r := &Result{}
	_, err := Convert(r, "unknown")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
