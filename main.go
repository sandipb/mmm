package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type stringArray []string

func (slPtr *stringArray) String() string {
	return fmt.Sprintf("%v", *slPtr)
}

func (slPtr *stringArray) Set(value string) error {
	*slPtr = append(*slPtr, value)
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "USAGE: mmm SRC_DIR TYPE_A=DST_DIR_A TYPE_B=DST_DIR_B ...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Example: mmm ~/import image=~/photos video=~/videos")
}

func validateDir(path string) {
	if fi, err := os.Stat(path); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	} else if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a directory\n", path)
		os.Exit(1)
	}
}

func main() {

	var sources, dests stringArray

	flag.Var(&sources, "src", "Source directory `path`. Can be repeated.")
	flag.Var(&dests, "dst", "Map of mime type to  destination directory in the form `type=dir`. Can be repeated.")
	flag.Parse()

	mimeMap := map[string]string{}

	for _, dir := range sources {
		validateDir(dir)
	}

	for _, m := range dests {
		if !strings.Contains(m, "=") {
			fmt.Fprintf(os.Stderr, "ERROR: Invalid destination mapping %v. See help.\n", m)
			os.Exit(1)
		}
		fields := strings.SplitN(m, "=", 2)
		validateDir(fields[1])
		mimeMap[fields[0]] = fields[1]
	}

	if len(sources) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: No sources specified")
		os.Exit(1)
	}

	if len(mimeMap) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: No destinations specified")
		os.Exit(1)
	}

	fmt.Printf("*** Will read files from directories: %v\n", sources)
	for t, d := range mimeMap {
		fmt.Printf("*** ... and copy files of type %s/* to: %s\n", t, d)
	}
}
