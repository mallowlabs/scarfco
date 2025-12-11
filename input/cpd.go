package input

import (
	"encoding/xml"
	"fmt"

	"github.com/mallowlabs/scarfco/output"
)

func ConvertCPD(content []byte) *output.Result {
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

	result := output.Result{}
	for k, v := range m {
		file := result.AddFile(k)
		for _, rf := range v {
			message := fmt.Sprint(rf.Duplication.Lines) +
				" lines duplicated codes detected: " +
				rf.Another.Path + ":" + fmt.Sprint(rf.Another.Line)
			file.AddError("cpd", "warning", message, rf.File.Line)
		}
	}
	return &result
}
