package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/jeromehadorn/vss"
)

var (
	drive       = flag.String("D", "", "Drive letter to copy (required)")
	symlinkPath = flag.String("S", "", "Path of symlink folder that points to the snapshot")
	force       = flag.Bool("f", false, "Creates snapshots if available shadow storage is low. Could delete old copies!")
	timeout     = flag.Int("timeout", 180, "Snapshot creation timeout in seconds (minimum 180)")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	checkUsage(flag.NArg())
	validate()

	Snapshotter := vss.Snapshotter{}
	snapshot, err := Snapshotter.CreateSnapshot(*drive, *timeout, *force)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Snapshot created: %s\n", snapshot.Id)

	if symlinkPath != nil {
		// SYMLINK here...
		res, error := SymlinkSnapshot(*symlinkPath, snapshot.Id, snapshot.DeviceObjectPath)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
		fmt.Printf("Symlink to snapshot created: %s\n", res)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:  ........
	`)
	flag.PrintDefaults()
	os.Exit(1)
}

func validate() {
	if *drive == "" {
		fmt.Fprintln(os.Stderr, `Error: Drive letter is required.`)
		usage()
	}
	drivePattern := "[d-zD-Z]:"
	match, _ := regexp.MatchString(drivePattern, *drive)
	if !match {
		fmt.Fprintln(os.Stderr, `Drive letter is invalid.`)
		usage()
		os.Exit(1)
	}

	if *symlinkPath != "" {
		// TODO: Ensure it's a valid path
	}
}

func checkUsage(nargs int) {
	if nargs > 0 {
		fmt.Fprintln(os.Stderr, `Unexpected arguments. Please see below all accepted arguments and their default values.`)
		usage()
	}
}
