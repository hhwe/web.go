package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	host   = "127.0.0.1:27017"
	source = "admin"
	user   = "user"
	pass   = "123456"
	database = "web"
)

var (
	globalSession *mgo.Session
)

// initial global session
func init() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		//Source:   source,
		//Username: user,
		//Password: pass,
		//Database:database,
		Timeout:time.Second * 1,
		PoolLimit:2,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ", err)
	}

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	globalSession = session
}

// build index of mongodb
func ensureIndex(s *mgo.Session) {
	// ensure index of collections which name is user
	c := s.DB(database).C("user")
	index := mgo.Index{
		Key:        []string{"email", "phone", "username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
