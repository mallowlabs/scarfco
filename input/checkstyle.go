package input

import (
	"encoding/xml"
	"strings"

	"github.com/mallowlabs/scarfco/output"
)

func convertCheckstyle(content []byte) *output.Result {
	type Error struct {
		Line     int    `xml:"line,attr"`
		Severity string `xml:"severity,attr"`
		Message  string `xml:"message,attr"`
		Source   string `xml:"source,attr"`
	}

	type File struct {
		Name   string  `xml:"name,attr"`
		Errors []Error `xml:"error"`
	}

	type Checkstyle struct {
		Files []File `xml:"file"`
	}

	var checkstyle Checkstyle
	xml.Unmarshal(content, &checkstyle)

	result := output.Result{}
	for _, file := range checkstyle.Files {
		f := result.AddFile(file.Name)
		for _, err := range file.Errors {
			f.AddError(err.Source, strings.ToLower(err.Severity), err.Message, err.Line)
		}
	}
	return &result
}
