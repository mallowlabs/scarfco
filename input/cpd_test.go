package input

import (
	"strings"
	"testing"
)

func TestConvertCPD_Basic(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<pmd-cpd>
  <duplication lines="10" tokens="50">
    <file line="44" path="/src/App.java"/>
    <file line="19" path="/src/Greeter.java"/>
  </duplication>
</pmd-cpd>`)

	r := convertCPD(xml)

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
			t.Errorf("file %s: expected severity warning, got %s", f.Name, f.Errors[0].Severity)
		}
		if f.Errors[0].Source != "cpd" {
			t.Errorf("file %s: expected source cpd, got %s", f.Name, f.Errors[0].Source)
		}
	}

	if !paths["/src/App.java"] || !paths["/src/Greeter.java"] {
		t.Errorf("expected both file paths in result, got: %v", paths)
	}
}

func TestConvertCPD_Message(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<pmd-cpd>
  <duplication lines="10" tokens="50">
    <file line="44" path="/src/App.java"/>
    <file line="19" path="/src/Greeter.java"/>
  </duplication>
</pmd-cpd>`)

	r := convertCPD(xml)

	for _, f := range r.Files {
		msg := f.Errors[0].Message
		if !strings.Contains(msg, "10 lines duplicated codes detected") {
			t.Errorf("unexpected message format: %q", msg)
		}
		if f.Name == "/src/App.java" {
			if !strings.Contains(msg, "/src/Greeter.java:19") {
				t.Errorf("App.java message should reference Greeter.java:19, got: %q", msg)
			}
		}
		if f.Name == "/src/Greeter.java" {
			if !strings.Contains(msg, "/src/App.java:44") {
				t.Errorf("Greeter.java message should reference App.java:44, got: %q", msg)
			}
		}
	}
}
