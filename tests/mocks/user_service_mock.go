// Code generated by MockGen. DO NOT EDIT.
// Source: user_service.go
//
// Generated by this command:
//
//	mockgen -source=user_service.go -destination=../../tests/mocks/user_service_mock.go -package=mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockUserBandwidthControllerI is a mock of UserBandwidthControllerI interface.
type MockUserBandwidthControllerI struct {
	ctrl     *gomock.Controller
	recorder *MockUserBandwidthControllerIMockRecorder
}

// MockUserBandwidthControllerIMockRecorder is the mock recorder for MockUserBandwidthControllerI.
type MockUserBandwidthControllerIMockRecorder struct {
	mock *MockUserBandwidthControllerI
}

// NewMockUserBandwidthControllerI creates a new mock instance.
func NewMockUserBandwidthControllerI(ctrl *gomock.Controller) *MockUserBandwidthControllerI {
	mock := &MockUserBandwidthControllerI{ctrl: ctrl}
	mock.recorder = &MockUserBandwidthControllerIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserBandwidthControllerI) EXPECT() *MockUserBandwidthControllerIMockRecorder {
	return m.recorder
}

// GetAvailableBandwidth mocks base method.
func (m *MockUserBandwidthControllerI) GetAvailableBandwidth(username string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableBandwidth", username)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetAvailableBandwidth indicates an expected call of GetAvailableBandwidth.
func (mr *MockUserBandwidthControllerIMockRecorder) GetAvailableBandwidth(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableBandwidth", reflect.TypeOf((*MockUserBandwidthControllerI)(nil).GetAvailableBandwidth), username)
}

// UpdateBandwidthUsed mocks base method.
func (m *MockUserBandwidthControllerI) UpdateBandwidthUsed(username string, used int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBandwidthUsed", username, used)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBandwidthUsed indicates an expected call of UpdateBandwidthUsed.
func (mr *MockUserBandwidthControllerIMockRecorder) UpdateBandwidthUsed(username, used any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBandwidthUsed", reflect.TypeOf((*MockUserBandwidthControllerI)(nil).UpdateBandwidthUsed), username, used)
}

// MockUserServiceI is a mock of UserServiceI interface.
type MockUserServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceIMockRecorder
}

// MockUserServiceIMockRecorder is the mock recorder for MockUserServiceI.
type MockUserServiceIMockRecorder struct {
	mock *MockUserServiceI
}

// NewMockUserServiceI creates a new mock instance.
func NewMockUserServiceI(ctrl *gomock.Controller) *MockUserServiceI {
	mock := &MockUserServiceI{ctrl: ctrl}
	mock.recorder = &MockUserServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceI) EXPECT() *MockUserServiceIMockRecorder {
	return m.recorder
}

// AddToHistory mocks base method.
func (m *MockUserServiceI) AddToHistory(user, host string, time time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToHistory", user, host, time)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToHistory indicates an expected call of AddToHistory.
func (mr *MockUserServiceIMockRecorder) AddToHistory(user, host, time any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToHistory", reflect.TypeOf((*MockUserServiceI)(nil).AddToHistory), user, host, time)
}

// GetAvailableBandwidth mocks base method.
func (m *MockUserServiceI) GetAvailableBandwidth(username string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableBandwidth", username)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetAvailableBandwidth indicates an expected call of GetAvailableBandwidth.
func (mr *MockUserServiceIMockRecorder) GetAvailableBandwidth(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableBandwidth", reflect.TypeOf((*MockUserServiceI)(nil).GetAvailableBandwidth), username)
}

// GetBandwidthUsed mocks base method.
func (m *MockUserServiceI) GetBandwidthUsed(username string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBandwidthUsed", username)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetBandwidthUsed indicates an expected call of GetBandwidthUsed.
func (mr *MockUserServiceIMockRecorder) GetBandwidthUsed(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBandwidthUsed", reflect.TypeOf((*MockUserServiceI)(nil).GetBandwidthUsed), username)
}

// UpdateBandwidthUsed mocks base method.
func (m *MockUserServiceI) UpdateBandwidthUsed(username string, used int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBandwidthUsed", username, used)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBandwidthUsed indicates an expected call of UpdateBandwidthUsed.
func (mr *MockUserServiceIMockRecorder) UpdateBandwidthUsed(username, used any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBandwidthUsed", reflect.TypeOf((*MockUserServiceI)(nil).UpdateBandwidthUsed), username, used)
}
