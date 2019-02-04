package httpClient

import (
	"context"
	"net/http"
	"time"
)

func Get(url string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
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

func Post(url string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil{
		return http.Response{}, err
	}
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}

func Put(url string)(response http.Response, err error){
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
	defer cncl()

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil{
		return http.Response{}, err
	}
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	return *resp, nil
}

//func request()(http.Request, error){
//	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
//	defer cncl()
//
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil{
//		return *req, err
//	}
//}