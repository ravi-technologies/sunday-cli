// Package output provides formatters for displaying data to users.
//
// Two formatters are available:
//   - HumanFormatter: Produces human-readable, colored terminal output
//   - JSONFormatter: Produces machine-parseable JSON output
//
// The active formatter is controlled by the --json flag and can be
// switched using SetJSON(). All commands should use the Current()
// formatter to respect the user's output preference.
//
// Example:
//
//	output.SetJSON(jsonFlag)
//	output.Current().Print(data)
package output
