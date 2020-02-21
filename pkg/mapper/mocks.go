package mapper

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockBody struct {
	mock.Mock
	Body io.Reader
}

func (m *MockBody) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBody) Read(p []byte) (n int, err error) {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return 0, err
	}

	return m.Body.Read(p)
}
