package models

import (
	"gopkg.in/mgo.v2"
)

var (
	mongodb = MongoDB{SERVER, DATABASE}
	db      = mongodb.Connect()
)

const (
	// SERVER 服务器地址
	SERVER = "127.0.0.1"
	// DATABASE 使用的数据库
	DATABASE = "test"
)

// MongoDB 定义一个mongodb数据库连接结构
type MongoDB struct {
	Server   string
	Database string
}

// Connect 建立一个mongodb链接
func (m *MongoDB) Connect() (db *mgo.Database) {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		panic(err)
	}

	db = session.DB(m.Database)
	return
}
