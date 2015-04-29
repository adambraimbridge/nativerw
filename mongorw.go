package main

import (
	"code.google.com/p/go-uuid/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const uuidName = "uuid"

type Resource struct {
	UUID        string
	Content     interface{}
	ContentType string
}

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

func (ma *MgoApi) Write(collection string, resource Resource) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(resource.UUID))}
	bsonResource := map[string]interface{}{
		"uuid":         bsonUUID,
		"content":      resource.Content,
		"content-type": resource.ContentType,
	}

	_, err := coll.Upsert(bson.D{{uuidName, bsonUUID}}, bsonResource)

	return err
}

func (ma *MgoApi) Read(collection string, uuidString string) (found bool, resource Resource, err error) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(uuidString))}

	var bsonResource map[string]interface{}

	if err = coll.Find(bson.M{uuidName: bsonUUID}).One(&bsonResource); err != nil {
		if err == mgo.ErrNotFound {
			return false, resource, nil
		}
		return false, resource, err
	}

	uuidData := bsonResource["uuid"].(bson.Binary).Data

	resource = Resource{
		UUID:        uuid.UUID(uuidData).String(),
		Content:     bsonResource["content"],
		ContentType: bsonResource["content-type"].(string),
	}

	return true, resource, nil
}
