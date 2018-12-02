# gome
Golang based home automation
GoLang Home (gome)

## About it
   I wanted to automate my home and all my IoT devices, I buy cheap and build from scratch.  There are many home
   automation systems out there but I didn't want any vendor or protocol lock. So im building my own starting with an 
   API only and later adding a web-ui for tablet control. since I have many raspberry Pi's doing many things around my 
   house this can tie it all together for central management and control. 

   starts up with a base devices.json in the root dir to load into a redis database from there adding more endpoints
   to add/remove devices fro the database and update the json.
   

### Work in progress as I learn GOLang and have time to play. :)

### Custom RaspberryPi setup
   
   each one of my rPI's has an API endpoint to control it
   
   Example:  [API ansible Role to setup an API on a raspberry pi](https://github.com/RebelIT/ansible-piDAK)  Just add 
   more as needed to control your rPI functions and apps
 

### supported Devices

   * Custom RaspberryPi API
   * Roku
   * Tuya WiFi outlets (with external tuya-cli dependency)
   * Scheduler for outlets to auto turn on and off on date/time
   * (Coming Soon) Tuya WiFi light switches (with external tuya-cli dependency)
   * (Coming Soon) Plex API - not sure yet what i can do with it
   * (Coming Soon) Amazon Alexa integration - voice control this central management beast
   * (Coming Soon) Scheduler for light switches to auto turn on and off on date/time
   
## Usage notes
Need to automate these yet
   * devices.json and secrets.json required in the root gome directory
   * still dependent in redis (working on an easier lightweight key:value runtime DB)
   * refer to examples in example directory for api call examples