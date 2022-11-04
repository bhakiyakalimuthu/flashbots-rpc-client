package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
)

type HttpClient struct {
	client  *http.Client
	url     string
	mu      sync.Mutex // protects headers
	headers http.Header
}

func DialHttpClient(rawURL string) (*HttpClient, error) {
	_, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	headers := make(http.Header, 2)
	headers.Set("accept", "application/json")
	headers.Set("content-type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	return &HttpClient{
		client:  client,
		url:     rawURL,
		headers: headers,
	}, nil
}

func DialHttpClientWithLocalHost(rawURL string) (*HttpClient, error) {
	return DialHttpClient("http://localhost:8080")
}

func (hc *HttpClient) doRequest(ctx context.Context, msg interface{}, signature string) (io.ReadCloser, error) {
	// prepare body
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// create request
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, hc.url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}

	// set headers
	hc.mu.Lock()
	request.Header = hc.headers.Clone()
	request.Header.Set("x-flashbots-signature", signature)
	hc.mu.Unlock()

	// send request
	resp, err := hc.client.Do(request)
	if err != nil {
		return nil, err
	}

	// handle response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var _body []byte
		if _, err = buf.ReadFrom(resp.Body); err == nil {
			_body = buf.Bytes()
		}

		return nil, common.HTTPError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       _body,
		}
	}
	return resp.Body, nil
}

func (hc *HttpClient) CallContext(ctx context.Context, msg interface{}, signature string) (*common.JSONRPCMessage, error) {
	respBody, err := hc.doRequest(ctx, msg, signature)
	respBody.Close()
	if err != nil {
		return nil, err
	}

	var resp *common.JSONRPCMessage
	if err = json.NewDecoder(respBody).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}
