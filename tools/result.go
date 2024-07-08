package tools

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type Result struct {
	Files []ResultFile
}

type ResultFile struct {
	Name   string
	Errors []ResultError
}

type ResultError struct {
	Source   string
	Severity string
	Message  string
	Line     int
}

func (r *Result) AddFile(name string) *ResultFile {
	f := ResultFile{Name: name}
	r.Files = append(r.Files, f)
	return &r.Files[len(r.Files)-1]
}

func (f *ResultFile) AddError(source, severity, message string, line int) {
	e := ResultError{Source: source, Severity: severity, Message: message, Line: line}
	f.Errors = append(f.Errors, e)
}

func (r *Result) ConvertToCheckstyle() string {
	var buf bytes.Buffer
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<checkstyle version=\"5.0\">\n")

	for _, f := range r.Files {
		buf.WriteString("  <file name=\"")
		xml.Escape(&buf, []byte(f.Name))
		buf.WriteString("\">\n")

		for _, e := range f.Errors {
			buf.WriteString("    <error ")

			buf.WriteString("line=\"")
			buf.WriteString(fmt.Sprint(e.Line))
			buf.WriteString("\" ")

			buf.WriteString("severity=\"")
			buf.WriteString(e.Severity)
			buf.WriteString("\" ")

			buf.WriteString("message=\"")
			xml.Escape(&buf, []byte(e.Message))
			buf.WriteString("\" ")

			buf.WriteString("source=\"")
			xml.Escape(&buf, []byte(e.Source))
			buf.WriteString("\"/>\n")
		}
		buf.WriteString("  </file>\n")
	}

	buf.WriteString("</checkstyle>")
	return buf.String()
}
