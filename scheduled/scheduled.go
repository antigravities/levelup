package scheduled

import (
	"fmt"

	"get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/fetch"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/util"
	"github.com/carlescere/scheduler"
)

// Start initializes the scheduler functions
func Start() {
	util.Info("Initializing scheduled tasks")

	scheduler.Every(6).Hours().Run(func() {
		search.Refresh()
	})

	scheduler.Every(30).Minutes().Run(RefreshStaleApps)
}

// RefreshStaleApps refreshes all of the apps that are stale (LastUpdated > 1 hour ago).
func RefreshStaleApps() {
	apps := dynamodb.FindStaleApps()

	for _, app := range apps {
		if err := fetch.AllRegions(&app); err != nil {
			util.Warn("Hit an error, backing off for now!")
			util.Warn(fmt.Sprintf("%v", err))
			return
		}

		dynamodb.PutApp(app)
	}
}
