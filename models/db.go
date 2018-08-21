// db.go - 数据库连接脚本
//
// 可以通过建立一个MongoDB对象来兴建一个session并且返回
package models

import (
	"gopkg.in/mgo.v2"
)

const (
	SERVER     = "127.0.0.1"
	DATABASE   = "test"
	COLLECTION = "users"
)

// 定义一个mongodb数据库连接结构
type MongoDB struct {
	Server   string
	Database string
}

func (m *MongoDB) Connect() (db *mgo.Database) {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	db = session.DB(m.Database)
}

var mongoDB = MongoDB{SERVER, DATABASE}
var db = mongoDB.Connect()
