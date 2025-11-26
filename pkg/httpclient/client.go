package httpclient

import (
	"errors"
	"net"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("not found")

type Config struct {
	Host string `default:"localhost" envconfig:"HTTP_CLIENT_HOST"`
	Port string `default:"8080"      envconfig:"HTTP_CLIENT_PORT"`
}

type Client struct {
	client http.Client
	host   string
}

func New(c Config) *Client {
	return &Client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		host: net.JoinHostPort(c.Host, c.Port),
	}
}
