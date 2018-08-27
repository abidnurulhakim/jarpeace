package database

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

type MongoDB struct {
	Session  *mgo.Session
	Database string
}

// Option holds all necessary options for database.
type MongoDBOption struct {
	User     string
	Password string
	Host     string
	Database string
	Timeout  int
}

func NewMongoDB(opt *MongoDBOption) (*MongoDB, error) {
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    []string{opt.Host},
		Database: opt.Database,
		Username: opt.User,
		Password: opt.Password,
		Timeout:  opt.Timeout * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDialInfo)
	if err != nil {
		panic(err)
	}
	return &MongoDB{Session: session, Database: opt.Database}, err
}

func (db *MongoDB) Collection(collection string) *mgo.Collection {
	c := db.Session.DB(db.Database).C(collection)
	return c
}

func (db *MongoDB) Close() {
	db.Session.Close()
}

func (db *MongoDB) Copy() *MongoDB {
	return &MongoDB{Session: db.Session.Copy(), Database: db.Database}
}
