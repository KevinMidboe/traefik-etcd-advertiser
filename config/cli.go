package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	filename string
)

func ParseCli() (string, *bool) {
	flag.StringVar(&filename, "filename", "", "path of config")
	publish := flag.Bool("publish", false, "publish etcd messages")

	flag.Parse()

	args := os.Args[1:]
	if len(args) == 0 {
		// no command, exit with code 2 (invalid usage)
	  fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

    flag.PrintDefaults()
		os.Exit(2)
	}

	if len(filename) < 1 {
		fmt.Fprintf(os.Stderr, "Filename required. Usage:\n")

		flag.PrintDefaults()
		os.Exit(2)
	}

	return filename, publish
}
