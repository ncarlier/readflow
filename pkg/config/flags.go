package config

// Flags of the command line
type Flags struct {
	Config    string `flag:"config" desc:"Configuration file to load" default:""`
	LogLevel  string `flag:"log-level" desc:"Log level (debug, info, warn, error)" default:"info"`
	LogPretty bool   `flag:"log-pretty" desc:"Output human readable logs" default:"false"`
}
