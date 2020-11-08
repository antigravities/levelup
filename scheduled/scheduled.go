package scheduled

import (
	"fmt"

	"get.cutie.cafe/levelup/conf"
	"get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/fetch"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
	"get.cutie.cafe/levelup/www"
	"github.com/carlescere/scheduler"
)

// Start initializes the scheduler functions
func Start() {
	util.Info("Initializing scheduled tasks")

	scheduler.Every(6).Hours().Run(func() {
		search.Refresh()
	})

	scheduler.Every(30).Minutes().Run(RefreshStaleApps)

	scheduler.Every(6).Hours().Run(www.ResetHelpfulRateLimit)
}

// RefreshStaleApps refreshes all of the apps that are stale (LastUpdated > 1 hour ago).
func RefreshStaleApps() {
	if conf.Fetch {
		var apps []types.App

		if !conf.ForceFetch {
			apps = dynamodb.FindStaleApps()
		} else {
			pApps := dynamodb.GetFullApps(false)

			for _, app := range pApps {
				apps = append(apps, *app)
			}
		}

		defer func() {
			conf.ForceFetch = false
			conf.ForceDiscord = false
		}()

		for _, app := range apps {
			shouldWebhook := app.LastUpdate == 0

			if err := fetch.AllRegions(&app); err != nil {
				util.Warn("Hit an error, backing off for now!")
				util.Warn(fmt.Sprintf("%v", err))
				return
			}

			dynamodb.PutApp(app)

			if shouldWebhook || conf.ForceDiscord {
				if err := fetch.PostDiscord(app.AppID); err != nil {
					util.Warn(fmt.Sprintf("Error: %v", err))
				}
			}
		}
	} else {
		util.Warn("Skipping stale app refresh, we're only serving.")
	}
}
