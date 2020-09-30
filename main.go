package main

import (
	"os"
	"time"

	"get.cutie.cafe/levelup/conf"
	db "get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/scheduled"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/util"
	"get.cutie.cafe/levelup/www"
)

func main() {
	util.LoadEnv()

	conf.Init()

	util.LogOpen()
	defer util.LogClose()

	util.Info("Level Up")
	util.Info("Copyright (c) 2020 Cutie Cafe")

	if !conf.Serve {
		util.Warn(" ========================================================================")
		util.Warn("| Running in fetch mode! This may not be what you're expecting.          |")
		util.Warn("| A web server will not be started. Ensure you're using a US IP address! |")
		util.Warn(" ========================================================================")
	}

	if !conf.Fetch {
		util.Warn(" ========================================================================================== ")
		util.Warn("| Running in serve mode! This may not be what you're expecting.                            |")
		util.Warn("| Price updates will not be fetched. Run another levelup in fetch/all mode to fetch prices.|")
		util.Warn(" ========================================================================================== ")
	}

	db.Initialize()

	if conf.Serve {
		if _, err := os.Stat("map"); os.IsNotExist(err) {
			search.Refresh()
		} else {
			util.Debug("Skipping initial search engine refresh (for now...)")
		}
	}

	for _, v := range db.GetApps(false) {
		db.GetApp(v)
	}

	scheduled.Start()

	if conf.Serve {
		www.Start()
	} else {
		for true {
			time.Sleep(3000)
		}
	}

	util.Info("Shutting down")
}
