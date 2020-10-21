package fetch

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

func fetchAppPageInfo(appid int) (int, int, bool) {
	appPage, err := HTTPGet(fmt.Sprintf("https://store.steampowered.com/app/%d", appid))
	if err != nil {
		util.Warn(fmt.Sprintf("Error fetching reviews: %v", err))
		return 0, 0, false
	}

	pg := string(appPage)

	totalRegex, err := regexp.Compile("review_summary_num_reviews\\\" value=\\\"(\\d*)\\\"")
	if err != nil {
		util.Warn(fmt.Sprintf("Error fetching reviews: %v", err))
		return 0, 0, false
	}

	trRes := totalRegex.FindStringSubmatch(pg)
	if len(trRes) < 2 {
		util.Warn(fmt.Sprintf("Could not find review count"))
		return 0, 0, false
	}

	total, err := strconv.Atoi(trRes[1])
	if err != nil {
		total = 0
	}

	positiveRegex, err := regexp.Compile("review_summary_num_positive_reviews\\\" value=\\\"(\\d*)\\\"")
	if err != nil {
		util.Warn(fmt.Sprintf("Error fetching reviews: %v", err))
		return 0, 0, false
	}

	prRes := positiveRegex.FindStringSubmatch(pg)
	if len(prRes) < 2 {
		util.Warn(fmt.Sprintf("Could not find positive review count"))
		return 0, 0, false
	}

	positive, err := strconv.Atoi(prRes[1])
	if err != nil {
		positive = 0
	}

	demo := strings.Index(pg, "demo_above_purchase") > -1

	return total, positive, demo
}

// Steam fetches an app's information from Steam and will update the passed App
func Steam(app *types.App, cc string) error {
	util.Info(fmt.Sprintf("Fetching Steam info for app %d in %s", app.AppID, cc))
	sf := map[string]types.SteamStorefront{}

	if err := httpJSON(fmt.Sprintf("https://store.steampowered.com/api/appdetails/?appids=%d&cc=%s", app.AppID, cc), &sf); err != nil {
		return err
	}

	if !sf[strconv.Itoa(app.AppID)].Success {
		return fmt.Errorf("App lookup unsuccessful")
	}

	sfapp := sf[strconv.Itoa(app.AppID)].Data

	if strings.ToLower(cc) == "us" {
		app.Name = sfapp.Name
		app.Developers = sfapp.Developers
		app.Publishers = sfapp.Publishers
		app.Description = sfapp.ShortDescription
		app.Screenshot = sfapp.Screenshots[rand.Intn(len(sfapp.Screenshots))].Path

		genres := []string{}

		for _, g := range sfapp.Genres {
			genres = append(genres, g.Description)
		}

		for _, t := range AppInfo(app.AppID) {
			genres = append(genres, t)
		}

		app.Genres = genres

		app.Platforms.Windows = sfapp.Platforms.Windows
		app.Platforms.MacOS = sfapp.Platforms.Mac
		app.Platforms.Linux = sfapp.Platforms.Linux

		total, positive, demo := fetchAppPageInfo(app.AppID)

		app.ReviewsPositive = positive
		app.ReviewsTotal = total
		app.Demo = demo
		app.Score = util.RateWilson(float64(positive), float64(total))

		if math.IsNaN(app.Score) {
			app.Score = 0
		}

		util.Info("Fetched app page. Waiting a second...")
		time.Sleep(2 * time.Second)
	}

	if app.Prices == nil {
		app.Prices = &types.AppPrices{}
	}

	if app.Prices.Steam == nil {
		app.Prices.Steam = make(map[string]types.AppPrice)
	}

	app.Prices.Steam[strings.ToLower(cc)] = types.AppPrice{
		Price:         sfapp.Price.Final,
		Discount:      sfapp.Price.DiscountPercent,
		URL:           fmt.Sprintf("https://store.steampowered.com/app/%d", app.AppID),
		OriginalPrice: sfapp.Price.Initial,
	}

	return nil
}

// SteamAppList fetches the SteamAppList
func SteamAppList() (map[string]types.SteamGame, error) {
	isa := &types.SteamAppList{}

	err := httpJSON("https://api.steampowered.com/ISteamApps/GetAppList/v2/", &isa)
	if err != nil {
		return nil, err
	}

	list := make(map[string]types.SteamGame)
	for _, val := range isa.AppList.Apps {
		list[strconv.Itoa(val.AppID)] = val
	}

	return list, nil
}
