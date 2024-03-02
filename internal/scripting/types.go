package scripting

// ScriptInput is the paylod passed to the script
type ScriptInput struct {
	URL    string
	HTML   string
	Text   string
	Title  string
	Origin string
	Tags   []string
}
