package main

import (
	"os"

	db "get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/scheduled"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/util"
	"get.cutie.cafe/levelup/www"
)

func main() {
	util.LoadEnv()

	util.LogOpen()
	defer util.LogClose()

	util.Info("Level Up")
	util.Info("Copyright (c) 2020 Cutie Cafe")

	db.Initialize()

	if _, err := os.Stat("map"); os.IsNotExist(err) {
		search.Refresh()
	} else {
		util.Debug("Skipping initial search engine refresh (for now...)")
	}

	for _, v := range db.GetApps(false) {
		db.GetApp(v)
	}

	scheduled.Start()

	www.Start()

	util.Info("Shutting down")
}
