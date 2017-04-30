package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

// Config struct is for containing the variables extracted from the proGY config file.
type Config struct {
	// ListenAddress specifies the port at which proGY is listening to locally.
	ListenAddress string `json:"listenaddress"`
	// Creds is an Array of Proxy  Struct which contains al the proxy related details
	// like remote address, username, password.
	Creds []Proxy `json:"Creds"`
	// DomainCacheFile points to a file that contains the domain name to IP cache.
	DomainCacheFile string `json:"domaincachefile"`
	// Verbose states whether proGY output is verbose or not.
	// Legacy code in proGY. To be removed in here and in proGY too.
	Verbose bool `json:"verbose"`
	// LoggerPort is the port at which TCP port the proGY logger runs.
	LoggerPort int `json:"loggerport"`
	// ControlPort is the unix socket at which the control server of proGY runs.
	ControlPort string `json:"controlport"`
	// filename string for storing the config file name.
	Filename string `json:"-"`
}

// Check is used to check all the proxies in the Creds.
// This inturn calls the `Check` method on each of the Proxy structs,
// with a context passed for cancellation during timeout.
func (c *Config) Check(proxyList *ProxyList) error {
	changed := false
	for i, val := range c.Creds {
		result := val.Check()
		switch {
		case result == ErrDial:
			fmt.Printf("%s@%s fails : Unable to connect\n", val.Username, val.RemoteProxyAddress)
		case result == ErrConnWrite:
			fmt.Printf("%s@%s fails : Unable to write in the connection\n", val.Username, val.RemoteProxyAddress)
		case result == ErrConnRead:
			fmt.Printf("%s@%s fails : Unable to read from the connection\n", val.Username, val.RemoteProxyAddress)
		case result >= 200 && result < 300:
			fmt.Printf("%s@%s passes : %d\n", val.Username, val.RemoteProxyAddress, result)
		case result == ErrProxyRedirection:
		case result == ErrProxySuccessfulRedirection:
			fmt.Printf("%s@%s fails : Quota Over\n", val.Username, val.RemoteProxyAddress)
		case result == ErrForbidden:
			fmt.Printf("%s@%s fails : Forbidden\n", val.Username, val.RemoteProxyAddress)
		case result == ErrAuthWrong:
			fmt.Printf("%s@%s fails : Wrong Proxy Cedentials\n", val.Username, val.RemoteProxyAddress)
		case result == ErrServerError:
			fmt.Printf("%s@%s fails : Poll Check Server Problem\n", val.Username, val.RemoteProxyAddress)
		}
		if (result >= 300 && result < 500) || result < 0 {
			changed = true
			p := proxyList.Random()
			c.Creds[i].RemoteProxyAddress = p.RemoteProxyAddress
			c.Creds[i].Username = p.Username
			c.Creds[i].Password = p.Password
		}
	}
	if changed {
		barr, err := json.Marshal(c)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(c.Filename, barr, 0644)
		if err != nil {
			return err
		}
		conn, err := net.Dial("unix", c.ControlPort)
		if err != nil {
			return err
		}
		_, err = conn.Write([]byte("RELOAD " + c.Filename + "\n"))
		if err != nil {
			return err
		}
		resp, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		if string(resp) != "SUCCESS\n" {
			return errors.New("Unable to reload the proGY daemon")
		}
		fmt.Println("Daemon reloaded")
	}
	return nil
}

// NewConfigFromFile is used to open the specified config (.progy) file, read contents
// and parse it into the Config struct.
func NewConfigFromFile(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	conf.Filename = filename
	if conf.ControlPort == "" {
		conf.ControlPort = "/tmp/proGY-control"
	}
	return conf, nil
}
