package main

import (
	db "get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/scheduled"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/util"
	"get.cutie.cafe/levelup/www"
)

func main() {
	util.Info("Level Up")
	util.Info("Copyright (c) 2020 Cutie Cafe")

	db.Initialize()

	search.Refresh()

	for _, v := range db.GetApps() {
		db.GetApp(v)
	}

	scheduled.Start()

	www.Start()
}
