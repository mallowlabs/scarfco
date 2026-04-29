package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/mallowlabs/scarfco/input"
	"github.com/mallowlabs/scarfco/output"
)

func run(format string) error {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	bytes, err := input.Read(filename)
	if err != nil {
		return err
	}

	result, err := input.Convert(bytes)
	if err != nil {
		return err
	}
	if result != nil {
		converted, err := output.Convert(result, format)
		if err != nil {
			return err
		}
		fmt.Println(converted)
	}
	return nil
}

func main() {
	var showVersion bool
	var format string
	flag.BoolVar(&showVersion, "v", false, "print the version")
	flag.StringVar(&format, "format", "checkstyle", "output format: checkstyle, sarif")
	flag.Parse()

	if showVersion {
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(info.Main.Version)
		} else {
			fmt.Fprintln(os.Stderr, "Error: could not read build info")
			os.Exit(1)
		}
		return
	}

	err := run(format)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
