package cmhp_net

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
	Error      error  `json:"error"`
}

type HttpArgs struct {
	Url     string
	Headers map[string]string
	Method  string

	InputData  []byte
	InputJSON  interface{}
	OutputJSON interface{}
}

func Request(args HttpArgs) HttpResponse {
	response := HttpResponse{}

	// Create request
	client := &http.Client{}
	var req *http.Request
	if args.Method == "POST" {
		if args.InputJSON != nil {
			out, _ := json.Marshal(args.InputJSON)
			r, err := http.NewRequest(args.Method, args.Url, bytes.NewBuffer(out))
			if err != nil {
				response.Error = err
				return response
			}
			req = r
		} else {
			r, err := http.NewRequest(args.Method, args.Url, bytes.NewBuffer(args.InputData))
			if err != nil {
				response.Error = err
				return response
			}
			req = r
		}
	} else {
		r, err := http.NewRequest(args.Method, args.Url, nil)
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

	// Out to JSON
	if args.OutputJSON != nil {
		json.Unmarshal(body, args.OutputJSON)
	}

	return response
}
