// Package output provides formatters for CLI output in human-readable and JSON formats.
package output

// Formatter defines the interface for outputting data in different formats.
type Formatter interface {
	// Print outputs data to stdout
	Print(data interface{}) error
	// PrintError outputs an error message
	PrintError(err error)
	// PrintMessage outputs a simple message
	PrintMessage(msg string)
	// PrintTable outputs tabular data with headers
	PrintTable(headers []string, rows [][]string)
}

// Current is the global formatter, set based on --json flag.
var Current Formatter = &HumanFormatter{}

// SetJSON switches between JSON and human-readable output modes.
func SetJSON(useJSON bool) {
	if useJSON {
		Current = &JSONFormatter{}
	} else {
		Current = &HumanFormatter{}
	}
}
