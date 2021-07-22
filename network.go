package cmhp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
}

type HttpArgs struct {
	Url     string
	Headers map[string]string
	Method  string
	Data    []byte
	Params  map[string]interface{}
	JSON    interface{}
}

func Request(args HttpArgs) HttpResponse {
	response := HttpResponse{}

	// Create request
	client := &http.Client{}
	var req *http.Request
	if args.Method == "POST" {
		if args.JSON != nil {
			out, _ := json.Marshal(args.JSON)
			r, err := http.NewRequest(args.Method, args.Url, bytes.NewBuffer(out))
			if err != nil {
				return response
			}
			req = r
		} else {
			r, err := http.NewRequest(args.Method, args.Url, bytes.NewBuffer(args.Data))
			if err != nil {
				return response
			}
			req = r
		}
	} else {
		r, err := http.NewRequest(args.Method, args.Url, nil)
		if err != nil {
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
	if args.JSON != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return response
	}

	// Read
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response
	}

	// Fill
	response.StatusCode = resp.StatusCode
	response.Body = body

	return response
}

func RequestGetBin(url string, headers map[string]string) HttpResponse {
	response := HttpResponse{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response
	}

	response.StatusCode = resp.StatusCode
	response.Body = body

	return response
}

func RequestGetText(url string, headers map[string]string) HttpResponse {
	response := RequestGetBin(url, headers)
	return response
}

func RequestGetJSON(url string, headers map[string]string, s interface{}) HttpResponse {
	response := RequestGetBin(url, headers)
	err := json.Unmarshal(response.Body, &s)
	if err != nil {
		response.StatusCode = 0
	}
	return response
}

func RequestPostBin(url string, headers map[string]string, data []byte) HttpResponse {
	response := HttpResponse{}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response
	}

	response.StatusCode = resp.StatusCode
	response.Body = body

	return response
}

func RequestPostJSON(url string, headers map[string]string, v interface{}) HttpResponse {
	out, err := json.Marshal(v)
	if err != nil {
		return HttpResponse{}
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	response := RequestPostBin(url, headers, out)
	return response
}
