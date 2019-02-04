package common

import (
	"bytes"
	"context"
	"github.com/rebelit/gome/notify"
	"net/http"
	"time"
)

//NOTE:  b, _ := json.Marshal(body)


func HttpPost(url string, body []byte, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil{
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		notify.MetricHttpOut(url, http.MethodPost, notify.FAILED)
		return http.Response{}, err
	}

	notify.MetricHttpOut(url, http.MethodPost, notify.SUCCESS)
	return *resp, nil
}

func HttpPut(url string, body []byte, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil{
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		notify.MetricHttpOut(url, http.MethodPut, notify.FAILED)
		return http.Response{}, err
	}

	notify.MetricHttpOut(url, http.MethodPut, notify.SUCCESS)
	return *resp, nil
}

func HttpDelete(url string, body []byte, headers map[string]string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))
	if err != nil{
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		notify.MetricHttpOut(url, http.MethodDelete, notify.FAILED)
		return http.Response{}, err
	}

	notify.MetricHttpOut(url, http.MethodDelete, notify.SUCCESS)
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
		notify.MetricHttpOut(url, http.MethodGet, notify.FAILED)
		return http.Response{}, err
	}

	notify.MetricHttpOut(url, http.MethodGet, notify.FAILED)
	return *resp, nil
}