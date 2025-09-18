package client

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func DoRequest(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}

	// untuk membaca body response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
