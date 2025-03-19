// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/vishenosik/sso/internal/services/authentication/models"
	mock "github.com/stretchr/testify/mock"
)

// Authentication is an autogenerated mock type for the Authentication type
type Authentication struct {
	mock.Mock
}

// IsAdmin provides a mock function with given fields: ctx, userID
func (_m *Authentication) IsAdmin(ctx context.Context, userID string) (bool, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for IsAdmin")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, request, appID
func (_m *Authentication) Login(ctx context.Context, request models.LoginRequest, appID string) (string, error) {
	ret := _m.Called(ctx, request, appID)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginRequest, string) (string, error)); ok {
		return rf(ctx, request, appID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginRequest, string) string); ok {
		r0 = rf(ctx, request, appID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.LoginRequest, string) error); ok {
		r1 = rf(ctx, request, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterNewUser provides a mock function with given fields: ctx, request
func (_m *Authentication) RegisterNewUser(ctx context.Context, request models.RegisterRequest) (string, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for RegisterNewUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterRequest) (string, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.RegisterRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthentication creates a new instance of Authentication. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthentication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Authentication {
	mock := &Authentication{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
