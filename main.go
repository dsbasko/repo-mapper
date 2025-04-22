package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	summarizer "github.com/dsbasko/repo-mapper/internal"
)

func main() {
	rootDir := flag.String("dir", ".", "root directory of the project for scanning")
	outputFileName := flag.String("o", "summary.md", "Output summary file name")

	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %v\n", flag.Args())
		fmt.Fprintf(os.Stderr, "use -h or --help for help.\n")
		os.Exit(1)
	}

	absRootDir, err := filepath.Abs(*rootDir)
	if err != nil {
		_ = fmt.Errorf("could not get absolute path for directory '%s': %v", *rootDir, err)
	}

	info, err := os.Stat(absRootDir)
	if err != nil {
		if os.IsNotExist(err) {
			_ = fmt.Errorf("directory '%s' not found.", absRootDir)
		}
		_ = fmt.Errorf("could not get information about directory '%s': %v", absRootDir, err)
	}
	if !info.IsDir() {
	_ = fmt.Errorf("path '%s' is not a directory.", absRootDir)
	}

	err = summarizer.SummarizeProject(absRootDir, *outputFileName)
	if err != nil {
		_ = fmt.Errorf("error when creating %s: %v", *outputFileName, err)
	}
}
