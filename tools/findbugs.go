package tools

import (
	"encoding/xml"
	"path"
)

func ConvertFindBugs(content []byte) *Result {
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

	result := Result{}
	for k, v := range m {
		file := result.AddFile(path.Join(srcDir, k))
		for _, bi := range v {
			file.AddError(bi.Type, severityFindBugs(bi.Priority), bi.LongMessage.Message, bi.SourceLine.Start)
		}
	}
	return &result
}

func severityFindBugs(priority int) string {
	if priority == 1 {
		return "error"
	} else if priority == 2 {
		return "warning"
	} else {
		return "info"
	}
}
