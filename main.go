package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"

	"github.com/fatih/color"
)

type stringArray []string

func (slPtr *stringArray) String() string {
	return fmt.Sprintf("%v", *slPtr)
}

func (slPtr *stringArray) Set(value string) error {
	*slPtr = append(*slPtr, value)
	return nil
}

func printError(format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr,
		"%s "+format+"\n",
		append([]interface{}{color.RedString("!!!")}, params...)...)
}

func printInfo(format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr,
		"%s "+format+"\n",
		append([]interface{}{color.CyanString("***")}, params...)...)
}

func usage() {
	fmt.Fprintln(os.Stderr, "USAGE: mmm SRC_DIR TYPE_A=DST_DIR_A TYPE_B=DST_DIR_B ...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Example: mmm ~/import image=~/photos video=~/videos")
}

func validateDir(path string) {
	if fi, err := os.Stat(path); err != nil {
		printError("%v", err)
		os.Exit(1)
	} else if !fi.IsDir() {
		printError("%s is not a directory", path)
		os.Exit(1)
	}
}

func main() {

	parser := argparse.NewParser("mmm", "Distributes file based on mime type")
	sources := parser.List("s", "src",
		&argparse.Options{Required: true, Help: "Source directory"})
	dests := parser.List("d", "dst",
		&argparse.Options{Required: true, Help: "Map of mime type to  destination directory in the form `type=dir`"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(errors.New(color.RedString("!!! ") + err.Error())))
		return
	}
	mimeMap := map[string]string{}

	for _, dir := range *sources {
		validateDir(dir)
	}

	for _, m := range *dests {
		if !strings.Contains(m, "=") {
			printError("Invalid destination mapping %v. See help.\n", m)
			os.Exit(1)
		}
		fields := strings.SplitN(m, "=", 2)
		validateDir(fields[1])
		mimeMap[fields[0]] = fields[1]
	}

	if len(*sources) == 0 {
		printError("No sources specified")
		os.Exit(1)
	}

	if len(mimeMap) == 0 {
		printError("No destinations specified")
		os.Exit(1)
	}

	printInfo("Will read files from directories: %v", *sources)
	for t, d := range mimeMap {
		printInfo("... and copy files of type %s/* to: %s", t, d)
	}
}
