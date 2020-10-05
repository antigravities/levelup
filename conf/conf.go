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

	// ForceDiscord determines whether we should automatically post Discord messages
	ForceDiscord bool = false

	// ForceFetch determines whether we should re-fetch all apps now
	ForceFetch bool = false
)

// Init finds the command line flags and sets FetchOnly/ServeOnly
func Init() {
	OpMode = flag.String("mode", "all", "Operation mode. Possible options are fetch, serve, and all (which is both).")
	forceDiscord := flag.Bool("force-discord", false, "Force discord reposting?")
	forceFetch := flag.Bool("force-fetch", false, "Force re-fetch of all apps?")
	flag.Parse()

	if forceDiscord != nil {
		ForceDiscord = *forceDiscord
	}

	if forceFetch != nil {
		ForceFetch = *forceFetch
	}

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
