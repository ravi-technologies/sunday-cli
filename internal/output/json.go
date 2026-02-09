package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// JSONFormatter outputs data in JSON format.
type JSONFormatter struct{}

// Print marshals data to indented JSON and outputs to stdout.
func (f *JSONFormatter) Print(data interface{}) error {
	output, err := marshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// PrintError outputs an error as JSON to stderr.
func (f *JSONFormatter) PrintError(err error) {
	output := map[string]string{
		"error": err.Error(),
	}
	data, marshalErr := marshalJSON(output)
	if marshalErr != nil {
		log.Printf("failed to marshal error JSON: %v", marshalErr)
		return
	}
	fmt.Fprintln(os.Stderr, string(data))
}

// PrintMessage outputs a message as JSON to stdout.
func (f *JSONFormatter) PrintMessage(msg string) {
	output := map[string]string{
		"message": msg,
	}
	data, marshalErr := marshalJSON(output)
	if marshalErr != nil {
		log.Printf("failed to marshal message JSON: %v", marshalErr)
		return
	}
	fmt.Println(string(data))
}

// TableOutput represents the JSON structure for table data.
type TableOutput struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

// PrintTable outputs tabular data as JSON to stdout.
func (f *JSONFormatter) PrintTable(headers []string, rows [][]string) {
	output := TableOutput{
		Headers: headers,
		Rows:    rows,
	}
	data, marshalErr := marshalJSON(output)
	if marshalErr != nil {
		log.Printf("failed to marshal table JSON: %v", marshalErr)
		return
	}
	fmt.Println(string(data))
}

// marshalJSON encodes data as indented JSON without escaping HTML characters.
// Go's json.Marshal escapes <, >, & as \u003c, \u003e, \u0026 by default,
// which breaks round-tripping of values like email Message-ID thread IDs.
func marshalJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// Encode appends a trailing newline; trim it since callers use Println
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}
