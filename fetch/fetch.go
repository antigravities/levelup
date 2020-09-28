package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

var rqclient *http.Client = &http.Client{}

const userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36"

// HTTPGet fetches a document using HTTP and returns the bytes of the file if found, or an error, if any
func HTTPGet(url string) ([]byte, error) {
	//	resp, err := http.Get(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []byte{}, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := rqclient.Do(req)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		util.Warn(fmt.Sprintf("Error: %v", fmt.Errorf("status code %v not OK", resp.StatusCode)))
		return []byte{}, fmt.Errorf("status code %v not OK", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func httpJSON(url string, cast interface{}) error {
	resp, err := HTTPGet(url)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, cast); err != nil {
		return err
	}

	return nil
}

func httpPostJSON(url string, data interface{}, cast interface{}) error {
	dataStr, err := json.Marshal(data)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataStr))
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := rqclient.Do(req)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return err
	}

	err = json.Unmarshal(bytes, cast)
	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
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

	if err := Fanatical(app, cc); err != nil {
		return err
	}

	app.LastUpdate = time.Now().Unix()

	return nil
}

// AllRegions works the same as All, but fetches everything for all regions.
func AllRegions(app *types.App) error {
	for _, region := range types.Regions {
		if err := All(app, region); err != nil {
			util.Warn(fmt.Sprintf("Error: %v", err))
			return err
		}

		util.Debug("Waiting a moment...")
		// Can only make 200 requests per 5 minutes to the Storefront API
		time.Sleep(2000 * time.Millisecond)
	}

	return nil
}
