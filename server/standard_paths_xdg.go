//go:build freebsd || linux || netbsd || openbsd || solaris

package coditra

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetUserConfigLocation(appName string) (string, error) {
	if appName == "" {
		return "", errors.New("'appName' is empty")
	}
	if envHome := os.Getenv("XDG_CONFIG_HOME"); envHome != "" {
		return filepath.Join(envHome, appName), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", appName), nil
}
