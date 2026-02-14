package cmd

// Creator function for a command
type Creator func() Cmd

// Commands registry
var Commands = map[string]Cmd{}

// Add command to the registry
func Add(name string, creator Creator) {
	Commands[name] = creator()
}
