package notify

import (
	"encoding/json"
	"errors"
	"github.com/rebelit/gome/common"
	"io/ioutil"
)

func getSecrets() (Secrets, error){
	s := Secrets{}
	secretsFile, err := ioutil.ReadFile(common.SECRETS)
	if err != nil {
		return Secrets{}, err
	}
	if err := json.Unmarshal(secretsFile, &s); err != nil{
		errorMsg := errors.New("unable to read secrets")
		return Secrets{}, errorMsg
	}
	return s, nil
}