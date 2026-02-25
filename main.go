package main

import "github.com/vsyaco/html2md/cmd"

// Injected via -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
