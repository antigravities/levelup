package search

import (
	"fmt"
	"os"
	"strconv"

	"get.cutie.cafe/levelup/fetch"
	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
)

var (
	index bleve.Index
	apps  map[string]types.SteamGame = make(map[string]types.SteamGame)
)

// Refresh the search index.
func Refresh() error {
	util.Info("Fetching Steam app list")

	wapps, err := fetch.SteamAppList()
	if err != nil {
		return err
	}

	if _, err := os.Stat("map"); !os.IsNotExist(err) {
		util.Debug("Using existing index")

		if index == nil {
			index, err = bleve.Open("map")
			if err != nil {
				util.Debug(fmt.Sprintf("Error: %v", err))
				return err
			}
		}
	} else {
		util.Debug("Using new index")

		if index == nil {
			index, err = bleve.NewUsing("map", bleve.NewIndexMapping(), scorch.Name, scorch.Name, nil)
			if err != nil {
				util.Debug(fmt.Sprintf("Error: %v", err))
				return err
			}
		}
	}

	util.Info("Indexing " + strconv.Itoa(len(wapps)) + " apps")

	util.Debug("Started batch")
	batch := index.NewBatch()
	for _, val := range wapps {
		batch.Index(strconv.Itoa(val.AppID), val)
		apps[strconv.Itoa(val.AppID)] = val
	}
	util.Debug("Ended batch")

	util.Debug("Running batch operations")
	index.Batch(batch)

	util.Info("Done indexing")
	return err
}

// Query performs a query against the index.
func Query(qs string) ([]types.SteamGame, error) {
	q := bleve.NewMatchQuery(qs)
	q.SetFuzziness(0)

	result, err := index.Search(bleve.NewSearchRequest(q))
	if err != nil {
		return nil, err
	}

	res := make([]types.SteamGame, len(result.Hits))
	for i, hit := range result.Hits {
		res[i] = apps[hit.ID]
	}

	return res, nil
}

// IsApp checks if an AppID is in the index.
func IsApp(appid int) bool {
	_, ok := apps[strconv.Itoa(appid)]
	return ok
}
