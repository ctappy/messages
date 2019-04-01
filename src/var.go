package configuration

type Config struct {
	SMTP struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	Slack struct {
		BotUserToken string `json:"bot_user_token"`
		ChannelName  string `json:"channel_name"`
	} `json:"slack"`
}

type BotInfo struct {
	// Bot info
	ID          string
	Name        string
	ChannelName string
	ChannelID   string
	// Bot options
	Disable bool
}

// default values
func DefaultArgs() BotInfo {
	return BotInfo{
		// Bot info
		"", // BotID   string
		"", // BotName string
		"", // ChannelName string
		"", // ChannelID string
		// Bot options
		false, // BotDisable bool
	}
}
