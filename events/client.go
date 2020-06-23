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

func (client *Client) delete(path string) ([]byte, int, error) {
	url := client.BaseURL + path
	return client.request(http.MethodDelete, url, []byte{})
}

func (client *Client) put(path string, body []byte) ([]byte, int, error) {
	url := client.BaseURL + path
	return client.request(http.MethodPut, url, body)
}

func (client *Client) post(path string, body []byte) ([]byte, int, error) {
	url := client.BaseURL + path
	return client.request(http.MethodPost, url, body)
}

func (client *Client) get(path string) ([]byte, int, error) {
	url := client.BaseURL + path
	return client.request(http.MethodGet, url, []byte{})
}

func (client *Client) request(method, url string, body []byte) ([]byte, int, error) {
	log.Printf("HTTP %s on URL %s", method, url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, 0, err
	}

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
