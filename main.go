/**
 * Copyright (c) 2020 Alexandra Frock, Cutie Café, contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"time"

	"get.cutie.cafe/levelup/conf"
	db "get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/scheduled"
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

	/*
		for _, v := range db.GetApps(false) {
			db.GetApp(v)
		}
	*/

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
