package resources

import (
	"github.com/Financial-Times/nativerw/mapper"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	CallArgs []interface{}
}

func (m *MockDB) EnsureIndex() {
	m.Called()
}

func (m *MockDB) GetSupportedCollections() map[string]bool {
	args := m.Called()
	return args.Get(0).(map[string]bool)
}

func (m *MockDB) Close() {
	m.Called()
}

func (m *MockDB) Delete(collection string, uuidString string) error {
	args := m.Called(collection, uuidString)
	return args.Error(0)
}

func (m *MockDB) Ids(collection string, stopChan chan struct{}, errChan chan error) chan string {
	args := m.Called(collection, stopChan, errChan)
	m.CallArgs = []interface{}{collection, stopChan, errChan}
	return args.Get(0).(chan string)
}

func (m *MockDB) Write(collection string, resource mapper.Resource) error {
	args := m.Called(collection, resource)
	return args.Error(0)
}

func (m *MockDB) Read(collection string, uuidString string) (found bool, res mapper.Resource, err error) {
	args := m.Called(collection, uuidString)
	return args.Bool(0), args.Get(1).(mapper.Resource), args.Error(2)
}
