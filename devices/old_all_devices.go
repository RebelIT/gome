package devices

import (
	db "github.com/rebelit/gome/database"
)

//func UpdateProfileStatus(name string, status bool) error { //Updates the device profile status in the database
//	profile, err := GetProfile(name)
//	if err != nil {
//		return err
//	}
//
//	profile.State.Status = status
//
//	err = db.Add(name, profile.structToString())
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func UpdateProfileAlive(name string, alive bool) error { //Updates the device profile alive in the database
//	profile, err := GetProfile(name)
//	if err != nil {
//		return err
//	}
//
//	profile.State.Alive = alive
//
//	err = db.Add(name, profile.structToString())
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func UpdateProfile(newProfile Profile) error { //Updates the entire profile in the database
	err := db.Add(newProfile.Name, newProfile.structToString())
	if err != nil {
		return err
	}

	return nil
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

//func RemoveProfile(name string) error { //Gets the device profile from the database
//
//	value, err := db.Get(name)
//	if err != nil {
//		return err
//	}
//	if value == "" {
//		return fmt.Errorf("%s not found", name)
//	}
//
//	err = db.Del(name)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
