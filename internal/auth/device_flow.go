package auth

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ravi-technologies/sunday-cli/internal/api"
	"github.com/ravi-technologies/sunday-cli/internal/config"
	"github.com/ravi-technologies/sunday-cli/internal/output"
)

const (
	// DefaultSpinnerCharSet is the Braille spinner pattern (index 14 in yacspin).
	DefaultSpinnerCharSet = 14
)

// DeviceFlow handles the device code authentication flow
type DeviceFlow struct {
	client  *api.Client
	spinner *spinner.Spinner
}

// NewDeviceFlow creates a new device flow handler
func NewDeviceFlow() (*DeviceFlow, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	s := spinner.New(spinner.CharSets[DefaultSpinnerCharSet], 100*time.Millisecond)
	s.Suffix = " Waiting for authorization..."

	return &DeviceFlow{
		client:  client,
		spinner: s,
	}, nil
}

// Run executes the device code flow
func (d *DeviceFlow) Run() error {
	// Request device code
	codeResp, err := d.client.RequestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}

	// Display instructions
	fmt.Println()
	fmt.Println("To authenticate, visit:")
	fmt.Printf("  %s\n", codeResp.VerificationURI)
	fmt.Println()
	fmt.Println("And enter the code:")
	fmt.Printf("  %s\n", codeResp.UserCode)
	fmt.Println()

	// Try to open browser
	if err := openBrowser(codeResp.VerificationURI + "?user_code=" + codeResp.UserCode); err != nil {
		// Not a fatal error, user can manually visit URL
		fmt.Println("(Could not open browser automatically)")
	}

	// Start polling with spinner
	d.spinner.Start()
	defer d.spinner.Stop()

	interval := time.Duration(codeResp.Interval) * time.Second
	deadline := time.Now().Add(time.Duration(codeResp.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		tokenResp, errCode, err := d.client.PollForToken(codeResp.DeviceCode)
		if err != nil {
			return fmt.Errorf("polling error: %w", err)
		}

		// Check error codes
		switch errCode {
		case "authorization_pending":
			// Still waiting, continue polling
			time.Sleep(interval)
			continue
		case "expired_token":
			return fmt.Errorf("device code expired. Please try again")
		case "":
			// Success! Save tokens
			d.spinner.Stop()

			cfg := &config.Config{
				AccessToken:  tokenResp.Access,
				RefreshToken: tokenResp.Refresh,
				ExpiresAt:    time.Now().Add(api.TokenExpiryBuffer), // Assume ~5 min expiry
				UserEmail:    tokenResp.User.Email,
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save tokens: %w", err)
			}

			output.Current.PrintMessage(fmt.Sprintf("Authenticated as %s", tokenResp.User.Email))
			return nil
		default:
			return fmt.Errorf("authentication error: %s", errCode)
		}
	}

	return fmt.Errorf("authentication timed out")
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}
