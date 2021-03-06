package common

// flip this to use when deployed
const FILE  = "/etc/apps/gome-server/devices.json"
const SECRETS = "/etc/apps/gome-server/secrets.json"
const DBPATH = "/etc/apps/gome-server/db/badger"
const STATSD = "127.0.0.1:8125"     // set to empty string to disable metrics
const SCHEDULE_MIN = 1
const STATUS_MIN = 5
const AWS_SEC = 2
const HTTP_TIMEOUT = 2
const SUCCESS = "success"
const FAILED = "failed"

// flip this on for local testing
//const FILE  = "./devices.json"
//const SECRETS = "./secrets.json"
//const STATSD = ""
//const DBPATH = "./testing/db/badger"