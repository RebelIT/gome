# gome
Golang based home automation
GoLang Home (gome)

###Current Development Status:
* Stable (it runs) Tuya schedules work. other device control endpoints work. 
* ToDo: see issues. 

## About it
   I wanted to automate my home and all my IoT devices, I buy cheap and build from scratch.  There are many home
   automation systems out there but I didn't want any vendor or protocol lock. So I'm building my own starting with an 
   API only and later adding a web-ui for tablet control. Since I have many raspberry Pi's doing many things around my 
   house this can tie it all together for central management and control. 

   Starts up with a base devices.json loaded into redis from there endpoints to add/remove devices from the database 
   and update the json.
   
#### Still a work in progress as I contintue to learn GOLang.  but it works! :)

## Test it:

```
docker build -t gometest .
docker run -v $PWD:/go/src/github.com/rebelit/gome -i -t gometest /bin/bash
```

## Deploy it
### Note:
* tuya wall outlets and light switches control is dependent on tuya-api, had issues with performance and blocking on 
rPI's due to cpu saturation. 
see issue #17 for fixes*
* tested on latest raspbian stretch

Ansible deployment: see doco in ansible dir


### Supported Devices
   [x] Custom RaspberryPi API
   
   [x] Roku
   
   [x] Tuya WiFi outlets (with external tuya-cli dependency)
   
   [x] Scheduler for outlets to auto turn on and off on date/time
   
   [x] Tuya WiFi light switches (with external tuya-cli dependency
   
   [ ] Plex API - not sure yet what i can do with it
   
   [x] Amazon Alexa integration - voice control this central management beast
      
   [ ] web-ui - feature issue added
   
   [ ] DEEBOT (ECOVACS Robotics) vaccum api integration, control & metrics
   
   
   
## RaspberryPi IoT devices
   
   each one of my rPI's has an API endpoint to control it
   
   Example:  [API ansible Role to setup an API on a raspberry pi](https://github.com/RebelIT/ansible-piDAK) this rPI api
   is used in this project to control the individual rPI's.
