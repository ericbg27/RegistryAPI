// Code generated by MockGen. DO NOT EDIT.
// Source: db/dbconnector.go

// Package mockdb is a generated GoMock package.
package mockdb

import (
	reflect "reflect"

	db "github.com/ericbg27/RegistryAPI/db"
	gomock "github.com/golang/mock/gomock"
)

// MockDBConnector is a mock of DBConnector interface.
type MockDBConnector struct {
	ctrl     *gomock.Controller
	recorder *MockDBConnectorMockRecorder
}

// MockDBConnectorMockRecorder is the mock recorder for MockDBConnector.
type MockDBConnectorMockRecorder struct {
	mock *MockDBConnector
}

// NewMockDBConnector creates a new mock instance.
func NewMockDBConnector(ctrl *gomock.Controller) *MockDBConnector {
	mock := &MockDBConnector{ctrl: ctrl}
	mock.recorder = &MockDBConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBConnector) EXPECT() *MockDBConnectorMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockDBConnector) CreateUser(userParams db.CreateUserParams) (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", userParams)
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockDBConnectorMockRecorder) CreateUser(userParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDBConnector)(nil).CreateUser), userParams)
}