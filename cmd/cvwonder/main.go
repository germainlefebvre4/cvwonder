package main

import "os"

var (
	Commit  = ""
	Version = "x.x.x"
	Date    = ""
)

func main() {
	os.Exit(run(os.Args[1:]))
}
