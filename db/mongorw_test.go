package db

import (
	"testing"

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
	defer mongo.Close()

	expectedResource := generateResource()

	err := mongo.Write("methode", expectedResource)
	assert.NoError(t, err)

	res, found, err := mongo.Read("methode", expectedResource.UUID)

	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, expectedResource.ContentType, res.ContentType)
	assert.Equal(t, expectedResource.UUID, res.UUID)
	assert.Equal(t, expectedResource.Content, res.Content)

	err = mongo.Delete("methode", expectedResource.UUID)
	assert.NoError(t, err)

	res, found, err = mongo.Read("methode", expectedResource.UUID)
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestGetSupportedCollections(t *testing.T) {
	mongo := startMongo(t)
	defer mongo.Close()

	expected := map[string]bool{"methode": true} // this is set in mongo_test.go
	actual := mongo.GetSupportedCollections()
	assert.Equal(t, expected, actual)
}

func TestEnsureIndexes(t *testing.T) {
	mongo := startMongo(t)
	defer mongo.Close()

	mongo.EnsureIndex()
	indexes, err := mongo.(*mongoDb).session.DB("native-store").C("methode").Indexes()

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

func TestGetIds(t *testing.T) {
	mongo := startMongo(t)
	defer mongo.Close()

	expectedResource := generateResource()
	mongo.Write("methode", expectedResource)

	stop := make(chan struct{})
	err := make(chan error)

	ids := mongo.Ids("methode", stop, err)

	for {
		actualID, ok := <-ids
		if !ok {
			t.Fail() // Should've found the expected ID by now
			break
		}

		if actualID == expectedResource.UUID {
			break
		}
	}
}

func TestGetIdsShouldStop(t *testing.T) {
	mongo := startMongo(t)
	defer mongo.Close()

	expectedResource := generateResource()
	mongo.Write("methode", expectedResource)

	stop := make(chan struct{})
	err := make(chan error)

	mongo.Ids("methode", stop, err)
	stop <- struct{}{} // we can see it stops by the coverage
}
