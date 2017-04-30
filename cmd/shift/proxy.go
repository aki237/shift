package main

import (
	"bufio"
	"encoding/base64"
	"net"
	"net/http"
	"time"
)

// Constants denoting return values from the Check Method of Proxy
const (
	ErrConnWrite                  int = -3
	ErrConnRead                       = -2
	ErrDial                           = -1
	ErrProxyRedirection               = 307
	ErrProxySuccessfulRedirection     = 302
	ErrForbidden                      = 403
	ErrAuthWrong                      = 407
	ErrServerError                    = 500
)

// Proxy struct is for containing the proxy details like remote proxy address, username and password.
type Proxy struct {
	RemoteProxyAddress string `json:"remoteproxyaddress"`
	Username           string `json:"username"`
	Password           string `json:"password"`
}

// Check method is used to check for any HTTP errors like 403, 407, 302, 307 etc.,
func (p Proxy) Check() int {
	conn, err := net.DialTimeout("tcp", p.RemoteProxyAddress, 3*time.Second)
	if err != nil {
		return -1
	}
	header := "HEAD http://gstatic.com/generate_204 HTTP/1.1\nProxy-Authorization: Basic "
	header += p.Base64Encode() + "\n\n"
	_, err = conn.Write([]byte(header))
	if err != nil {
		return -3
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return -2
	}
	return resp.StatusCode
}

// Base64Encode method is used to encode the username and password to Base64 standard encoding
func (p Proxy) Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(p.Username + ":" + p.Password))
}
