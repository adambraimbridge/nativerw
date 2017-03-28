package db

import (
	"errors"
	"testing"
	"time"

	"github.com/Financial-Times/nativerw/mapper"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func generateResource() mapper.Resource {
	return mapper.Resource{
		UUID:        uuid.NewUUID().String(),
		Content:     map[string]interface{}{"randomness": uuid.NewUUID().String()},
		ContentType: "application/json",
	}
}

func TestReadWriteDelete(t *testing.T) {
	mongo := startMongo(t)
	connection, err := mongo.Open()

	assert.NoError(t, err)

	defer connection.Close()

	expectedResource := generateResource()

	err = connection.Write("methode", expectedResource)
	assert.NoError(t, err)

	res, found, err := connection.Read("methode", expectedResource.UUID)

	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, expectedResource.ContentType, res.ContentType)
	assert.Equal(t, expectedResource.UUID, res.UUID)
	assert.Equal(t, expectedResource.Content, res.Content)

	err = connection.Delete("methode", expectedResource.UUID)
	assert.NoError(t, err)

	res, found, err = connection.Read("methode", expectedResource.UUID)
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestGetSupportedCollections(t *testing.T) {
	mongo := startMongo(t)
	connection, err := mongo.Open()
	assert.NoError(t, err)

	defer connection.Close()

	expected := map[string]bool{"methode": true} // this is set in mongo_test.go
	actual := connection.GetSupportedCollections()
	assert.Equal(t, expected, actual)
}

func TestEnsureIndexes(t *testing.T) {
	mongo := startMongo(t)
	connection, err := mongo.Open()
	assert.NoError(t, err)

	defer connection.Close()

	connection.EnsureIndex()
	indexes, err := connection.(*mongoConnection).session.DB("native-store").C("methode").Indexes()

	assert.NoError(t, err)
	count := 0
	for _, index := range indexes {
		if index.Name == "uuid-index" {
			assert.True(t, index.Background)
			assert.True(t, index.Unique)
			assert.Equal(t, []string{"uuid"}, index.Key)
			count = count + 1
		}
	}

	assert.Equal(t, 1, count)
}

func TestAwaitConnectionFailsIfNotOpened(t *testing.T) {
	mongo := startMongo(t)

	_, err := mongo.Await()
	assert.Error(t, err)
}

func TestAwaitConnectionBlocks(t *testing.T) {
	mongo := startMongo(t)

	mongo.(*mongoDB).connection = NewOptional(func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return &mongoConnection{dbName: "test-collection"}, nil
	})

	connection, err := mongo.Await()
	assert.NoError(t, err)
	assert.Equal(t, "test-collection", connection.(*mongoConnection).dbName)
}

func TestAwaitConnectionFails(t *testing.T) {
	mongo := startMongo(t)

	mongo.(*mongoDB).connection = NewOptional(func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return nil, errors.New("went spectacularly wrong")
	})

	connection, err := mongo.Await()
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestAwaitConnectionReturnsIfInitialised(t *testing.T) {
	mongo := startMongo(t).(*mongoDB)
	ch := make(chan bool, 1)
	defer func() {
		ch <- true
		close(ch)
	}()

	mongo.connection = NewOptional(func() (interface{}, error) {
		<-ch
		return &mongoConnection{dbName: "psych!-changed-it"}, nil
	})

	mongo.connection.val = &mongoConnection{dbName: "find-me-pls"}

	connection, err := mongo.Await()
	assert.NoError(t, err)
	assert.Equal(t, "find-me-pls", connection.(*mongoConnection).dbName)
}

func TestOpenWillReturnConnectionIfAlreadyInitialised(t *testing.T) {
	mongo := startMongo(t).(*mongoDB)

	mongo.connection = NewOptional(func() (interface{}, error) {
		return &mongoConnection{dbName: "faked"}, nil
	})

	mongo.connection.Block()

	connection, err := mongo.Open()
	assert.NoError(t, err)
	assert.Equal(t, "faked", connection.(*mongoConnection).dbName)
}

func TestOpenFailsIfInitialisationFailed(t *testing.T) {
	mongo := startMongo(t).(*mongoDB)

	mongo.connection = NewOptional(func() (interface{}, error) {
		return nil, errors.New("i failed")
	})

	mongo.connection.Block()

	connection, err := mongo.Open()
	assert.Error(t, err)
	assert.Nil(t, connection)
}
