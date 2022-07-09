package util

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func HTTPGet(reqURL string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	log.Printf("GET %s", reqURL)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Got status code: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}

func HTTPPost(reqURL string, headers map[string]string, formData map[string]string) (*http.Response, error) {
	urlValues := url.Values{}
	for k, v := range formData {
		urlValues.Set(k, v)
	}
	encodedURLValues := urlValues.Encode()

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(encodedURLValues))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("POST %s (%s)", reqURL, encodedURLValues)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Got status code: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
