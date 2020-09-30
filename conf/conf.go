package conf

import (
	"flag"
	"fmt"

	"get.cutie.cafe/levelup/util"
)

var (
	// OpMode describes the server's operation mode
	OpMode *string

	// Fetch enables fetching
	Fetch bool = true

	// Serve enables serving
	Serve bool = true
)

// Init finds the command line flags and sets FetchOnly/ServeOnly
func Init() {
	OpMode = flag.String("mode", "all", "Operation mode. Possible options are fetch, serve, and all (which is both).")
	flag.Parse()

	util.Debug(fmt.Sprintf("OpMode: %s", *OpMode))

	switch *OpMode {
	case "fetch":
		Fetch = true
		Serve = false
		break
	case "serve":
		Fetch = false
		Serve = true
		break
	case "all":
		Fetch = true
		Serve = true
		break
	}

	util.Debug(fmt.Sprintf("OpMode: %s, fetch: %t, serve: %t", *OpMode, Fetch, Serve))
}
