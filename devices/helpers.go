package devices

import (
	"encoding/json"
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

func stringToStruct(profile string) (device Profile, error error){
	d := Profile{}

	if err := json.Unmarshal([]byte(profile), &d); err != nil {
		return d, err
	}

	return d, nil
}

func (profile *Profile)structToString() string{
	value, err := json.Marshal(profile)
	if err != nil {
		return ""
	}
	return string(value)
}