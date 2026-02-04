package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdoutJSON captures stdout output from a function.
func captureStdoutJSON(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// captureStderrJSON captures stderr output from a function.
func captureStderrJSON(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestJSONFormatter_Print(t *testing.T) {
	formatter := &JSONFormatter{}

	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Count int    `json:"count"`
	}

	data := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Count: 42,
	}

	output := captureStdoutJSON(func() {
		err := formatter.Print(data)
		if err != nil {
			t.Errorf("Print() returned error: %v", err)
		}
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("Print() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify contents
	var result TestStruct
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal Print() output: %v", err)
	}

	if result.Name != "John Doe" {
		t.Errorf("Print() name = %q, want %q", result.Name, "John Doe")
	}
	if result.Email != "john@example.com" {
		t.Errorf("Print() email = %q, want %q", result.Email, "john@example.com")
	}
	if result.Count != 42 {
		t.Errorf("Print() count = %d, want %d", result.Count, 42)
	}
}

func TestJSONFormatter_PrintTable(t *testing.T) {
	formatter := &JSONFormatter{}

	headers := []string{"ID", "Name", "Status"}
	rows := [][]string{
		{"1", "Alice", "Active"},
		{"2", "Bob", "Inactive"},
	}

	output := captureStdoutJSON(func() {
		formatter.PrintTable(headers, rows)
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("PrintTable() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify structure
	var result TableOutput
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal PrintTable() output: %v", err)
	}

	// Verify headers
	if len(result.Headers) != 3 {
		t.Errorf("PrintTable() headers length = %d, want %d", len(result.Headers), 3)
	}
	expectedHeaders := []string{"ID", "Name", "Status"}
	for i, h := range expectedHeaders {
		if result.Headers[i] != h {
			t.Errorf("PrintTable() header[%d] = %q, want %q", i, result.Headers[i], h)
		}
	}

	// Verify rows
	if len(result.Rows) != 2 {
		t.Errorf("PrintTable() rows length = %d, want %d", len(result.Rows), 2)
	}
	if result.Rows[0][1] != "Alice" {
		t.Errorf("PrintTable() rows[0][1] = %q, want %q", result.Rows[0][1], "Alice")
	}
	if result.Rows[1][1] != "Bob" {
		t.Errorf("PrintTable() rows[1][1] = %q, want %q", result.Rows[1][1], "Bob")
	}
}

func TestJSONFormatter_PrintError(t *testing.T) {
	formatter := &JSONFormatter{}
	testErr := errors.New("something went wrong")

	output := captureStderrJSON(func() {
		formatter.PrintError(testErr)
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("PrintError() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify structure
	var result map[string]string
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal PrintError() output: %v", err)
	}

	// Verify error field
	if result["error"] != "something went wrong" {
		t.Errorf("PrintError() error = %q, want %q", result["error"], "something went wrong")
	}
}

func TestJSONFormatter_PrintMessage(t *testing.T) {
	formatter := &JSONFormatter{}
	msg := "Operation completed successfully"

	output := captureStdoutJSON(func() {
		formatter.PrintMessage(msg)
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("PrintMessage() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify structure
	var result map[string]string
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal PrintMessage() output: %v", err)
	}

	// Verify message field
	if result["message"] != msg {
		t.Errorf("PrintMessage() message = %q, want %q", result["message"], msg)
	}
}

func TestJSONFormatter_Print_InvalidInput(t *testing.T) {
	formatter := &JSONFormatter{}

	// Create an unmarshalable type (channel)
	ch := make(chan int)

	var printErr error
	output := captureStdoutJSON(func() {
		printErr = formatter.Print(ch)
	})

	// Verify error is returned for unmarshalable input
	if printErr == nil {
		t.Error("Print() should return error for unmarshalable input")
	}

	// Verify error message mentions JSON
	if !strings.Contains(printErr.Error(), "JSON") {
		t.Errorf("Print() error should mention JSON, got: %v", printErr)
	}

	// Output should be empty since it failed
	if strings.TrimSpace(output) != "" {
		t.Errorf("Print() should not output anything on error, got: %s", output)
	}
}

func TestJSONFormatter_Print_Nil(t *testing.T) {
	formatter := &JSONFormatter{}

	output := captureStdoutJSON(func() {
		err := formatter.Print(nil)
		if err != nil {
			t.Errorf("Print(nil) returned error: %v", err)
		}
	})

	// nil should marshal to "null" in JSON
	trimmed := strings.TrimSpace(output)
	if trimmed != "null" {
		t.Errorf("Print(nil) = %q, want %q", trimmed, "null")
	}
}

func TestJSONFormatter_Print_Slice(t *testing.T) {
	formatter := &JSONFormatter{}

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := []Item{
		{ID: 1, Name: "First"},
		{ID: 2, Name: "Second"},
	}

	output := captureStdoutJSON(func() {
		err := formatter.Print(data)
		if err != nil {
			t.Errorf("Print() returned error: %v", err)
		}
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("Print() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify contents
	var result []Item
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal Print() output: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Print() slice length = %d, want %d", len(result), 2)
	}
	if result[0].Name != "First" {
		t.Errorf("Print() result[0].Name = %q, want %q", result[0].Name, "First")
	}
}

func TestJSONFormatter_Print_Map(t *testing.T) {
	formatter := &JSONFormatter{}

	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	output := captureStdoutJSON(func() {
		err := formatter.Print(data)
		if err != nil {
			t.Errorf("Print() returned error: %v", err)
		}
	})

	// Verify output is valid JSON
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		t.Errorf("Print() output is not valid JSON: %s", output)
	}

	// Unmarshal and verify contents
	var result map[string]interface{}
	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal Print() output: %v", err)
	}

	if result["key1"] != "value1" {
		t.Errorf("Print() key1 = %v, want %v", result["key1"], "value1")
	}
}

func TestJSONFormatter_Print_Indented(t *testing.T) {
	formatter := &JSONFormatter{}

	data := map[string]string{"key": "value"}

	output := captureStdoutJSON(func() {
		err := formatter.Print(data)
		if err != nil {
			t.Errorf("Print() returned error: %v", err)
		}
	})

	// Verify output is indented (contains newlines and spaces)
	if !strings.Contains(output, "\n") {
		t.Errorf("Print() output should be indented with newlines, got: %s", output)
	}
	if !strings.Contains(output, "  ") {
		t.Errorf("Print() output should be indented with spaces, got: %s", output)
	}
}
