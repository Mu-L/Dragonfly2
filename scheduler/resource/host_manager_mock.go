// Code generated by MockGen. DO NOT EDIT.
// Source: host_manager.go
//
// Generated by this command:
//
//	mockgen -destination host_manager_mock.go -source host_manager.go -package resource
//
// Package resource is a generated GoMock package.
package resource

import (
	reflect "reflect"

	set "d7y.io/dragonfly/v2/pkg/container/set"
	gomock "go.uber.org/mock/gomock"
)

// MockHostManager is a mock of HostManager interface.
type MockHostManager struct {
	ctrl     *gomock.Controller
	recorder *MockHostManagerMockRecorder
}

// MockHostManagerMockRecorder is the mock recorder for MockHostManager.
type MockHostManagerMockRecorder struct {
	mock *MockHostManager
}

// NewMockHostManager creates a new mock instance.
func NewMockHostManager(ctrl *gomock.Controller) *MockHostManager {
	mock := &MockHostManager{ctrl: ctrl}
	mock.recorder = &MockHostManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHostManager) EXPECT() *MockHostManagerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockHostManager) Delete(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0)
}

// Delete indicates an expected call of Delete.
func (mr *MockHostManagerMockRecorder) Delete(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockHostManager)(nil).Delete), arg0)
}

// Load mocks base method.
func (m *MockHostManager) Load(arg0 string) (*Host, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].(*Host)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockHostManagerMockRecorder) Load(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockHostManager)(nil).Load), arg0)
}

// LoadOrStore mocks base method.
func (m *MockHostManager) LoadOrStore(arg0 *Host) (*Host, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadOrStore", arg0)
	ret0, _ := ret[0].(*Host)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// LoadOrStore indicates an expected call of LoadOrStore.
func (mr *MockHostManagerMockRecorder) LoadOrStore(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadOrStore", reflect.TypeOf((*MockHostManager)(nil).LoadOrStore), arg0)
}

// LoadRandomHosts mocks base method.
func (m *MockHostManager) LoadRandomHosts(arg0 int, arg1 set.SafeSet[string]) []*Host {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadRandomHosts", arg0, arg1)
	ret0, _ := ret[0].([]*Host)
	return ret0
}

// LoadRandomHosts indicates an expected call of LoadRandomHosts.
func (mr *MockHostManagerMockRecorder) LoadRandomHosts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadRandomHosts", reflect.TypeOf((*MockHostManager)(nil).LoadRandomHosts), arg0, arg1)
}

// Range mocks base method.
func (m *MockHostManager) Range(f func(any, any) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Range", f)
}

// Range indicates an expected call of Range.
func (mr *MockHostManagerMockRecorder) Range(f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Range", reflect.TypeOf((*MockHostManager)(nil).Range), f)
}

// RunGC mocks base method.
func (m *MockHostManager) RunGC() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunGC")
	ret0, _ := ret[0].(error)
	return ret0
}

// RunGC indicates an expected call of RunGC.
func (mr *MockHostManagerMockRecorder) RunGC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunGC", reflect.TypeOf((*MockHostManager)(nil).RunGC))
}

// Store mocks base method.
func (m *MockHostManager) Store(arg0 *Host) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", arg0)
}

// Store indicates an expected call of Store.
func (mr *MockHostManagerMockRecorder) Store(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockHostManager)(nil).Store), arg0)
}
