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
	// xdoc format: BugInstance nested inside <file classname="..."> elements
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

	// native format: BugInstance is a direct child of BugCollection with nested SourceLine
	type NativeSourceLine struct {
		Sourcepath string `xml:"sourcepath,attr"`
		Start      int    `xml:"start,attr"`
		Synthetic  string `xml:"synthetic,attr"`
	}

	type NativeBugInstance struct {
		Type        string             `xml:"type,attr"`
		Priority    string             `xml:"priority,attr"`
		LongMessage string             `xml:"LongMessage"`
		SourceLines []NativeSourceLine `xml:"SourceLine"`
	}

	// nativeBug holds a parsed bug entry from the native format before grouping by file
	type nativeBug struct {
		filePath string
		bugType  string
		severity string
		message  string
		line     int
	}

	type BugCollection struct {
		XMLName    xml.Name           `xml:"BugCollection"`
		Files      []File             `xml:"file"`
		NativeBugs []NativeBugInstance `xml:"BugInstance"`
		Project    Project            `xml:"Project"`
	}

	var bc BugCollection
	xml.Unmarshal(content, &bc)

	srcDir := ""
	if len(bc.Project.SrcDirs) > 0 {
		srcDir = bc.Project.SrcDirs[0].Path
	}

	result := output.Result{}

	if len(bc.Files) > 0 {
		// xdoc format
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

	// native format: collect bugs then group by file path to produce one ResultFile per source file
	var bugs []nativeBug
	for _, bi := range bc.NativeBugs {
		sourcepath, start := "", 0
		for _, sl := range bi.SourceLines {
			if sl.Synthetic == "true" {
				sourcepath = sl.Sourcepath
				start = sl.Start
				break
			}
		}
		if sourcepath == "" && len(bi.SourceLines) > 0 {
			sourcepath = bi.SourceLines[0].Sourcepath
			start = bi.SourceLines[0].Start
		}
		if sourcepath == "" {
			continue
		}
		bugs = append(bugs, nativeBug{
			filePath: path.Join(srcDir, sourcepath),
			bugType:  bi.Type,
			severity: severityFindBugsNative(bi.Priority),
			message:  bi.LongMessage,
			line:     start,
		})
	}

	seen := map[string]bool{}
	var filePaths []string
	for _, b := range bugs {
		if !seen[b.filePath] {
			seen[b.filePath] = true
			filePaths = append(filePaths, b.filePath)
		}
	}
	for _, fp := range filePaths {
		f := result.AddFile(fp)
		for _, b := range bugs {
			if b.filePath == fp {
				f.AddError(b.bugType, b.severity, b.message, b.line)
			}
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

func severityFindBugsNative(priority string) string {
	switch priority {
	case "1":
		return "error"
	case "2":
		return "warning"
	default:
		return "info"
	}
}
