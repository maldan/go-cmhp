package cmhp_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
	Error      error  `json:"error"`
	Url        string `json:"url"`
}

type HttpArgs struct {
	Url     string
	Headers map[string]string
	Method  string
	Proxy   string

	InputData  []byte
	InputJSON  map[string]any
	OutputJSON *map[string]any
}

func Request(args HttpArgs) HttpResponse {
	response := HttpResponse{}

	// Create client
	client := &http.Client{}

	// Set proxy
	if args.Proxy != "" {
		proxyUrl, _ := url.Parse(args.Proxy)
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	newUrl := args.Url

	// Create request
	var req *http.Request
	if args.Method == "GET" || args.Method == "DELETE" {
		newUrl += "?"
		for k, v := range args.InputJSON {
			newUrl += fmt.Sprintf("%v=%v&", k, v)
		}
	}
	if args.Method == "POST" || args.Method == "PATCH" || args.Method == "PUT" {
		if args.InputJSON != nil {
			out, _ := json.Marshal(args.InputJSON)
			r, err := http.NewRequest(args.Method, newUrl, bytes.NewBuffer(out))
			if err != nil {
				response.Error = err
				return response
			}
			req = r
		} else {
			r, err := http.NewRequest(args.Method, newUrl, bytes.NewBuffer(args.InputData))
			if err != nil {
				response.Error = err
				return response
			}
			req = r
		}
	} else {
		r, err := http.NewRequest(args.Method, newUrl, nil)
		if err != nil {
			response.Error = err
			return response
		}
		req = r
	}

	// Fill headers
	if args.Headers == nil {
		args.Headers = make(map[string]string)
	}
	for k, v := range args.Headers {
		req.Header.Set(k, v)
	}
	if args.InputJSON != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		response.Error = err
		return response
	}

	// Read
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return response
	}

	// Fill
	response.StatusCode = resp.StatusCode
	response.Body = body
	response.Url = newUrl

	// Out to JSON
	if args.OutputJSON != nil {
		json.Unmarshal(body, args.OutputJSON)
	}

	return response
}
