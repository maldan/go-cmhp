package cmhp_net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	hurl "net/url"
	"os"
	"path/filepath"
)

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
	Error      error  `json:"error"`
	Url        string `json:"url"`
	Headers    map[string]any
}

type HttpArgs struct {
	Url     string
	Headers map[string]string
	Method  string
	Proxy   string

	InputData  []byte
	InputJSON  map[string]any
	OutputJSON any
}

type RequestOptions struct {
	Headers map[string]string
	Proxy   string
}

type Response[T any] struct {
	StatusCode  int
	RawBody     []byte
	GenericBody T
	Error       error
	Url         string
}

func Request(args HttpArgs) HttpResponse {
	response := HttpResponse{}

	// Create client
	client := &http.Client{}

	// Set proxy
	if args.Proxy != "" {
		proxyUrl, _ := hurl.Parse(args.Proxy)
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
	response.Headers = make(map[string]any, 0)
	for k, v := range resp.Header {
		response.Headers[k] = v
	}

	// response.Headers = resp.Header.
	response.StatusCode = resp.StatusCode
	response.Body = body
	response.Url = newUrl

	// Out to JSON
	if args.OutputJSON != nil {
		json.Unmarshal(body, args.OutputJSON)
	}

	return response
}

func buildQuery(data map[string]any) string {
	out := "?"
	for k, v := range data {
		out += fmt.Sprintf("%v=%v&", k, v)
	}
	return out
}

func GetJson[T any](url string, method string, data map[string]any, opts RequestOptions) Response[T] {
	response := Response[T]{
		Url: url,
	}

	// Create client
	client := &http.Client{}

	// Set proxy
	if opts.Proxy != "" {
		proxyUrl, _ := hurl.Parse(opts.Proxy)
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Build query
	if method == "GET" || method == "DELETE" {
		response.Url += buildQuery(data)
	}

	// Prepare data
	out, _ := json.Marshal(data)
	inputData := bytes.NewBuffer(out)

	// Create request
	request, err := http.NewRequest(method, response.Url, inputData)
	if err != nil {
		response.Error = err
		return response
	}

	// Fill headers
	request.Header.Set("Content-Type", "application/json")
	for k, v := range opts.Headers {
		request.Header.Set(k, v)
	}

	// Do request
	resp, err := client.Do(request)
	if err != nil {
		response.Error = err
		return response
	}

	// Read
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		response.Error = err
		return response
	}

	// Fill
	response.StatusCode = resp.StatusCode
	response.RawBody = body
	json.Unmarshal(body, &response.GenericBody)

	return response
}

func DownloadFile(url string, dest string) error {
	r := Request(HttpArgs{Url: url, Method: "GET"})
	if r.StatusCode != 200 {
		return errors.New("can't get url")
	}

	// Create path for file
	err := os.MkdirAll(filepath.Dir(dest), 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, r.Body, 0777)
	return err
}
