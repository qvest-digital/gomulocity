package inventory

import (
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Username   string
	Password   string
}