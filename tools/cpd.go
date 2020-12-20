package tools

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

func ConvertCPD(content []byte) string {
	type File struct {
		Line int    `xml:"line,attr"`
		Path string `xml:"path,attr"`
	}

	type Duplication struct {
		Files []File `xml:"file"`
		Lines int    `xml:"lines,attr"`
	}

	type PmdCpd struct {
		XMLName      xml.Name      `xml:"pmd-cpd"`
		Duplications []Duplication `xml:"duplication"`
	}

	type ResultFile struct {
		File        File
		Another     File
		Duplication Duplication
	}

	var cpd PmdCpd
	xml.Unmarshal(content, &cpd)

	var buf bytes.Buffer
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<checkstyle version=\"5.0\">\n")

	m := map[string][]ResultFile{}

	for _, duplication := range cpd.Duplications {
		for index, file := range duplication.Files {
			another := duplication.Files[1]
			if index == 1 {
				another = duplication.Files[0]
			}

			rf := ResultFile{File: file, Another: another, Duplication: duplication}

			v, ok := m[file.Path]
			if ok {
				m[file.Path] = append(v, rf)
			} else {
				m[file.Path] = []ResultFile{rf}
			}
		}
	}

	for k, v := range m {
		buf.WriteString("  <file name=\"")
		xml.Escape(&buf, []byte(k))
		buf.WriteString("\">\n")

		for _, rf := range v {
			buf.WriteString("    <error ")

			buf.WriteString("line=\"")
			buf.WriteString(fmt.Sprint(rf.File.Line))
			buf.WriteString("\" ")

			buf.WriteString("severity=\"warning\" ")

			buf.WriteString("message=\"")
			buf.WriteString(fmt.Sprint(rf.Duplication.Lines))
			buf.WriteString(" lines duplicated codes detected: ")

			xml.Escape(&buf, []byte(rf.Another.Path))
			buf.WriteString(":")
			buf.WriteString(fmt.Sprint(rf.Another.Line))
			buf.WriteString("\" ")

			buf.WriteString("source=\"cpd\"/>\n")
		}
		buf.WriteString("  </file>\n")
	}

	buf.WriteString("</checkstyle>")
	return buf.String()
}
