package dht

// Root Config
var (
	DefaultLogLevel  = "info"
	DefaultLogFormat = "text"

	CustomUserAgent = "probelab-dht-crawler"
)

type RootConfig struct {
	LogLevel  string
	LogFormat string
}

// Lookup Config
var (
	DefaultNetwork     = Mainnet
	DefaultIsNamespace = false
	DefaultNamespace   = NsFull
)

type LookupCmdConfig struct {
	Network string

	IsCustomNamespace bool
	Namespace         string
}
