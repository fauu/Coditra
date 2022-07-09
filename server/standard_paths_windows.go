//go:build windows

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
	if envHome := os.Getenv("APPDATA"); envHome != "" {
		return filepath.Join(envHome, appName), nil
	}
	return "", errors.New("getting user appdata directory")
}
