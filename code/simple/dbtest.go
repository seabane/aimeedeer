package main

import (
        "fmt"
        "launchpad.net/gobson/bson"
        "launchpad.net/mgo"
)

type Person struct {
        Name string
        Phone string
}

func main() {
        session, err := mgo.Mongo("localhost")
        if err != nil {
                panic(err)
        }
        defer session.Close()

        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        c := session.DB("test").C("people")
        err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	               &Person{"Cla", "+55 53 8402 8510"})
        if err != nil {
                panic(err)
        }

        result := Person{}
        err = c.Find(bson.M{"name": "Ale"}).One(&result)
        if err != nil {
                panic(err)
        }

        fmt.Println("Phone:", result.Phone)
}