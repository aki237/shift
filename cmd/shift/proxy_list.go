package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
)

// ProxyList is a type (map or structs) for containing Credentials of different remote proxies.
type ProxyList struct {
	Proxies []Proxy
}

// Random method is used to get a random proxy from a random sublist.
// This will continue till it gets a successfull proxy (meaning no 403, 302, 307 or 407s)
func (p ProxyList) Random() Proxy {
	px := Proxy{}
	for {
		x := rand.Intn(len(p.Proxies))
		fmt.Println(len(p.Proxies), x)
		px = p.Proxies[x]
		result := px.Check()
		if result >= 200 && result < 300 {
			break
		}
	}
	return px
}

// NewProxyListFromFile function is used to read the proxy list from the given file.
// The file has a TOML like syntax.
func NewProxyListFromFile(filename string) (*ProxyList, error) {
	p := &ProxyList{}
	p.Proxies = make([]Proxy, 0)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	sectionRegex := regexp.MustCompile(`^\[\s*(?P<address>[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{1,5})\s*\]$`)
	currentSection := ""
	for _, line := range lines {
		if sectionRegex.MatchString(line) {
			currentSection = sectionRegex.FindStringSubmatch(line)[len(sectionRegex.SubexpNames())-1]
			continue
		}
		username, password := getEqualityPairs(line)
		if username == "" && password == "" {
			continue
		}

		c := Proxy{Username: username, Password: password}
		c.RemoteProxyAddress = currentSection
		if currentSection == "" {
			fmt.Print("Shift : Warning : Proxies found in unspecified section. ")
			fmt.Print("Check your proxy list file for any inconsistencies.\n")
			c.RemoteProxyAddress = "untitled"
		}
		p.Proxies = append(p.Proxies, c)
	}
	return p, nil
}

// getEqualityPairs functions is used to split a string by "=" and return the parts, space trimmed.
// This is a bad parser. Soon will be replaced with a regex Matching one.
func getEqualityPairs(line string) (string, string) {
	splitted := strings.Split(strings.TrimSpace(line), "=")
	if len(splitted) != 2 {
		return "", ""
	}
	return strings.TrimSpace(splitted[0]), strings.TrimSpace(splitted[1])
}
