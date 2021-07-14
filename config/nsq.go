package config

var (
	AddrNsqlookupd string
	AddrNsqd       string
)

const (
	TopicName           = "chat"
	MaxInFlight         = 10
	LookupdPollInterval = 30
	ArchiveChannelName  = "archive"
	BotChannelName      = "bot"
)

func init() {
	AddrNsqlookupd = "0.0.0.0:4161"
	AddrNsqd = "0.0.0.0:4150"
}
