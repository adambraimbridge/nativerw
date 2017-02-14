package db

import (
	"net"
	"strings"
	"time"

	"github.com/Financial-Times/nativerw/config"
	"github.com/Financial-Times/nativerw/mapper"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const uuidName = "uuid"

type mongoDb struct {
	dbName      string
	session     *mgo.Session
	collections map[string]bool
}

// DB contains all mongo request logic, including reads, writes and deletes.
type DB interface {
	EnsureIndex()
	GetSupportedCollections() map[string]bool
	Delete(collection string, uuidString string) error
	Ids(collection string, stopChan chan struct{}, errChan chan error) chan string
	Write(collection string, resource mapper.Resource) error
	Read(collection string, uuidString string) (res mapper.Resource, found bool, err error)
	Close()
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

// NewDBConnection dials the mongo cluster, and returns a new handler DB instance
func NewDBConnection(config *config.Configuration) (DB, error) {
	info := mgo.DialInfo{
		Timeout:    5 * time.Second,
		Addrs:      strings.Split(config.Mongos, ","),
		DialServer: tcpDialServer,
	}

	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Strong, true)
	collections := createMapWithAllowedCollections(config.Collections)

	return &mongoDb{config.DbName, session, collections}, nil
}

func (ma *mongoDb) GetSupportedCollections() map[string]bool {
	return ma.collections
}

func (ma *mongoDb) Close() {
	ma.session.Close()
}

func createMapWithAllowedCollections(collections []string) map[string]bool {
	var collectionMap = make(map[string]bool)
	for _, coll := range collections {
		collectionMap[coll] = true
	}
	return collectionMap
}

func (ma *mongoDb) EnsureIndex() {
	newSession := ma.session.Copy()
	defer newSession.Close()

	index := mgo.Index{
		Name:       "uuid-index",
		Key:        []string{"uuid"},
		Background: true,
		Unique:     true,
	}

	for coll := range ma.collections {
		newSession.DB(ma.dbName).C(coll).EnsureIndex(index)
	}
}

func (ma *mongoDb) Delete(collection string, uuidString string) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)
	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(uuidString))}

	return coll.Remove(bson.D{{uuidName, bsonUUID}})
}

func (ma *mongoDb) Write(collection string, resource mapper.Resource) error {
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

func (ma *mongoDb) Read(collection string, uuidString string) (res mapper.Resource, found bool, err error) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	bsonUUID := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(uuidString))}

	var bsonResource map[string]interface{}

	if err = coll.Find(bson.M{uuidName: bsonUUID}).One(&bsonResource); err != nil {
		if err == mgo.ErrNotFound {
			return res, false, nil
		}
		return res, false, err
	}

	uuidData := bsonResource["uuid"].(bson.Binary).Data

	res = mapper.Resource{
		UUID:        uuid.UUID(uuidData).String(),
		Content:     bsonResource["content"],
		ContentType: bsonResource["content-type"].(string),
	}

	return res, true, nil
}

func (ma *mongoDb) Ids(collection string, stopChan chan struct{}, errChan chan error) chan string {
	ids := make(chan string)
	go func() {
		defer close(ids)

		newSession := ma.session.Copy()
		newSession.SetSocketTimeout(30 * time.Second)
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
