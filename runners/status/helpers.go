package status

import (
	"math/rand"
	"time"
)

func randomizeCollection() time.Duration {
	min := 0
	max := 300

	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(max - min) + min

	return time.Duration(int64(i))
}