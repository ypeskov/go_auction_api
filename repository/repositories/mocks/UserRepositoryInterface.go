// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	models "ypeskov/go_hillel_9/repository/models"

	mock "github.com/stretchr/testify/mock"
)

// UserRepositoryInterface is an autogenerated mock type for the UserRepositoryInterface type
type UserRepositoryInterface struct {
	mock.Mock
}

// AddOrUpdateRefreshToken provides a mock function with given fields: userId, token
func (_m *UserRepositoryInterface) AddOrUpdateRefreshToken(userId int, token string) error {
	ret := _m.Called(userId, token)

	if len(ret) == 0 {
		panic("no return value specified for AddOrUpdateRefreshToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(userId, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: srcUser
func (_m *UserRepositoryInterface) CreateUser(srcUser *models.User) (*models.User, error) {
	ret := _m.Called(srcUser)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.User) (*models.User, error)); ok {
		return rf(srcUser)
	}
	if rf, ok := ret.Get(0).(func(*models.User) *models.User); ok {
		r0 = rf(srcUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(srcUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *UserRepositoryInterface) GetUserByEmail(email string) *models.User {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	return r0
}

// GetUserByRefreshToken provides a mock function with given fields: token
func (_m *UserRepositoryInterface) GetUserByRefreshToken(token string) *models.User {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByRefreshToken")
	}

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	return r0
}

// GetUsersList provides a mock function with given fields:
func (_m *UserRepositoryInterface) GetUsersList() ([]*models.User, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetUsersList")
	}

	var r0 []*models.User
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*models.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*models.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRepositoryInterface creates a new instance of UserRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepositoryInterface {
	mock := &UserRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
