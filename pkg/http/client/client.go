package client

import (
	"bytes"
	"sora_landing_be/pkg/logger"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"
)

type HTTPClient struct {
	Client *http.Client
}

var (
	instance *HTTPClient
	once     = &sync.Once{}
)

type RequestAttributes struct {
	Method      string
	Headers     map[string]string
	QueryParams url.Values
	Body        interface{}
}

func GetClient() *HTTPClient {
	once.Do(func() {
		instance = &HTTPClient{
			Client: &http.Client{
				Timeout: time.Second * 20,
			},
		}
	})
	return instance
}

// Call is wrapper for http call, supports any http call method eg, POST,GET,PUT etc.
// Parameters:
//   - urls: Url target which want to be called
//   - response: Response body from target url, response has to be pointer type
//   - attributes : HTTP client attributes when do call operation
//
// Returns:
//   - statusCode : response status code, default value is 0
//   - err : error
func (c *HTTPClient) Call(urls string, response interface{}, attributes RequestAttributes) (err error) {
	if attributes.Method == "" {
		return fmt.Errorf("http client call method is required")
	}

	if response != nil && reflect.ValueOf(response).Kind() != reflect.Ptr {
		return fmt.Errorf("response must be a pointer")
	}

	parsedURL, err := url.Parse(urls)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	if attributes.QueryParams != nil {
		q := parsedURL.Query()
		for k, v := range attributes.QueryParams {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		parsedURL.RawQuery = q.Encode()
	}

	var reqBody []byte
	if attributes.Body != nil {
		reqBody, err = json.Marshal(attributes.Body)
		if err != nil {
			return NewHTTPClientError(0, fmt.Sprintf("error marshaling body: %w", err))
		}
	}

	req, err := http.NewRequest(attributes.Method, parsedURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range attributes.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Body != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}
			return NewHTTPClientError(resp.StatusCode, fmt.Sprintf("unexpected status code: %d, response : %s", resp.StatusCode, string(body)))
		}

		return NewHTTPClientError(resp.StatusCode, fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return NewHTTPClientError(resp.StatusCode, fmt.Sprintf("error decoding response: %w", err))
		}
	}

	return nil
}
