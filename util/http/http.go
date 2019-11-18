package http

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

//NOTE:
// *****************************************************************
//	common http functions with metrics for project use
func HttpPost(url string, body []byte, headers map[string]string) (response http.Response, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*2)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpPut(url string, body []byte, headers map[string]string) (response http.Response, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*2)
	defer cncl()

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpDelete(url string, body []byte, headers map[string]string) (response http.Response, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*2)
	defer cncl()

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))
	if err != nil {
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpGet(url string, headers map[string]string) (response http.Response, error error) {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*2)
	defer cncl()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil
}
