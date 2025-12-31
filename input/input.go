package input

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"github.com/mallowlabs/scarfco/output"
)

// Converter is a function that converts from a byte slice to a Result.
type Converter func([]byte) *output.Result

var converters = make(map[string]Converter)

// RegisterConverter registers a converter for a given format.
func RegisterConverter(format string, converter Converter) {
	converters[format] = converter
}

func Read(filename string) ([]byte, error) {
	var r io.Reader
	switch filename {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func Convert(content []byte) (*output.Result, error) {
	format, err := selectFormat(content)
	if err != nil {
		return nil, err
	}

	converter, ok := converters[format]
	if !ok {
		return nil, errors.New("unknown format: " + format)
	}

	return converter(content), nil
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
		default:
		}
		if format != "" {
			break
		}
	}
	return format, nil
}
