package devices

import (
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Roku App ID's
const NETFLIX = 12
const PLEX = 13535
const SLING = 46041
const PANDORA = 28
const PRIME_VIDEO = 13
const GOOGLE_PLAY = 50025
const HBOGO = 8378
const YOUTUBE = 837

func getAppId(app string)(string, error){
	id := 0
	switch app {
	case "netflix":
		id = NETFLIX
	case "plex":
		id = PLEX
	case "sling":
		id = SLING
	case "pandora":
		id = PANDORA
	case "prime":
		id = PRIME_VIDEO
	case "google":
		id = GOOGLE_PLAY
	case "hbo":
		id = HBOGO
	case "youtube":
		id = YOUTUBE
	default:
		return "", errors.New("no app "+app+" found")
	}

	return strconv.Itoa(id), nil
}

func RokuDeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	uriPart := "/"
	alive := false

	resp, err := rokuGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		alive = true
	}

	if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

func launchApp(deviceName string, app string) error {
	id, err := getAppId(app)
	if err != nil{
		return err
	}
	uri := "/launch/"+id
	resp, err := rokuPost(uri, deviceName)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return errors.New("non-200 status code return")
	}
	return nil
}


// http wrappers
func rokuPost(uriPart string, deviceName string) (http.Response, error) {
	d, err := GetDevice(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpPost(url,nil,nil)
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}

func rokuGet(uriPart string, deviceName string) (http.Response, error) {
	d, err := GetDevice(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpGet(url, nil)
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}