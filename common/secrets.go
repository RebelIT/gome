package common

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

func GetSecrets() (Secrets, error){
	s := Secrets{}
	secretsFile, err := ioutil.ReadFile(SECRETS)
	if err != nil {
		return Secrets{}, err
	}
	if err := json.Unmarshal(secretsFile, &s); err != nil{
		errorMsg := errors.New("unable to read secrets")
		return Secrets{}, errorMsg
	}
	return s, nil
}