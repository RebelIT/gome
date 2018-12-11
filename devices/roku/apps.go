package roku

import (
	"github.com/pkg/errors"
	"strconv"
)

//Roku App ID's
const NETFLIX = 12
const PLEX = 12
const SLING = 12
const PANDORA = 12
const PRIME_VIDEO = 12
const GOOGLE_PLAY = 12
const HBOGO = 12
const YOUTUBE = 12

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