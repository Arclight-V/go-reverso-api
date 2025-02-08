package config

import (
	"os"
	"path/filepath"
)

func DefaultFilePath() string {
	dir, _ := os.Getwd()
	// TODO:: add getting directory and join with current directory
	//if err != nil {
	//	return "", err
	//}
	return filepath.Join(dir, File)
}
