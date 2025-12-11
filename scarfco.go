package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mallowlabs/scarfco/input"
	"github.com/mallowlabs/scarfco/output"
)

func run() error {
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
		converted, err := output.Convert(result, "checkstyle")
		if err != nil {
			return err
		}
		fmt.Println(converted)
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
