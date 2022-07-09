//go:build darwin

package coditra

import (
	"errors"
	"os"
	"path/filepath"
)

func GetUserConfigLocation(appName string) (string, error) {
	if appName == "" {
		return "", errors.New("'appName' is empty")
	}
	if envHome := os.Getenv("HOME"); envHome != "" {
		return filepath.Join(envHome, "Library", "Application Support", appName), nil
	}
	return "", errors.New("getting user appdata directory")
}
