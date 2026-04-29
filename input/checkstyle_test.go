package input

import "testing"

func TestConvertCheckstyle_Basic(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<checkstyle version="9.3">
  <file name="/src/App.java">
    <error line="10" severity="error" message="Missing Javadoc" source="com.example.JavadocCheck"/>
    <error line="20" severity="warning" message="Line too long" source="com.example.LineLengthCheck"/>
  </file>
</checkstyle>`)

	r := convertCheckstyle(xml)

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if r.Files[0].Name != "/src/App.java" {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(r.Files[0].Errors))
	}

	e0 := r.Files[0].Errors[0]
	if e0.Line != 10 || e0.Message != "Missing Javadoc" || e0.Source != "com.example.JavadocCheck" {
		t.Errorf("unexpected first error: %+v", e0)
	}
}

func TestConvertCheckstyle_Severity(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"Warning", "warning"},
		{"INFO", "info"},
	}

	for _, tc := range cases {
		xml := []byte(`<?xml version="1.0"?>
<checkstyle>
  <file name="/src/A.java">
    <error line="1" severity="` + tc.input + `" message="m" source="s"/>
  </file>
</checkstyle>`)

		r := convertCheckstyle(xml)
		got := r.Files[0].Errors[0].Severity
		if got != tc.want {
			t.Errorf("severity %q: expected %q, got %q", tc.input, tc.want, got)
		}
	}
}
