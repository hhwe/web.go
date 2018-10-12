package main

import "gopkg.in/mgo.v2"

// build index of mongodb
func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("web").C("user")

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
