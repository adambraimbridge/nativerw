package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

const uuidName = "uuid"

type MgoApi struct {
	dbName      string
	session     *mgo.Session
	beforeWrite propertyConverter
	afterRead   propertyConverter
}

func NewMgoApi(urls, dbName string, beforeWrite, afterRead propertyConverter) (*MgoApi, error) {
	session, err := mgo.DialWithTimeout(urls, time.Duration(3*time.Second))
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &MgoApi{dbName, session, beforeWrite, afterRead}, nil
}

func (ma *MgoApi) Write(collection string, resource map[string]interface{}) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	ma.mongoizeAll(resource)

	_, err := coll.Upsert(bson.D{{uuidName, resource[uuidName]}}, resource)

	return err
}

func (ma *MgoApi) Read(collection string, resourceId string) (bool, interface{}) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	// convert resource id to mgo friendly form if needed
	props := make(map[string]interface{})
	props[uuidName] = resourceId
	ma.mongoizeAll(props)
	uuid := props[uuidName]

	var resource map[string]interface{}
	if err := coll.Find(bson.M{uuidName: uuid}).One(&resource); err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		panic(err)
	}

	ma.unmongoizeAll(resource)

	return true, resource
}

func (ma *MgoApi) mongoizeAll(resource map[string]interface{}) {
	for k, v := range resource {
		if reflect.ValueOf(v).Type() == mapStrIfType {
			ma.mongoizeAll(v.(map[string]interface{}))
		} else {
			pm := simplePropertyModifier{resource, k}
			ma.beforeWrite(pm, resource, k, v)
		}
	}
}

var mapStrIfType = reflect.ValueOf(make(map[string]interface{})).Type()

func (ma *MgoApi) unmongoizeAll(resource map[string]interface{}) {
	for k, v := range resource {
		if reflect.ValueOf(v).Type() == mapStrIfType {
			ma.unmongoizeAll(v.(map[string]interface{}))
		} else {
			pm := simplePropertyModifier{resource, k}
			ma.afterRead(pm, resource, k, v)
		}
	}
}
