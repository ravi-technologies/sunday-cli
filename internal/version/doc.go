// Package version provides build-time version information and configuration.
//
// Variables in this package are set at build time using -ldflags:
//
//	go build -ldflags "-X internal/version.Version=1.0.0 -X internal/version.APIBaseURL=https://api.example.com"
//
// The package provides:
//   - Version: The application version string
//   - APIBaseURL: The backend API base URL (required)
//   - GetAPIBaseURL(): Returns the API URL or error if not set
//   - GetVersion(): Returns version or "dev" if not set
package version
