package rpi

import "net/http"

type Pi struct {
	address string
	client  *http.Client
}