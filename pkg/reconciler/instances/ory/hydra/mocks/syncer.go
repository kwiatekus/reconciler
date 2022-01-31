// Code generated by mockery v2.9.4. DO NOT EDIT.

package mock

import (
	context "context"

	kubernetes "github.com/kyma-incubator/reconciler/pkg/reconciler/kubernetes"

	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// Syncer is an autogenerated mock type for the Syncer type
type Syncer struct {
	mock.Mock
}

// TriggerSynchronization provides a mock function with given fields: _a0, client, logger, namespace, forceSync
func (_m *Syncer) TriggerSynchronization(_a0 context.Context, client kubernetes.Client, logger *zap.SugaredLogger, namespace string, forceSync bool) error {
	ret := _m.Called(_a0, client, logger, namespace, forceSync)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, kubernetes.Client, *zap.SugaredLogger, string, bool) error); ok {
		r0 = rf(_a0, client, logger, namespace, forceSync)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}