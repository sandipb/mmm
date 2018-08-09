package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "USAGE: mmm SRC_DIR TYPE_A=DST_DIR_A TYPE_B=DST_DIR_B ...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Example: mmm ~/import image=~/photos video=~/videos")
}

func main() {

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
}
