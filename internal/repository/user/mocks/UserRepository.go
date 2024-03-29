// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, login
func (_m *UserRepository) Find(ctx context.Context, login string) (*user_model.User, error) {
	ret := _m.Called(ctx, login)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *user_model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user_model.User, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user_model.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user_model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
