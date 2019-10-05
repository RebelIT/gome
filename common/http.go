package common

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//NOTE:
// *****************************************************************
//	common http functions with metrics for project use
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
		return http.Response{}, err
	}

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
		return http.Response{}, err
	}

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
		return http.Response{}, err
	}

	return *resp, nil
}

func HttpGet(url string, headers map[string]string)(response http.Response, error error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second * HTTP_TIMEOUT)
	defer cncl()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	for key, value := range headers {
		req.Header.Set(key,value)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}


// *****************************************************************
// Http Response helper functions
func ReturnOk(w http.ResponseWriter, r *http.Request, response interface{}){
	code := http.StatusOK
	MetricHttpIn(r.RequestURI, code, r.Method)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
	}
	return
}

func ReturnBad(w http.ResponseWriter, r *http.Request){
	code := http.StatusBadRequest
	MetricHttpIn(r.RequestURI, code, r.Method)
	w.WriteHeader(code)
	return
}

func ReturnInternalError(w http.ResponseWriter, r *http.Request){
	code := http.StatusInternalServerError
	MetricHttpIn(r.RequestURI, code, r.Method)
	w.WriteHeader(code)
	return
}