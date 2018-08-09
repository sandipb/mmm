package main

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
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

// WorkData info
type WorkData struct {
	src string
	dst string
}

func workFiles(src string, mimeMap map[string]string) chan WorkData {
	wd := make(chan WorkData)
	go func() {
		for fpath := range fileList(src) {
			mimeType := mime.TypeByExtension(filepath.Ext(fpath))
			if mimeType == "" {
				//printError("No type for '%s' found", f.path)
				continue
			}
			mimeTypeMajor := mimeType[:strings.IndexByte(mimeType, '/')]
			if d, ok := mimeMap[mimeTypeMajor]; ok {
				relPath, err := filepath.Rel(src, fpath)
				if err != nil {
					printError("Error finding relpath of %s: %v", fpath, err)
					continue
				}

				wd <- WorkData{src: fpath, dst: filepath.Join(d, relPath)}
			}

		}
		close(wd)
	}()
	return wd
}

func fileList(src string) chan string {
	fChan := make(chan string)
	go func() {
		filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
			if err == nil && !fi.IsDir() {
				fChan <- path
			}
			return nil
		})
		close(fChan)
	}()
	return fChan
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

	for _, source := range *sources {
		for w := range workFiles(source, mimeMap) {
			fmt.Println(w.src, "->", w.dst)
		}
	}
}
