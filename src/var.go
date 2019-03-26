package variables

type Config struct {
	SMTP struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	Slack struct {
		SlackKey  string `json:"slack_key"`
		ChannelID string `json:"channel_id"`
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
func New(name string) Arguments {
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
