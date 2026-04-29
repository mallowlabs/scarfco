package input

import "testing"

func TestSelectFormat_FindBugs(t *testing.T) {
	xml := []byte(`<?xml version='1.0'?><BugCollection version='4.9.8'></BugCollection>`)
	got, err := selectFormat(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "BugCollection" {
		t.Errorf("expected BugCollection, got %q", got)
	}
}

func TestSelectFormat_PMD(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?><pmd version="7.0.0"></pmd>`)
	got, err := selectFormat(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "pmd" {
		t.Errorf("expected pmd, got %q", got)
	}
}

func TestSelectFormat_CPD(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?><pmd-cpd></pmd-cpd>`)
	got, err := selectFormat(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "pmd-cpd" {
		t.Errorf("expected pmd-cpd, got %q", got)
	}
}

func TestSelectFormat_Checkstyle(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?><checkstyle version="9.3"></checkstyle>`)
	got, err := selectFormat(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "checkstyle" {
		t.Errorf("expected checkstyle, got %q", got)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?><unknown></unknown>`)
	_, err := Convert(xml)
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
