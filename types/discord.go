package types

// DiscordOutgoingWebhook represents an object we execute on a Discord web hook
type DiscordOutgoingWebhook struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}
