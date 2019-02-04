package common

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

//NOTE:
//	Project uses json all request body gets marshaled into json
//

func HttpPost(url string, body interface{}, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	b, _ := json.Marshal(body)
	req.Body.Read(b)

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpPut(url string, body interface{}, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	b, _ := json.Marshal(body)
	req.Body.Read(b)

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpDelete(url string, body interface{}, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	b, _ := json.Marshal(body)
	req.Body.Read(b)

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpGet(url string)(response http.Response, error error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}