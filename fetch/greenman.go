package fetch

import (
	"fmt"
	"net/url"
	"strings"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// GreenMan fetches the Green Man Gaming information for an app
func GreenMan(app *types.App, cc string) error {
	util.Info(fmt.Sprintf("Fetching Green Man Gaming info for app %d", app.AppID))

	query := &types.GreenManSearch{}

	if err := httpJSON(fmt.Sprintf("https://api.greenmangaming.com/api/v2/quick_search_results/%s", url.QueryEscape(app.Name)), query); err != nil {
		if strings.Index(err.Error(), "status code 404") < 0 {
			return err
		}

		util.Warn(fmt.Sprintf("Could not find Green Man Gaming page for %s", app.Name))
		return nil
	}

	next := ""
	for _, v := range query.Results {
		if strings.ToLower(v.Name) == strings.ToLower(app.Name) {
			next = v.URL
			break
		}
	}

	if next == "" {
		util.Warn(fmt.Sprintf("Could not find Green Man Gaming page for %s", app.Name))
		return nil
	}

	util.Info(next)
	return nil
}
