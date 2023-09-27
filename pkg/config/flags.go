package config

// Flags of the command line
type Flags struct {
	Config string `flag:"config" desc:"Configuration file to load" default:""`
}
