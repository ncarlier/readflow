package version

import (
	"fmt"
)

// Version of the app
var Version = "snapshot"

func Print() {
	fmt.Printf(`webhookd (%s)
Copyright (C) 2018 Nunux, Org.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by Nicolas Carlier.`, Version)
}
