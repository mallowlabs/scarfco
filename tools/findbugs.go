package tools

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
)

func ConvertFindBugs(content []byte) string {
	type SourceLine struct {
		Start      int    `xml:"start,attr"`
		Sourcepath string `xml:"sourcepath,attr"`
	}

	type LongMessage struct {
		Message string `xml:",cdata"`
	}

	type BugInstance struct {
		Type        string      `xml:"type,attr"`
		LongMessage LongMessage `xml:"LongMessage"`
		SourceLine  SourceLine  `xml:"SourceLine"`
		Priority    int         `xml:"priority,attr"`
	}

	type SrcDir struct {
		Path string `xml:",cdata"`
	}
	type Project struct {
		SrcDirs []SrcDir `xml:"SrcDir"`
	}

	type BugCollection struct {
		XMLName      xml.Name      `xml:"BugCollection"`
		Project      Project       `xml:"Project"`
		BugInstances []BugInstance `xml:"BugInstance"`
	}

	var bc BugCollection
	xml.Unmarshal(content, &bc)

	var buf bytes.Buffer
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<checkstyle version=\"5.0\">\n")

	m := map[string][]BugInstance{}

	srcDir := bc.Project.SrcDirs[0].Path

	for _, bi := range bc.BugInstances {
		v, ok := m[bi.SourceLine.Sourcepath]
		if ok {
			m[bi.SourceLine.Sourcepath] = append(v, bi)
		} else {
			m[bi.SourceLine.Sourcepath] = []BugInstance{bi}
		}
	}

	for k, v := range m {
		buf.WriteString("  <file name=\"")
		xml.Escape(&buf, []byte(path.Join(srcDir, k)))
		buf.WriteString("\">\n")

		for _, bi := range v {
			buf.WriteString("    <error ")

			buf.WriteString("line=\"")
			buf.WriteString(fmt.Sprint(bi.SourceLine.Start))
			buf.WriteString("\" ")

			buf.WriteString("severity=\"")
			buf.WriteString(fmt.Sprint(bi.Priority))
			buf.WriteString("\" ")

			buf.WriteString("message=\"")
			xml.Escape(&buf, []byte(bi.LongMessage.Message))
			buf.WriteString("\" ")

			buf.WriteString("source=\"")
			xml.Escape(&buf, []byte(bi.Type))
			buf.WriteString("\"/>\n")
		}
		buf.WriteString("  </file>\n")
	}

	buf.WriteString("</checkstyle>")
	return buf.String()
}
