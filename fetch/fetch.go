package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"get.cutie.cafe/levelup/types"
)

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status code %v not OK", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func httpJSON(url string, cast interface{}) error {
	resp, err := httpGet(url)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, cast); err != nil {
		return err
	}

	return nil
}

// All fetches all of the app's prices and info
func All(app *types.App, cc string) error {
	if err := Steam(app, cc); err != nil {
		return err
	}

	if cc == "us" {
		if err := Humble(app, cc); err != nil {
			return err
		}
	}

	app.LastUpdate = time.Now().Unix()

	return nil
}

// AllRegions works the same as All, but fetches everything for all regions.
func AllRegions(app *types.App) error {
	for _, region := range types.Regions {
		if err := All(app, region); err != nil {
			return err
		}
	}

	return nil
}
