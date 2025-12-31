package input

import (
	"encoding/xml"
	"strings"

	"github.com/mallowlabs/scarfco/output"
)

func init() {
	RegisterConverter("pmd", convertPMD)
}

func convertPMD(content []byte) *output.Result {
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

	result := output.Result{}
	for _, file := range pmd.Files {
		f := result.AddFile(file.Name)
		for _, violation := range file.Violations {
			f.AddError(violation.Rule, severityPMD(violation.Priority), strings.TrimSpace(violation.Message), violation.Beginline)
		}
	}
	return &result
}

func severityPMD(priority int) string {
	switch priority {
	case 1, 2:
		return "error"
	case 3, 4:
		return "warning"
	default:
		return "info"
	}
}
