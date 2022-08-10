// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"
	zap "go.uber.org/zap"
)

// OryFinalizersHandler is an autogenerated mock type for the OryFinalizersHandler type
type OryFinalizersHandler struct {
	mock.Mock
}

// FindAndDeleteOryFinalizers provides a mock function with given fields: kubeconfigData, logger
func (_m *OryFinalizersHandler) FindAndDeleteOryFinalizers(kubeconfigData string, logger *zap.SugaredLogger) error {
	ret := _m.Called(kubeconfigData, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *zap.SugaredLogger) error); ok {
		r0 = rf(kubeconfigData, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOryFinalizersHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewOryFinalizersHandler creates a new instance of OryFinalizersHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOryFinalizersHandler(t mockConstructorTestingTNewOryFinalizersHandler) *OryFinalizersHandler {
	mock := &OryFinalizersHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}