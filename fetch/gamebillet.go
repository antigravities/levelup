package fetch

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"

	"github.com/PuerkitoBio/goquery"
)

// GameBillet fetches app info from GameBillet
func GameBillet(app *types.App, cc string) error {
	util.Info(fmt.Sprintf("Fetching GameBillet info for app %d", app.AppID))

	results := &[]types.GameBilletSearch{}

	if err := httpJSON(fmt.Sprintf("https://gamebillet.com/catalog/searchtermautocomplete?term=%s", url.QueryEscape(app.Name)), results); err != nil {
		util.Warn(fmt.Sprintf("Suppressed error: %v", err))
		return nil
	}

	next := ""

	for _, result := range *results {
		if strings.TrimSpace(strings.ToLower(result.Label)) == strings.TrimSpace(strings.ToLower(app.Name)) {
			next = result.ProductURL
			break
		}
	}

	if next == "" {
		util.Warn(fmt.Sprintf("Couldn't find GameBillet page for %s", app.Name))
		return nil
	}

	util.Debug(next)

	body, err := HTTPGet(fmt.Sprintf("https://gamebillet.com%s", next))
	if err != nil {
		return err
	}

	j, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	btn := j.Find("[id^='add-to-cart-button']")
	if btn == nil {
		util.Warn(fmt.Sprintf("Could not find add to cart button on %s. Setting price to 99999...", app.Name))
		app.Prices.GameBillet[cc] = types.AppPrice{
			Price:         99999,
			OriginalPrice: 99999,
			Discount:      0,
			URL:           fmt.Sprintf("https://gamebillet.com%s", next),
		}
		return nil
	}

	discount := uint64(0)
	price := uint64(99999)
	originalPrice := uint64(99999)

	ins := btn.Find("ins")
	if ins != nil {
		text := ins.Text()
		discount, err = strconv.ParseUint(text[1:len(text)-1], 10, 32)
		if err != nil {
			util.Warn(fmt.Sprintf("Could not parse discount percent on %s: %v", app.Name, err))
			discount = 0
		}
	}

	span := btn.Parent().Find("span")
	if span != nil {
		text := span.Text()
		fPrice, err := strconv.ParseFloat(text[1:], 32)
		if err != nil {
			util.Warn(fmt.Sprintf("Could not parse price on %s.", app.Name))
			discount = 0
		} else {
			price = uint64(fPrice * 100)
			originalPrice = price
		}
	}

	sp := btn.Parent().Parent().Find(".buy-timer-sale")
	if sp != nil {
		text := span.Text()
		fPrice, err := strconv.ParseFloat(text[1:], 32)
		if err != nil {
			util.Warn(fmt.Sprintf("Could not parse original price on %s.", app.Name))
		} else {
			originalPrice = uint64(fPrice * 100)
		}
	}

	if app.Prices.GameBillet == nil {
		app.Prices.GameBillet = make(map[string]types.AppPrice)
	}

	app.Prices.GameBillet[cc] = types.AppPrice{
		Price:         price,
		OriginalPrice: originalPrice,
		Discount:      uint(discount),
		URL:           fmt.Sprintf("https://gamebillet.com%s", next),
	}

	return nil
}
