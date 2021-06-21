package cmhp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func RequestGetAsBin(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func RequestGetAsText(url string, headers map[string]string) (string, error) {
	data, err := RequestGetAsBin(url, headers)
	return string(data), err
}

func RequestGetAsJSON(url string, headers map[string]string, s interface{}) error {
	data, err := RequestGetAsBin(url, headers)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &s)
	return err
}
