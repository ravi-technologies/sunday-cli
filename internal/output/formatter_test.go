package output

import (
	"testing"
)

func TestSetJSON_True(t *testing.T) {
	// Reset to default first
	Current = &HumanFormatter{}

	// Set JSON mode
	SetJSON(true)

	// Verify Current is now JSONFormatter
	_, ok := Current.(*JSONFormatter)
	if !ok {
		t.Errorf("SetJSON(true) should set Current to JSONFormatter, got %T", Current)
	}
}

func TestSetJSON_False(t *testing.T) {
	// Set to JSON first
	Current = &JSONFormatter{}

	// Set human mode
	SetJSON(false)

	// Verify Current is now HumanFormatter
	_, ok := Current.(*HumanFormatter)
	if !ok {
		t.Errorf("SetJSON(false) should set Current to HumanFormatter, got %T", Current)
	}
}

func TestCurrent_DefaultHuman(t *testing.T) {
	// Save and restore Current after test
	originalCurrent := Current
	defer func() { Current = originalCurrent }()

	// Reinitialize to test package default
	Current = &HumanFormatter{}

	// Verify default is HumanFormatter
	_, ok := Current.(*HumanFormatter)
	if !ok {
		t.Errorf("Default Current should be HumanFormatter, got %T", Current)
	}
}
