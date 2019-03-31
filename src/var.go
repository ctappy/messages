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

type args struct {
	// debug flags
	Info  bool
	Warn  bool
	Debug bool
	Trace bool
	// Bot info
	BotID   string
	BotName string
	// Bot options
	BotDisable bool
}

// default values
func DefaultArgs() args {
	return args{
		false, // Info  bool
		false, // Warn  bool
		false, // Debug bool
		false, // Trace bool
		// Bot info
		"", // BotID   string
		"", // BotName string
		// Bot options
		false, // BotDisable bool
	}
}
