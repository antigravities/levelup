package fetch

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

var (
	fanSearchKey          string
	fanSearchKeyFetchTime int64
)

// Fanatical searches and updates apps from Fanatical
func Fanatical(app *types.App, cc string) error {
	util.Info(fmt.Sprintf("Fetching Fanatical info for app %d in %s", app.AppID, cc))

	if fanSearchKey == "" || time.Now().Unix()-fanSearchKeyFetchTime > 3600 {
		fanSearchKey = ""

		body, err := HTTPGet("https://www.fanatical.com/")
		if err != nil {
			return err
		}

		regex, err := regexp.Compile("window.searchKey = \\'(.*)\\';")
		if err != nil {
			return err
		}

		sk := regex.FindStringSubmatch(string(body))

		if len(sk) < 2 || sk[1] == "" {
			util.Warn("Could not find Fanatical search key")
			return err
		}

		fanSearchKey = strings.Split(sk[1], ("';"))[0]
		fanSearchKeyFetchTime = time.Now().Unix()

		util.Info("Found a new Fanatical search key")
	}

	res := &types.FanAlgoliaIncoming{}
	if err := httpPostJSON("https://w2m9492ddv-dsn.algolia.net/1/indexes/fan_alt_rank/query?x-algolia-agent=Algolia%20for%20JavaScript%20(3.35.1)%3B%20Browser%20(lite)&x-algolia-application-id=W2M9492DDV", types.FanAlgoliaOutgoing{
		Params: "query=" + url.QueryEscape(app.Name) + "&hitsPerPage=5",
		APIKey: fanSearchKey,
	}, res); err != nil {
		return err
	}

	for _, v := range res.Hits {
		target := ""

		if strings.ToLower(v.Name) == strings.ToLower(app.Name) {
			switch cc {
			case "us":
				target = "USD"
				break
			case "fr":
				target = "EUR"
				break
			case "uk":
				target = "GBP"
				break
			}

			if target == "" {
				util.Warn(fmt.Sprintf("Region not implemented in Fanatical: %v", cc))
				return nil
			}

			// probably a bad assumption. there's Price and FullPrice fields that I'll need to rely on existing
			if _, ok := v.Price[target]; ok {
				if app.Prices.Fanatical == nil {
					app.Prices.Fanatical = make(map[string]types.AppPrice)
				}

				app.Prices.Fanatical[cc] = types.AppPrice{
					OriginalPrice: uint64(v.FullPrice[target] * 100), // should be good as we're not dealing with insane amounts of digits
					Price:         uint64(v.Price[target] * 100),
					Discount:      v.DiscountPercent,
					URL:           "https://www.fanatical.com/en/game/" + v.Slug,
				}
			}

			return nil
		}
	}

	return nil
}
