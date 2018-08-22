package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Book struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Title   string        `bson:"title" json:"title"`
	Authors []string      `bson:"authors" json:"authors"`
	Price   string        `bson:"price" json:"price"`
}

const (
	bookCollection = "books"
)

func (b *Book) FindAll() (books []Book, err error) {
	err = db.C(bookCollection).Find(bson.M{}).All(&books)
	return
}

func (b *Book) FindById(id string) (book Book, err error) {
	err = db.C(bookCollection).Find(bson.ObjectIdHex(id)).One(&book)
	return
}

func (b *Book) Insert(book Book) (err error) {
	err = db.C(bookCollection).Insert(&book)
	return
}

func (b *Book) Delete(book Book) (err error) {
	err = db.C(bookCollection).Remove(&book)
	return
}

func (m *Book) Update(book Book) (err error) {
	err = db.C(bookCollection).UpdateId(book.ID, &book)
	return
}
