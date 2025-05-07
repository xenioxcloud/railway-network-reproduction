package oxylabs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type HttpClient interface {
	Do(ctx context.Context, req *http.Request) (resBody []byte, err error)
}

var (
	Err404 = errors.New("error 404")
)

type oxylabsHttpClient struct {
	username string
	password string
	entry    string

	client *http.Client
}

const (
	MAX_RETRY_COUNT = 5
)

func NewOxylabsHttpClient(username, password, entry string) HttpClient {
	proxy, err := url.Parse(fmt.Sprintf("http://user-%s:%s@%s", username, password, entry))
	if err != nil {
		log.Fatalln(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	client := &oxylabsHttpClient{
		username: username,
		password: password,
		entry:    entry,
		client:   httpClient,
	}
	return client
}

func (h *oxylabsHttpClient) Do(ctx context.Context, req *http.Request) (resBody []byte, err error) {
	var retry bool
	retries := 0
	backoff := time.Second
	for {
		resBody, retry, err = h.do(req)
		if !retry {
			break
		}
		if retries == MAX_RETRY_COUNT {
			err = fmt.Errorf("failed %d consecutive retries: %v", MAX_RETRY_COUNT, err)
			return nil, err
		}
		retries++
		log.Printf("Request failed %d times (%v), trying again\n", retries, err)
		time.Sleep(backoff)
		backoff *= 2
	}
	return
}

func (h *oxylabsHttpClient) do(req *http.Request) (resBody []byte, retry bool, err error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := h.client.Do(req)
	if err != nil {
		return nil, true, err
	}

	if res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= http.StatusInternalServerError {
		return nil, true, nil
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, false, Err404
	}
	if res.StatusCode != http.StatusOK {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, false, err
		}
		err = fmt.Errorf("response status code is %d: %v", res.StatusCode, string(resBody))
		return nil, false, err
	}

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, true, err
	}
	return resBody, false, nil
}
