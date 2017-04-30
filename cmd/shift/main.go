package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// homeDir gets the home directory from environement
var homeDir = os.Getenv("HOME")

// Build Variables
var appVersion string = "v0.00"
var appBuildDate string
var appBuildPlatform string

// usage function is used to print the usage like printing the commandline options
func usage() {
	fmt.Println("shift : proxy manager for proGY,", appVersion, appBuildPlatform, appBuildDate)
	fmt.Println("Usage :", os.Args[0], "<options>")
	flag.PrintDefaults()
}

// main function is the entry
func main() {
	// progyConfigFile is the commandline option pointing to .progy config file
	var progyConfigFile = flag.String("config", filepath.Join(homeDir, ".progy"), "Path pointing to the .progy config file.")

	// proxListFile is the commandline option pointing to proxy list file in Toml Format
	var proxyListFile = flag.String("proxylist", filepath.Join(homeDir, ".config", "proxy"), "Path pointing to the proxy list file. (TOML only accepted)")

	flag.Usage = usage
	var version = flag.Bool("version", false, "Prints the version of the application")
	flag.Parse()
	if *version {
		fmt.Println("shift : proxy manager for proGY,", appVersion, appBuildPlatform, appBuildDate)
		return
	}

	if !IsFileExist(*progyConfigFile) || !IsFileExist(*proxyListFile) {
		fmt.Println("shift : Error : One or more of the files passed as config is not found.")
		return
	}
	conf, err := NewConfigFromFile(*progyConfigFile)
	if err != nil {
		fmt.Println("shift : Error : ", err)
		return
	}

	plist, err := NewProxyListFromFile(*proxyListFile)
	if err != nil {
		fmt.Println("shift : Error : ", err)
		return
	}

	conf.Check(plist)
	for {
		select {
		case <-time.After(time.Second * 5):
			err = conf.Check(plist)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
