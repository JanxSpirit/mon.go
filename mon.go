package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

var mongoHost string
var mongoPort int
var mongoDb string
var mongoColl string
var cmd string
var mainSession *mgo.Session
var err error

func init() {
	flag.StringVar(&mongoHost, "host", "localhost", "the host running the MongoDB instance")
	flag.IntVar(&mongoPort, "port", 27017, "the port the MongoDB instance is listening on")
	flag.StringVar(&mongoDb, "db", "", "the db to use")
	flag.StringVar(&mongoColl, "coll", "", "the collection to use")
	flag.Parse()
	cmd = os.Args[len(os.Args)-1]
	mainSession, err = mgo.Dial(fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort))
	if err != nil {
		panic(err)
	}
}

func main() {
	session := getSession()
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	switch cmd {
	case "colls":
		listColls(session)
	case "dbs":
		listDbs(session)
	case "last":
		lastDocument(session)
	}

}

func listColls(session *mgo.Session) {
	db := session.DB(mongoDb)

	cNames, err := db.CollectionNames()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(cNames); i++ {
		fmt.Println(cNames[i])
	}
}

func listDbs(session *mgo.Session) {
	dbNames, err := session.DatabaseNames()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(dbNames); i++ {
		fmt.Println(dbNames[i])
	}
}

func lastDocument(session *mgo.Session) {
	var m bson.M
	coll := session.DB(mongoDb).C(mongoColl)
	coll.Find(nil).Sort("-_id").One(&m)
	js, _ := json.MarshalIndent(&m, "", "  ")
	fmt.Println(string(js))
}

func getSession() (session *mgo.Session) {
	return mainSession.Copy()
}
