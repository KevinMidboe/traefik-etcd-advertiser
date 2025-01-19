package config

import (
	"flag"
	"fmt"
	"os"
)

var filename string
const applicationName = "traefik-etcd-advertiser"

func setupHelpMessage() {
	flag.Usage = func() {
		helpMessage()
	}
}

func helpMessage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [OPTIONS]\n\n", applicationName)
	fmt.Fprintf(os.Stderr, "Options:\n")

	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nHelp Options:\n  -h, --help\tShow this help message")
}

func printVersion(v string) {
	fmt.Printf("Version: %s\n", v)
}

func ParseCli(version string) (string, *bool) {
	setupHelpMessage()

	flag.StringVar(&filename, "filename", "", "path of config")
	publish := flag.Bool("publish", false, "publish etcd messages")
	versionFlag := flag.Bool("version", false, "print version information")

	flag.Parse()

	args := os.Args[1:]
	if *versionFlag {
		printVersion(version)
		os.Exit(0)
	}

	if len(args) == 0 {
		// no command, exit with code 2 (invalid usage)
		helpMessage()

		os.Exit(2)
	}

	if len(filename) < 1 {
		fmt.Fprintf(os.Stderr, "Filename required. Usage:\n")

		helpMessage()
		os.Exit(2)
	}

	return filename, publish
}
