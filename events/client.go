package gomulocity_event

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Username   string
	Password   string
}

func (client *Client) post(path string, body string) ([]byte, int, error) {
	url := client.BaseURL + path

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, 0, err
	}

	return client.request(req)
}

func (client *Client) get(path string) ([]byte, int, error) {
	url := client.BaseURL + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, 0, err
	}

	return client.request(req)
}

func (client *Client) request(req *http.Request) ([]byte, int, error) {
	log.Printf("HTTP %s on URL %s", req.Method, req.URL)

	req.SetBasicAuth(client.Username, client.Password)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		log.Printf("An error occured: %s", err.Error())
		return nil, 0, err
	}
	log.Printf("Got status %d", resp.StatusCode)
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading from stream: %s", err.Error())
		return nil, 0, err
	}

	log.Printf("Debug: Response body was: %s", result)
	return result, resp.StatusCode, nil
}
