package roku

import (
	"github.com/pkg/errors"
	"strconv"
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