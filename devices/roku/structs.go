package roku

import "net/http"

type Roku struct {
	address string
	client  *http.Client
}