package main_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mallowlabs/scarfco/input"
)

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("failed to read testdata/%s: %v", name, err)
	}
	return b
}

func TestIntegration_Spotbugs(t *testing.T) {
	content := readTestdata(t, "spotbugs.xml")
	r, err := input.Convert(content)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if !strings.HasSuffix(r.Files[0].Name, "example/App.java") {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(r.Files[0].Errors))
	}
	e := r.Files[0].Errors[0]
	if e.Source != "HE_EQUALS_USE_HASHCODE" {
		t.Errorf("expected source HE_EQUALS_USE_HASHCODE, got %q", e.Source)
	}
	if e.Severity != "error" {
		t.Errorf("expected severity error, got %q", e.Severity)
	}
	if e.Line != 19 {
		t.Errorf("expected line 19, got %d", e.Line)
	}
}

func TestIntegration_PMD(t *testing.T) {
	content := readTestdata(t, "pmd.xml")
	r, err := input.Convert(content)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if len(r.Files[0].Errors) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(r.Files[0].Errors))
	}

	expected := []struct {
		source string
		line   int
	}{
		{"OverrideBothEqualsAndHashcode", 18},
		{"UnusedLocalVariable", 37},
		{"EmptyCatchBlock", 38},
	}
	for i, ex := range expected {
		e := r.Files[0].Errors[i]
		if e.Source != ex.source || e.Line != ex.line {
			t.Errorf("error[%d]: expected {%s line=%d}, got {%s line=%d}", i, ex.source, ex.line, e.Source, e.Line)
		}
		if e.Severity != "warning" {
			t.Errorf("error[%d]: expected severity warning, got %q", i, e.Severity)
		}
	}
}

func TestIntegration_CPD(t *testing.T) {
	content := readTestdata(t, "cpd.xml")
	r, err := input.Convert(content)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	if len(r.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(r.Files))
	}

	paths := map[string]bool{}
	for _, f := range r.Files {
		paths[f.Name] = true
		if len(f.Errors) != 1 {
			t.Errorf("file %s: expected 1 error, got %d", f.Name, len(f.Errors))
		}
		if f.Errors[0].Severity != "warning" {
			t.Errorf("file %s: expected severity warning, got %q", f.Name, f.Errors[0].Severity)
		}
		if !strings.Contains(f.Errors[0].Message, "10 lines duplicated") {
			t.Errorf("file %s: unexpected message: %q", f.Name, f.Errors[0].Message)
		}
	}

	if !paths["/home/runner/work/scarfco-example/scarfco-example/src/main/java/example/App.java"] {
		t.Errorf("App.java not found in result files: %v", paths)
	}
	if !paths["/home/runner/work/scarfco-example/scarfco-example/src/main/java/example/Greeter.java"] {
		t.Errorf("Greeter.java not found in result files: %v", paths)
	}
}

func TestIntegration_SpotbugsXml(t *testing.T) {
	content := readTestdata(t, "spotbugsXml.xml")
	r, err := input.Convert(content)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if !strings.HasSuffix(r.Files[0].Name, "example/App.java") {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(r.Files[0].Errors))
	}
	e := r.Files[0].Errors[0]
	if e.Source != "HE_EQUALS_USE_HASHCODE" {
		t.Errorf("expected source HE_EQUALS_USE_HASHCODE, got %q", e.Source)
	}
	if e.Severity != "error" {
		t.Errorf("expected severity error, got %q", e.Severity)
	}
	if e.Line != 19 {
		t.Errorf("expected line 19, got %d", e.Line)
	}
}

func TestIntegration_Checkstyle(t *testing.T) {
	content := readTestdata(t, "checkstyle-result.xml")
	r, err := input.Convert(content)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	if len(r.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(r.Files))
	}

	total := 0
	for _, f := range r.Files {
		total += len(f.Errors)
	}
	if total != 26 {
		t.Errorf("expected 26 errors total, got %d", total)
	}

	for _, f := range r.Files {
		for _, e := range f.Errors {
			if e.Severity != "error" {
				t.Errorf("file %s: expected all severities to be 'error', got %q", f.Name, e.Severity)
			}
		}
	}
}
