package output

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToSARIF_Basic(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/Foo.java")
	f.AddError("MyRule", "warning", "something bad", 42)

	got := toSARIF(r)

	if !strings.Contains(got, `"version": "2.1.0"`) {
		t.Errorf("expected SARIF version, got:\n%s", got)
	}
	if !strings.Contains(got, `"ruleId": "MyRule"`) {
		t.Errorf("expected ruleId, got:\n%s", got)
	}
	if !strings.Contains(got, `"uri": "/src/Foo.java"`) {
		t.Errorf("expected uri, got:\n%s", got)
	}
	if !strings.Contains(got, `"startLine": 42`) {
		t.Errorf("expected startLine, got:\n%s", got)
	}
	if !strings.Contains(got, `"text": "something bad"`) {
		t.Errorf("expected message text, got:\n%s", got)
	}
}

func TestToSARIF_SeverityLevel(t *testing.T) {
	cases := []struct {
		severity string
		want     string
	}{
		{"error", "error"},
		{"warning", "warning"},
		{"info", "note"},
		{"unknown", "note"},
	}
	for _, tc := range cases {
		got := sarifLevel(tc.severity)
		if got != tc.want {
			t.Errorf("severity %q: expected %q, got %q", tc.severity, tc.want, got)
		}
	}
}

func TestToSARIF_Rules(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/Foo.java")
	f.AddError("DupRule", "error", "msg1", 1)
	f.AddError("DupRule", "warning", "msg2", 2)
	f.AddError("OtherRule", "info", "msg3", 3)

	got := toSARIF(r)

	if count := strings.Count(got, `"id": "DupRule"`); count != 1 {
		t.Errorf("expected DupRule once in rules, got %d times", count)
	}
	if !strings.Contains(got, `"id": "OtherRule"`) {
		t.Errorf("expected OtherRule in rules, got:\n%s", got)
	}
}

func TestToSARIF_JSONValid(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/A.java")
	f.AddError("Rule", "error", "msg", 10)

	got := toSARIF(r)

	var v interface{}
	if err := json.Unmarshal([]byte(got), &v); err != nil {
		t.Errorf("invalid JSON: %v\n%s", err, got)
	}
}

func TestToSARIF_EmptyResult(t *testing.T) {
	r := &Result{}
	got := toSARIF(r)

	var v interface{}
	if err := json.Unmarshal([]byte(got), &v); err != nil {
		t.Errorf("invalid JSON for empty result: %v\n%s", err, got)
	}
	if !strings.Contains(got, `"rules": []`) {
		t.Errorf("expected empty rules array, got:\n%s", got)
	}
	if !strings.Contains(got, `"results": []`) {
		t.Errorf("expected empty results array, got:\n%s", got)
	}
}

func TestToSARIF_ToolName(t *testing.T) {
	r := &Result{Tool: "SpotBugs", ToolURI: "https://spotbugs.github.io/"}
	f := r.AddFile("/src/Foo.java")
	f.AddError("Rule", "error", "msg", 1)

	got := toSARIF(r)

	if !strings.Contains(got, `"name": "SpotBugs"`) {
		t.Errorf("expected tool name SpotBugs, got:\n%s", got)
	}
	if !strings.Contains(got, `"informationUri": "https://spotbugs.github.io/"`) {
		t.Errorf("expected informationUri spotbugs.github.io, got:\n%s", got)
	}
}

func TestToSARIF_ToolNameFallback(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/Foo.java")
	f.AddError("Rule", "error", "msg", 1)

	got := toSARIF(r)

	if !strings.Contains(got, `"name": "scarfco"`) {
		t.Errorf("expected fallback tool name scarfco, got:\n%s", got)
	}
	if !strings.Contains(got, `"informationUri": "https://github.com/mallowlabs/scarfco"`) {
		t.Errorf("expected fallback informationUri scarfco, got:\n%s", got)
	}
}

func TestConvert_SARIF(t *testing.T) {
	r := &Result{}
	f := r.AddFile("/src/A.java")
	f.AddError("Rule", "error", "msg", 10)

	got, err := Convert(r, "sarif")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty output")
	}
	if !strings.Contains(got, "2.1.0") {
		t.Errorf("expected SARIF version in output, got:\n%s", got)
	}
}
