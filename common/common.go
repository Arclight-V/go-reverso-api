package common

import (
	"os"
	"path/filepath"
)

func GetDefaultDataDir() string {
	dir, _ := os.Getwd()
	// TODO:: add getting directory and join with current directory
	//if err != nil {
	//	return "", err
	//}
	return filepath.Join(dir, DefaultDataDir)
}
