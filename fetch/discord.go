package fetch

import (
	"fmt"
	"os"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"
)

// PostDiscord submits an app to LU_POST_APPROVAL
func PostDiscord(AppID int) error {
	if os.Getenv("LU_POST_APPROVAL") == "" || os.Getenv("LU_WEBROOT") == "" {
		util.Warn("Not sending a Web hook because empty LU_POST_APPROVAL or LU_WEBROOT")
		return nil
	}

	return httpPostJSON(os.Getenv("LU_POST_APPROVAL"), &types.DiscordOutgoingWebhook{
		Content:  fmt.Sprintf("%s/api/image/%d.png", os.Getenv("LU_WEBROOT"), AppID),
		Username: "recommendations.steamsal.es",
	}, nil)
}
