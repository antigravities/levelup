package fetch

import (
	"fmt"
	"os"
	"strconv"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// AppInfo fetches the tags and review score of an app (generically named because we may fetch more from AppInfo soon)
func AppInfo(appid int) ([]string, int) {
	if os.Getenv("LU_APPINFO") == "" {
		util.Warn("LU_APPINFO is not defined, skipping fetch")
		return []string{}, 0
	}

	appinfo := &types.AppInfo{}
	err := httpJSON(fmt.Sprintf("%s%d", os.Getenv("LU_APPINFO"), appid), appinfo)

	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []string{}, 0
	}

	if appinfo.Error != "" {
		util.Warn(fmt.Sprintf("MicroAppInfo error: %s", appinfo.Error))
		return []string{}, 0
	}

	reviews, err := strconv.Atoi(appinfo.Common.ReviewPercentage)
	if err != nil {
		util.Warn(fmt.Sprintf("Could not marshal %s to string", appinfo.Common.ReviewPercentage))
		reviews = 0
	}

	return appinfo.Common.StoreTags, reviews
}
