package cmd

// Creator function for an output
type Creator func() Cmd

// Commands registry
var Commands = map[string]Cmd{}

// Add output to the registry
func Add(name string, creator Creator) {
	Commands[name] = creator()
}
