package main

import (
	"flag"
	"fmt"
	"os"
	
	vss "github.com/jeromehadorn/vss"
)

var (
	drive       = flag.String("D", "", "Drive letter to copy (example: C:\\)")
	symlinkPath = flag.String("S", "", "Path of symlink folder")
	force       = flag.Bool("f", false, "Creates snapshots if available shadow storage is low. Warning: Replaces older snapshots.")
	timeout     = flag.Int("timeout", 180, "Snapshot creation timeout in seconds (min 180)")
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
		res, error := SymlinkSnapshot(*symlinkPath, snapshot.Id, snapshot.DeviceObjectPath)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
		fmt.Printf("Symlink to snapshot created: %s\n", res)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:  vss [options]
	`)
	flag.PrintDefaults()
	os.Exit(1)
}

func validate() {
	if *drive == "" {
		fmt.Fprintln(os.Stderr, `Error: Drive letter is required.`)
		usage()
	}
}

func checkUsage(nargs int) {
	if nargs > 0 {
		fmt.Fprintln(os.Stderr, `Unexpected arguments. Please see below all accepted arguments and their default values.`)
		usage()
	}
}
