package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Movie struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	CoverImage  string        `bson:"cover_image" json:"cover_image"`
	Description string        `bson:"description" json:"description"`
}

const (
	movieCollection = "movies"
)

func (m *Movie) FindAll() (movies []Movie, err error) {
	err = db.C(movieCollection).Find(bson.M{}).All(&movies)
	return
}

func (m *Movie) FindById(id string) (movie Movie, err error) {
	err = db.C(movieCollection).Find(bson.ObjectIdHex(id)).One(&movie)
	return
}

func (m *Movie) Insert(movie Movie) (err error) {
	err = db.C(movieCollection).Insert(&movie)
	return
}

func (m *Movie) Delete(movie Movie) (err error) {
	err = db.C(movieCollection).Remove(&movie)
	return
}

func (m *Movie) Update(movie Movie) (err error) {
	err = db.C(movieCollection).UpdateId(movie.ID, &movie)
	return
}
