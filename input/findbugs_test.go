package input

import (
	"strings"
	"testing"
)

func TestConvertFindBugs_Basic(t *testing.T) {
	xml := []byte(`<?xml version='1.0'?>
<BugCollection>
  <file classname='com.example.App'>
    <BugInstance type='NP_NULL' priority='High' category='CORRECTNESS' message='null deref' lineNumber='10'/>
  </file>
  <Project/>
</BugCollection>`)

	r := convertFindBugs(xml)

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if !strings.HasSuffix(r.Files[0].Name, "com/example/App.java") {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(r.Files[0].Errors))
	}
	e := r.Files[0].Errors[0]
	if e.Source != "NP_NULL" {
		t.Errorf("expected source NP_NULL, got %s", e.Source)
	}
	if e.Line != 10 {
		t.Errorf("expected line 10, got %d", e.Line)
	}
}

func TestConvertFindBugs_SrcDir(t *testing.T) {
	xml := []byte(`<?xml version='1.0'?>
<BugCollection>
  <file classname='example.App'>
    <BugInstance type='HE' priority='Normal' category='BAD_PRACTICE' message='msg' lineNumber='5'/>
  </file>
  <Project>
    <SrcDir>/src/main/java</SrcDir>
  </Project>
</BugCollection>`)

	r := convertFindBugs(xml)

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if r.Files[0].Name != "/src/main/java/example/App.java" {
		t.Errorf("unexpected file path: %s", r.Files[0].Name)
	}
}

func TestConvertFindBugs_InnerClass(t *testing.T) {
	xml := []byte(`<?xml version='1.0'?>
<BugCollection>
  <file classname='example.App$Inner'>
    <BugInstance type='HE' priority='Low' category='BAD_PRACTICE' message='msg' lineNumber='1'/>
  </file>
  <Project/>
</BugCollection>`)

	r := convertFindBugs(xml)

	if !strings.HasSuffix(r.Files[0].Name, "example/App.java") {
		t.Errorf("inner class suffix should be stripped, got: %s", r.Files[0].Name)
	}
}

func TestConvertFindBugs_Priority(t *testing.T) {
	cases := []struct {
		priority string
		want     string
	}{
		{"High", "error"},
		{"Normal", "warning"},
		{"Low", "info"},
		{"Unknown", "info"},
	}

	for _, tc := range cases {
		got := severityFindBugs(tc.priority)
		if got != tc.want {
			t.Errorf("priority %q: expected %q, got %q", tc.priority, tc.want, got)
		}
	}
}

func TestConvertFindBugs_NativeFormat(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="utf-8"?>
<BugCollection>
  <Project>
    <SrcDir>/src/main/java</SrcDir>
  </Project>
  <BugInstance type='HE_EQUALS_USE_HASHCODE' priority='1'>
    <LongMessage>example.App defines equals and uses Object.hashCode()</LongMessage>
    <Class classname='example.App' primary='true'>
      <SourceLine classname='example.App' start='12' end='58' sourcepath='example/App.java'/>
    </Class>
    <Method primary='true' name='equals'>
      <SourceLine start='19' end='22' sourcepath='example/App.java'/>
    </Method>
    <SourceLine synthetic='true' start='19' end='22' sourcepath='example/App.java'/>
  </BugInstance>
</BugCollection>`)

	r := convertFindBugs(xml)

	if len(r.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(r.Files))
	}
	if r.Files[0].Name != "/src/main/java/example/App.java" {
		t.Errorf("unexpected file name: %s", r.Files[0].Name)
	}
	if len(r.Files[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(r.Files[0].Errors))
	}
	e := r.Files[0].Errors[0]
	if e.Source != "HE_EQUALS_USE_HASHCODE" {
		t.Errorf("expected source HE_EQUALS_USE_HASHCODE, got %s", e.Source)
	}
	if e.Severity != "error" {
		t.Errorf("expected severity error, got %s", e.Severity)
	}
	if e.Line != 19 {
		t.Errorf("expected line 19, got %d", e.Line)
	}
}

func TestConvertFindBugs_NativePriority(t *testing.T) {
	cases := []struct {
		priority string
		want     string
	}{
		{"1", "error"},
		{"2", "warning"},
		{"3", "info"},
		{"5", "info"},
	}

	for _, tc := range cases {
		got := severityFindBugsNative(tc.priority)
		if got != tc.want {
			t.Errorf("native priority %q: expected %q, got %q", tc.priority, tc.want, got)
		}
	}
}
