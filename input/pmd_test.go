package input

import "testing"

func TestConvertPMD_Basic(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<pmd version="7.0.0">
  <file name="/src/App.java">
    <violation beginline="18" priority="3" rule="OverrideBothEqualsAndHashcode">
      Ensure you override both equals() and hashCode()
    </violation>
    <violation beginline="37" priority="3" rule="UnusedLocalVariable">
      Avoid unused local variables such as 'b'.
    </violation>
    <violation beginline="38" priority="3" rule="EmptyCatchBlock">
      Avoid empty catch blocks
    </violation>
  </file>
</pmd>`)

	r := convertPMD(xml)

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if r.Files[0].Name != "/src/App.java" {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(r.Files[0].Errors))
	}

	e0 := r.Files[0].Errors[0]
	if e0.Source != "OverrideBothEqualsAndHashcode" || e0.Line != 18 {
		t.Errorf("unexpected first error: %+v", e0)
	}
}

func TestConvertPMD_Priority(t *testing.T) {
	cases := []struct {
		priority int
		want     string
	}{
		{1, "error"},
		{2, "error"},
		{3, "warning"},
		{4, "warning"},
		{5, "info"},
	}

	for _, tc := range cases {
		got := severityPMD(tc.priority)
		if got != tc.want {
			t.Errorf("priority %d: expected %q, got %q", tc.priority, tc.want, got)
		}
	}
}

func TestConvertPMD_Message(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<pmd version="7.0.0">
  <file name="/src/App.java">
    <violation beginline="1" priority="3" rule="SomeRule">
      ` + "   trimmed message   " + `
    </violation>
  </file>
</pmd>`)

	r := convertPMD(xml)

	if len(r.Files[0].Errors) != 1 {
		t.Fatalf("expected 1 error")
	}
	msg := r.Files[0].Errors[0].Message
	if msg != "trimmed message" {
		t.Errorf("expected trimmed message, got %q", msg)
	}
}
