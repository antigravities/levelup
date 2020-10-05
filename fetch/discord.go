package fetch

import (
	"fmt"
	"os"
	"time"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// PostDiscord submits an app to LU_POST_APPROVAL
func PostDiscord(AppID int) error {
	if os.Getenv("LU_POST_APPROVAL") == "" || os.Getenv("LU_WEBROOT") == "" {
		util.Warn("Not sending a Web hook because empty LU_POST_APPROVAL or LU_WEBROOT")
		return nil
	}

	if err := httpPostJSON(os.Getenv("LU_POST_APPROVAL"), &types.DiscordOutgoingWebhook{
		Content:  fmt.Sprintf("%sapi/image/%d.png", os.Getenv("LU_WEBROOT"), AppID),
		Username: "recommendations.steamsal.es",
	}, nil); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	return httpPostJSON(os.Getenv("LU_POST_APPROVAL"), &types.DiscordOutgoingWebhook{
		Content:  fmt.Sprintf("<https://store.steampowered.com/app/%d>", AppID),
		Username: "recommendations.steamsal.es",
	}, nil)
}
