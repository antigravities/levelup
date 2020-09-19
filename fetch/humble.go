package fetch

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// Humble fetches an app's prices from Humble and will update the passed App
func Humble(app *types.App, cc string) error {
	util.Info(fmt.Sprintf("Fetching Humble info for app %d", app.AppID))

	if app.Name == "" {
		return fmt.Errorf("app %d has no name", app.AppID)
	}

	q := "https://www.humblebundle.com/store/api/search?request=0&page=0&sort=bestselling&filter=all&search=" + url.QueryEscape(app.Name)

	if len(app.Developers) > 0 {
		q = q + "&developer=" + url.QueryEscape(app.Developers[0])
	}

	hs := types.HumbleSearch{}
	if err := httpJSON(q, &hs); err != nil {
		return err
	}

	if len(hs.Results) == 0 {
		return nil
	}

	for _, val := range hs.Results {
		if strings.ToLower(val.Name) == strings.ToLower(app.Name) {
			if app.Prices.Humble == nil {
				app.Prices.Humble = make(map[string]types.AppPrice)
			}

			app.Prices.Humble[cc] = types.AppPrice{
				Price:         uint64(math.Ceil(val.CurrentPrice.Amount * 100)),
				OriginalPrice: uint64(math.Ceil(val.FullPrice.Amount * 100)),
				Discount:      uint(math.Floor((1 - (val.CurrentPrice.Amount / val.FullPrice.Amount)) * 100)),
				URL:           "https://humble.com/store/" + val.HumanURL,
			}

			return nil
		}
	}

	return nil
}
