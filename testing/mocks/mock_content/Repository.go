// Code generated by mockery v2.15.0. DO NOT EDIT.

package mock_content

import (
	content "appota/web-builder/content"

	mock "github.com/stretchr/testify/mock"

	mongo "go.mongodb.org/mongo-driver/mongo"

	options "go.mongodb.org/mongo-driver/mongo/options"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// DeleteJSONFile provides a mock function with given fields: id, path
func (_m *Repository) DeleteJSONFile(id string, path string) error {
	ret := _m.Called(id, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTemplatesData provides a mock function with given fields: id
func (_m *Repository) DeleteTemplatesData(id string) (*mongo.DeleteResult, error) {
	ret := _m.Called(id)

	var r0 *mongo.DeleteResult
	if rf, ok := ret.Get(0).(func(string) *mongo.DeleteResult); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.DeleteResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTemplatesData provides a mock function with given fields: filter, opts
func (_m *Repository) FindTemplatesData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error) {
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

// GetJSONFile provides a mock function with given fields: id, path
func (_m *Repository) GetJSONFile(id string, path string) ([]byte, error) {
	ret := _m.Called(id, path)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string) []byte); ok {
		r0 = rf(id, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTemplatesData provides a mock function with given fields: data
func (_m *Repository) InsertTemplatesData(data *content.Template) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*content.Template) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveJSONFile provides a mock function with given fields: name, file, path
func (_m *Repository) SaveJSONFile(name string, file []byte, path string) error {
	ret := _m.Called(name, file, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, string) error); ok {
		r0 = rf(name, file, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTemplatesData provides a mock function with given fields: id, data
func (_m *Repository) UpdateTemplatesData(id string, data *content.Template) (*content.Template, error) {
	ret := _m.Called(id, data)

	var r0 *content.Template
	if rf, ok := ret.Get(0).(func(string, *content.Template) *content.Template); ok {
		r0 = rf(id, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*content.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *content.Template) error); ok {
		r1 = rf(id, data)
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
