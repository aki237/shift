package main

import "os"

// IsFileExist returs a boolean whether the passed filename does exist and is a file.
func IsFileExist(filename string) bool {
	stat, err := os.Stat(filename)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}
