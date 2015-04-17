package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const uuidName = "uuid"

type MgoApi struct {
	dbName  string
	session *mgo.Session
}

func NewMgoApi(urls, dbName string) (*MgoApi, error) {
	session, err := mgo.DialWithTimeout(urls, time.Duration(3*time.Second))
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &MgoApi{dbName, session}, nil
}

func (ma *MgoApi) Write(collection string, resource map[string]interface{}) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	_, err := coll.Upsert(bson.D{{uuidName, resource[uuidName]}}, resource)

	return err
}

func (ma *MgoApi) Read(collection string, uuid string) (bool, interface{}) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	var resource map[string]interface{}
	if err := coll.Find(bson.M{uuidName: uuid}).One(&resource); err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		panic(err)
	}

	return true, resource
}
