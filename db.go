package main

import (
	"log"

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
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ", err)
	}
	//defer session.Close() // clean up when we’re done

	globalSession = session
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)
}

// build index of mongodb
func ensureIndex(s *mgo.Session) {
	s = globalSession.Copy()
	defer s.Close()

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

func Insert(s *mgo.Session, collection string, docs ...interface{}) error {
	c := s.DB(database).C(collection)
	return c.Insert(docs...)
}

func FindOne(s *mgo.Session, collection string, query, selector, result interface{}) error {
	c := s.DB(database).C(collection)
	return c.Find(query).Select(selector).One(result)
}

func FindAll(s *mgo.Session, collection string, query, selector, result interface{}) error {
	c := s.DB(database).C(collection)
	return c.Find(query).Select(selector).All(result)
}

func Update(s *mgo.Session, collection string, query, update interface{}) error {
	c := s.DB(database).C(collection)
	return c.Update(query, update)
}

func Remove(s *mgo.Session, collection string, query interface{}) error {
	c := s.DB(database).C(collection)
	return c.Remove(query)
}


