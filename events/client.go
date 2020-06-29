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

// Returns an empty header map
func EmptyHeader() map[string][]string {
	return map[string][]string{}
}

// Returns a pre filled header map with an "Accept" header.
func AcceptHeader(accept string) map[string][]string {
	return map[string][]string{"Accept": {accept}}
}

func (client *Client) delete(path string, header map[string][]string) ([]byte, int, error) {
	return client.request(http.MethodDelete, path, []byte{}, header)
}

func (client *Client) put(path string, body []byte, header map[string][]string) ([]byte, int, error) {
	return client.request(http.MethodPut, path, body, header)
}

func (client *Client) post(path string, body []byte, header map[string][]string) ([]byte, int, error) {
	return client.request(http.MethodPost, path, body, header)
}

func (client *Client) get(path string, header map[string][]string) ([]byte, int, error) {
	return client.request(http.MethodGet, path, []byte{}, header)
}

func (client *Client) request(method, path string, body []byte, header map[string][]string) ([]byte, int, error) {
	url := client.BaseURL + path
	log.Printf("HTTP %s on URL %s", method, url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, 0, err
	}

	req.SetBasicAuth(client.Username, client.Password)
	for header, values := range header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

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
