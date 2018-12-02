package common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostWebReq (reqBody interface{}, reqUrl string) (respCode int, respError error){
	contentType := "application/json"
	body, err := json.Marshal(reqBody)

	if err != nil{
		return 0, err
	}

	resp, err := http.Post(reqUrl, contentType, bytes.NewReader(body))
	if err != nil{
		return 0, err
	}

	return resp.StatusCode, nil
}
