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

func (m *Movie) FindAll() (movies []Movie, err error) {
	err = db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return
}

func (m *Movie) FindById(id string) (movie Movie, err error) {
	err = db.C(COLLECTION).Find(bson.ObjectIdHex(id)).One(&movie)
	return
}

func (m *Movie) Insert(movie Movie) (err error) {
	err = db.C(COLLECTION).Insert(&movie)
	return
}

func (m *Movie) Delete(movie Movie) (err error) {
	err = db.C(COLLECTION).Remove(&movie)
	return
}

func (m *Movie) Update(movie Movie) (err error) {
	err = db.C(COLLECTION).UpdateId(movie.ID, &movie)
	return
}
