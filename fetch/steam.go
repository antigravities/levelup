package fetch

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

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

		app.Genres = genres
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
