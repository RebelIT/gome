package devices

import (
	"fmt"
	db "github.com/rebelit/gome/database"
)

func GetProfile(name string) (profile Profile, error error) { //Gets the device profile from the database

	value, err := db.Get(name)
	if err != nil {
		return Profile{}, err
	}

	profile, err = stringToStruct(value)
	if error != nil {
		return Profile{}, nil
	}

	return profile, nil
}

func UpdateProfileStatus(name string, status bool) error { //Updates the device profile status in the database
	profile, err := GetProfile(name)
	if err != nil {
		return err
	}

	profile.State.Status = status

	err = db.Add(name, profile.structToString())
	if err != nil {
		return err
	}

	return nil
}

func UpdateProfileAlive(name string, alive bool) error { //Updates the device profile alive in the database
	profile, err := GetProfile(name)
	if err != nil {
		return err
	}

	profile.State.Alive = alive

	err = db.Add(name, profile.structToString())
	if err != nil {
		return err
	}

	return nil
}

func UpdateProfile(newProfile Profile) error { //Updates the entire profile in the database
	err := db.Add(newProfile.Name, newProfile.structToString())
	if err != nil {
		return err
	}

	return nil
}

func GetAllProfiles() (profiles Devices, err error) { //Gets the full inventory of device profiles from the database
	profiles = Devices{}
	list := []Profile{}

	keys, err := db.GetAll()
	if err != nil {
		return Devices{}, err
	}

	for _, key := range keys {
		p, err := GetProfile(key)
		if err != nil {
			continue
		}
		list = append(list, p)
	}
	profiles.Profiles = list

	return profiles, nil
}

func GetAllDeviceNames() (names []string, err error) { //Gets a list of devices from the database
	list := []string{}

	keys, err := db.GetAll()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		list = append(list, key)
	}

	return list, nil
}

func RemoveProfile(name string) error { //Gets the device profile from the database

	value, err := db.Get(name)
	if err != nil {
		return err
	}
	if value == "" {
		return fmt.Errorf("%s not found", name)
	}

	err = db.Del(name)
	if err != nil {
		return err
	}

	return nil
}
