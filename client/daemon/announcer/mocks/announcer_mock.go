// Code generated by MockGen. DO NOT EDIT.
// Source: announcer.go
//
// Generated by this command:
//
//	mockgen -destination mocks/announcer_mock.go -source announcer.go -package mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAnnouncer is a mock of Announcer interface.
type MockAnnouncer struct {
	ctrl     *gomock.Controller
	recorder *MockAnnouncerMockRecorder
}

// MockAnnouncerMockRecorder is the mock recorder for MockAnnouncer.
type MockAnnouncerMockRecorder struct {
	mock *MockAnnouncer
}

// NewMockAnnouncer creates a new mock instance.
func NewMockAnnouncer(ctrl *gomock.Controller) *MockAnnouncer {
	mock := &MockAnnouncer{ctrl: ctrl}
	mock.recorder = &MockAnnouncerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAnnouncer) EXPECT() *MockAnnouncerMockRecorder {
	return m.recorder
}

// Serve mocks base method.
func (m *MockAnnouncer) Serve() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serve")
	ret0, _ := ret[0].(error)
	return ret0
}

// Serve indicates an expected call of Serve.
func (mr *MockAnnouncerMockRecorder) Serve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serve", reflect.TypeOf((*MockAnnouncer)(nil).Serve))
}

// Stop mocks base method.
func (m *MockAnnouncer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockAnnouncerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAnnouncer)(nil).Stop))
}
