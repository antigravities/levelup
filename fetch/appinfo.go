package fetch

import (
	"fmt"
	"os"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// AppInfo fetches the tags of an app (generically named because we may fetch more from AppInfo soon)
func AppInfo(appid int) []string {
	if os.Getenv("LU_APPINFO") == "" {
		util.Warn("LU_APPINFO is not defined, skipping fetch")
		return []string{}
	}

	appinfo := &types.AppInfo{}
	err := httpJSON(fmt.Sprintf("%s%d", os.Getenv("LU_APPINFO"), appid), appinfo)

	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []string{}
	}

	if appinfo.Error != "" {
		util.Warn(fmt.Sprintf("MicroAppInfo error: %s", appinfo.Error))
		return []string{}
	}

	return appinfo.Common.StoreTags
}
