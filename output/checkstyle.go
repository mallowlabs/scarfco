package output

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// ToChekstyle converts a Result object to a Checkstyle XML string.
func ToChekstyle(r *Result) string {
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
