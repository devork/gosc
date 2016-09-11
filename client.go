package osc

import (
	"fmt"
	"net/http"
	"strings"
)

// Camera type should be used to send commands to the remote server
type Camera struct {
	uri    string
	client http.Client
}

func (c *Camera) Info() (*Info, error) {
	resp, err := c.client.Get(c.uri + "/osc/info")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	info, err := parseInfo(resp.Body)

	if err != nil {
		return nil, err
	}

	return info, nil
}

// New creates a Camera instance
func New(host string, port int) (*Camera, error) {

	if host == "" {
		return nil, fmt.Errorf("no host provided to connect to")
	}

	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}

	if port < 0x1 || port > 0xFFFF {
		port = 80
	}

	uri := fmt.Sprintf("%s:%d", host, port)

	return &Camera{
		uri:    uri,
		client: http.Client{},
	}, nil
}
