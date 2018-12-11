package common

// flip this to use when deployed
const FILE  = "/etc/gome/devices.json"
const SECRETS = "/etc/gome/secrets.json"

// flip this on for local testing
//const FILE  = "./devices.json"
//const SECRETS = ".secrets.json"

// set to empty string to disable metrics
const STATSD = "127.0.0.1:8125"
