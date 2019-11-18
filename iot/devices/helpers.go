package devices

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func randomizeCollection() time.Duration {
	min := 0
	max := 300

	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(max-min) + min

	return time.Duration(int64(i))
}

func stringToStruct(profile string) (device Profile, error error) {
	d := Profile{}

	if err := json.Unmarshal([]byte(profile), &d); err != nil {
		return d, err
	}

	return d, nil
}

func (profile *Profile) structToString() string {
	value, err := json.Marshal(profile)
	if err != nil {
		return ""
	}
	return string(value)
}

func (a *Action) constructAction() string {
	p1 := ""
	p2 := ""
	p3 := ""
	p4 := ""
	p5 := ""

	if a.Arg1 != "" {
		p1 = fmt.Sprintf("%s", a.Arg1)
	}
	if a.Arg1 != "" {
		p2 = fmt.Sprintf("/%s", a.Arg2)
	}
	if a.Arg1 != "" {
		p3 = fmt.Sprintf("/%s", a.Arg3)
	}
	if a.Arg1 != "" {
		p4 = fmt.Sprintf("/%s", a.Arg4)
	}
	if a.Arg1 != "" {
		p5 = fmt.Sprintf("/%s", a.Arg5)
	}

	return fmt.Sprintf("%s%s%s%s%s", p1, p2, p3, p4, p5)
}

func validateAction(profile Profile, actionName string) (action Action, error error) {
	action = Action{}
	validated := false
	for _, a := range profile.Actions {
		if actionName == a.Name {
			validated = true
			actionName = a.Name
			action.Arg1 = a.Arg1
			action.Arg2 = a.Arg2
			action.Arg3 = a.Arg3
			action.Arg4 = a.Arg4
			action.Arg5 = a.Arg5
		}
	}
	if !validated {
		return action, fmt.Errorf("unable to validate action: %s for %s", actionName, profile.Name)
	}

	return action, nil
}
