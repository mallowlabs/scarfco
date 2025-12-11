package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mallowlabs/scarfco/input"
	"github.com/mallowlabs/scarfco/output"
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

	result, err := input.Convert(bytes)
	if err != nil {
		return err
	}
	if result != nil {
		fmt.Println(output.ToChekstyle(result))
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
