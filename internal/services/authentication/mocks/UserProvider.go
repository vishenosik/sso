// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/blacksmith-vish/sso/internal/store/models"
	mock "github.com/stretchr/testify/mock"
)

// UserProvider is an autogenerated mock type for the UserProvider type
type UserProvider struct {
	mock.Mock
}

// IsAdmin provides a mock function with given fields: ctx, userID
func (_m *UserProvider) IsAdmin(ctx context.Context, userID string) (bool, error) {
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

// User provides a mock function with given fields: ctx, email
func (_m *UserProvider) User(ctx context.Context, email string) (models.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for User")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserProvider creates a new instance of UserProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserProvider {
	mock := &UserProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
