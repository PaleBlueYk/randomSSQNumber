package utils

import "os"

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return true
}
