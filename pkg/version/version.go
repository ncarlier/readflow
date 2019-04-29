package version

import (
	"fmt"
)

// Version of the app
var Version = "snapshot"

// GitCommit is the GIT commit revision
var GitCommit = "n/a"

// Built is the built date
var Built = "n/a"

// Print version to stdout
func Print() {
	fmt.Printf(`Version:    %s
Git commit: %s
Built:      %s

Copyright (C) 2019 Nunux, Org.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by Nicolas Carlier.
`, Version, GitCommit, Built)
}
