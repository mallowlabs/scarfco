package input

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/mallowlabs/scarfco/output"
)

func Convert(content []byte) (*output.Result, error) {
	format, err := selectFormat(content)
	if err != nil {
		return nil, err
	}

	var result *output.Result = nil

	switch format {
	case "pmd":
		result = ConvertPMD(content)
	case "pmd-cpd":
		result = ConvertCPD(content)
	case "BugCollection":
		result = ConvertFindBugs(content)
	default:
		return nil, errors.New("unknown format error")
	}
	return result, nil
}

func selectFormat(content []byte) (string, error) {
	d := xml.NewDecoder(bytes.NewReader(content))
	format := ""

	for {
		token, err := d.Token()

		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return "", err
		}
		switch t := token.(type) {
		case xml.StartElement:
			format = t.Name.Local
			break
		default:
			break
		}
		if format != "" {
			break
		}
	}
	return format, nil
}
