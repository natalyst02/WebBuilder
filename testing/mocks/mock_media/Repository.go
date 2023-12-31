// Code generated by mockery v2.15.0. DO NOT EDIT.

package mock_media

import (
	media "appota/web-builder/media"

	mock "github.com/stretchr/testify/mock"

	mongo "go.mongodb.org/mongo-driver/mongo"

	options "go.mongodb.org/mongo-driver/mongo/options"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// DeleteMediaData provides a mock function with given fields: filter
func (_m *Repository) DeleteMediaData(filter primitive.D) (*mongo.DeleteResult, error) {
	ret := _m.Called(filter)

	var r0 *mongo.DeleteResult
	if rf, ok := ret.Get(0).(func(primitive.D) *mongo.DeleteResult); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.DeleteResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.D) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMediaData provides a mock function with given fields: filter, opts
func (_m *Repository) FindMediaData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error) {
	ret := _m.Called(filter, opts)

	var r0 *mongo.Cursor
	if rf, ok := ret.Get(0).(func(primitive.D, options.FindOptions) *mongo.Cursor); ok {
		r0 = rf(filter, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Cursor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.D, options.FindOptions) error); ok {
		r1 = rf(filter, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMedia provides a mock function with given fields: path
func (_m *Repository) GetMedia(path string) ([]byte, string, error) {
	ret := _m.Called(path)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(path)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// InsertMediaData provides a mock function with given fields: data
func (_m *Repository) InsertMediaData(data *media.Media) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*media.Media) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMedia provides a mock function with given fields: file, filename
func (_m *Repository) SaveMedia(file []byte, filename string) error {
	ret := _m.Called(file, filename)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, string) error); ok {
		r0 = rf(file, filename)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMediaData provides a mock function with given fields: filter, data
func (_m *Repository) UpdateMediaData(filter primitive.D, data *media.UpdateMedia) (*mongo.UpdateResult, error) {
	ret := _m.Called(filter, data)

	var r0 *mongo.UpdateResult
	if rf, ok := ret.Get(0).(func(primitive.D, *media.UpdateMedia) *mongo.UpdateResult); ok {
		r0 = rf(filter, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.D, *media.UpdateMedia) error); ok {
		r1 = rf(filter, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
