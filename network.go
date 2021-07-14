package cmhp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func RequestGetAsBin(url string, headers map[string]string) ([]byte, int, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, errors.New(string(body))
	}

	return body, resp.StatusCode, nil
}

func RequestGetAsText(url string, headers map[string]string) (string, int, error) {
	data, status, err := RequestGetAsBin(url, headers)
	return string(data), status, err
}

func RequestGetAsJSON(url string, headers map[string]string, int, s interface{}) (int, error) {
	data, status, err := RequestGetAsBin(url, headers)
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(data, &s)
	return status, err
}

func RequestPostAsBin(url string, headers map[string]string, data []byte) ([]byte, int, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, errors.New(string(body))
	}

	return body, resp.StatusCode, nil
}

func RequestPostAsJson(url string, headers map[string]string, v interface{}) (string, int, error) {
	out, err := json.Marshal(v)
	if err != nil {
		return "", 500, err
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	data, status, err := RequestPostAsBin(url, headers, out)
	return string(data), status, err
}
