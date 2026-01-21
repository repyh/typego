package ecosystem

import (
	"os"
	"path/filepath"
)

const (
	HiddenDirName = ".typego"
	BinaryName    = "typego-app.exe"
	HandoffEnvVar = "TYPEGO_HANDOFF"
)

// GetJITBinaryPath returns the path to the JIT binary if it exists
func GetJITBinaryPath(cwd string) (string, bool) {
	path := filepath.Join(cwd, HiddenDirName, "bin", BinaryName)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

// IsHandoffRequired returns true if we should delegate to the JIT binary
func IsHandoffRequired(cwd string) bool {
	// Prevent infinite loops where the JIT binary calls itself
	if os.Getenv(HandoffEnvVar) == "true" {
		return false
	}

	_, exists := GetJITBinaryPath(cwd)
	return exists
}
