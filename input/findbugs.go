package input

import (
	"encoding/xml"
	"path"
	"strings"

	"github.com/mallowlabs/scarfco/output"
)

func init() {
	RegisterConverter("BugCollection", convertFindBugs)
}

func convertFindBugs(content []byte) *output.Result {
	type BugInstance struct {
		Type       string `xml:"type,attr"`
		Priority   string `xml:"priority,attr"`
		Category   string `xml:"category,attr"`
		Message    string `xml:"message,attr"`
		LineNumber int    `xml:"lineNumber,attr"`
	}

	type File struct {
		ClassName    string        `xml:"classname,attr"`
		BugInstances []BugInstance `xml:"BugInstance"`
	}

	type SrcDir struct {
		Path string `xml:",cdata"`
	}
	type Project struct {
		SrcDirs []SrcDir `xml:"SrcDir"`
	}

	type BugCollection struct {
		XMLName xml.Name `xml:"BugCollection"`
		Files   []File   `xml:"file"`
		Project Project  `xml:"Project"`
	}

	var bc BugCollection
	xml.Unmarshal(content, &bc)

	srcDir := ""
	if len(bc.Project.SrcDirs) > 0 {
		srcDir = bc.Project.SrcDirs[0].Path
	}

	result := output.Result{}
	for _, file := range bc.Files {
		// com.example.MyClass$InnerClass -> com/example/MyClass.java
		className := file.ClassName
		if idx := strings.Index(className, "$"); idx != -1 {
			className = className[:idx]
		}
		filePath := path.Join(srcDir, strings.ReplaceAll(className, ".", "/")+".java")
		f := result.AddFile(filePath)
		for _, bi := range file.BugInstances {
			f.AddError(bi.Type, severityFindBugs(bi.Priority), bi.Message, bi.LineNumber)
		}
	}
	return &result
}

func severityFindBugs(priority string) string {
	switch priority {
	case "High":
		return "error"
	case "Normal":
		return "warning"
	case "Low":
		return "info"
	default:
		return "info"
	}
}
