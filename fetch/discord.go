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
		Content:  fmt.Sprintf("<https://s.team/a/%d>", AppID),
		Username: "recommendations.steamsal.es",
	}, nil)
}

// PostDiscordPreapprove posts a pre-approval message to LU_POST_PREAPPROVAL
func PostDiscordPreapprove(AppID int, recommendation string) error {
	if os.Getenv("LU_POST_PREAPPROVAL") == "" {
		util.Warn("Not sending a Web hook because empty LU_POST_PREAPPROVAL")
		return nil
	}

	return httpPostJSON(os.Getenv("LU_POST_PREAPPROVAL"), &types.DiscordOutgoingWebhook{
		Content:  fmt.Sprintf("https://store.steampowered.com/app/%d\n```%s```", AppID, recommendation),
		Username: "recommendations.steamsal.es",
	}, nil)
}
