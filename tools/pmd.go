package tools

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func ConvertPMD(content []byte) string {
	type Violation struct {
		Beginline int    `xml:"beginline,attr"`
		Priority  int    `xml:"priority,attr"`
		Rule      string `xml:"rule,attr"`
		Message   string `xml:",cdata"`
	}

	type File struct {
		Violations []Violation `xml:"violation"`
		Name       string      `xml:"name,attr"`
	}

	type Pmd struct {
		XMLName xml.Name `xml:"pmd"`
		Files   []File   `xml:"file"`
		Version string   `xml:"version,attr"`
	}

	var pmd Pmd
	xml.Unmarshal(content, &pmd)

	var buf bytes.Buffer
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<checkstyle version=\"5.0\">\n")

	for _, file := range pmd.Files {
		buf.WriteString("  <file name=\"")
		xml.Escape(&buf, []byte(file.Name))
		buf.WriteString("\">\n")
		for _, violation := range file.Violations {
			buf.WriteString("    <error ")

			buf.WriteString("line=\"")
			buf.WriteString(fmt.Sprint(violation.Beginline))
			buf.WriteString("\" ")

			buf.WriteString("severity=\"")
			buf.WriteString(severityPMD(violation.Priority))
			buf.WriteString("\" ")

			buf.WriteString("message=\"")
			xml.Escape(&buf, []byte(strings.TrimSpace(violation.Message)))
			buf.WriteString("\" ")

			buf.WriteString("source=\"")
			xml.Escape(&buf, []byte(violation.Rule))
			buf.WriteString("\"/>\n")
		}
		buf.WriteString("  </file>\n")
	}

	buf.WriteString("</checkstyle>")
	return buf.String()
}

func severityPMD(priority int) string {
	if priority == 1 || priority == 2 {
		return "error"
	} else if priority == 3 || priority == 4 {
		return "warning"
	} else {
		return "info"
	}
}
