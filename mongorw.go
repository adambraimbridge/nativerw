package main

import (
	"net"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"log"
)

const uuidName = "uuid"

type resource struct {
	UUID        string
	Content     interface{}
	ContentType string
}

type mgoAPI struct {
	dbName      string
	session     *mgo.Session
	collections map[string]bool
}

func tcpDialServer(addr *mgo.ServerAddr) (net.Conn, error) {
	ra, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, ra)
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(30 * time.Second)
	return conn, nil
}

func newMgoAPI(config *configuration) (*mgoAPI, error) {
	info := mgo.DialInfo{
		Timeout:    30 * time.Second,
		Addrs:      strings.Split(config.Mongos, ","),
		DialServer: tcpDialServer,
	}
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Strong, true)
	collections := createMapWithAllowedCollections(config.Collections)

	return &mgoAPI{config.DbName, session, collections}, nil
}

func createMapWithAllowedCollections(collections []string) map[string]bool {
	var collectionMap = make(map[string]bool)
	for _, coll := range collections {
		collectionMap[coll] = true
	}
	return collectionMap
}

func (ma *mgoAPI) EnsureIndex() {
	newSession := ma.session.Copy()
	defer newSession.Close()

	index := mgo.Index{
		Key:        []string{"uuid"},
		Background: true,
	}

	for coll := range ma.collections {
		newSession.DB(ma.dbName).C(coll).EnsureIndex(index)
	}
}

func (ma *mgoAPI) Delete(collection string, uuidString string) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)
	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(uuidString))}

	return coll.Remove(bson.D{{uuidName, bsonUUID}})
}

func (ma *mgoAPI) Write(collection string, resource resource) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	session := reflect.ValueOf(newSession)
	syncTimeout := session.FieldByName("syncTimeout")
	sockTimeout := session.FieldByName("sockTimeout")
	log.Printf("syncTimeout=%v", syncTimeout)
	log.Printf("sockTimeout=%v", sockTimeout)

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

func (ma *mgoAPI) Read(collection string, uuidString string) (found bool, res resource, err error) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(uuidString))}

	var bsonResource map[string]interface{}

	if err = coll.Find(bson.M{uuidName: bsonUUID}).One(&bsonResource); err != nil {
		if err == mgo.ErrNotFound {
			return false, res, nil
		}
		return false, res, err
	}

	uuidData := bsonResource["uuid"].(bson.Binary).Data

	res = resource{
		UUID:        uuid.UUID(uuidData).String(),
		Content:     bsonResource["content"],
		ContentType: bsonResource["content-type"].(string),
	}

	return true, res, nil
}

func (ma *mgoAPI) Ids(collection string, stopChan chan struct{}, errChan chan error) chan string {
	ids := make(chan string)
	go func() {
		defer close(ids)
		newSession := ma.session.Copy()
		defer newSession.Close()
		coll := newSession.DB(ma.dbName).C(collection)

		iter := coll.Find(nil).Select(bson.M{uuidName: true}).Iter()
		var result map[string]interface{}
		for iter.Next(&result) {
			select {
			case <-stopChan:
				break
			case ids <- uuid.UUID(result["uuid"].(bson.Binary).Data).String():
			}
		}
		if err := iter.Close(); err != nil {
			errChan <- err
		}
	}()
	return ids
}
