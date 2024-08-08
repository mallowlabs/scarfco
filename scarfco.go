package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mallowlabs/scarfco/tools"
)

func init() {
	flag.Parse()
}

func read(filename string) ([]byte, error) {
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
	bytes, _ := io.ReadAll(r)

	return bytes, nil
}

func run() error {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	bytes, _ := read(filename)

	format, err := selectFormat(bytes)
	if err != nil {
		return err
	}

	var result *tools.Result = nil

	switch format {
	case "pmd":
		result = tools.ConvertPMD(bytes)
	case "pmd-cpd":
		result = tools.ConvertCPD(bytes)
	case "BugCollection":
		result = tools.ConvertFindBugs(bytes)
	default:
		return errors.New("unknown format error")
	}
	if result != nil {
		fmt.Println(result.ConvertToCheckstyle())
	}
	return nil
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

func main() {
	err := run()
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
