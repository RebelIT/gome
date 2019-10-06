# gome_badger
BadgerDB wrapper to simplify the usage

# usage examples:
`go get github.com/rebelit/gome_badger`

basic usage:
```
//Initialize a new database
d, err := db.New("/tmp/database")
if err != nil{
    fmt.Println("I didn't create a new DB")
}
defer d.Close()


//Add some data to the database
if err := d.Set("coolKey", "coolData"); err != nil{
    fmt.Println("I didn't set any data")
}


//Add some more data to the database
if err := d.Set("coolKey2", "coolData2"); err != nil{
    fmt.Println("I didn't set any data")
}


//Get a single key's value
value, err := d.Get("coolKey")
if err != nil{
    fmt.Println("I didn't get any data")
}
fmt.Println(value)


//Get all keys from the database
keys, err := d.GetAllKeys()
if err != nil{
    fmt.Println("I didn't get any data")
}
for _, k := range keys{
    fmt.Println(k)
}


//Delete a k/v pair from the database
if err := d.Delete("coolKey"); err != nil{
    fmt.Println("I didn't delete any data")
}


//Delete everything
if err := d.DeleteAll(); err != nil{
    fmt.Println("I didn't delete any data")
}
```