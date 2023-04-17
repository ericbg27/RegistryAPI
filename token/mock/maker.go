// Code generated by MockGen. DO NOT EDIT.
// Source: token/maker.go

// Package mocktoken is a generated GoMock package.
package mocktoken

import (
	reflect "reflect"
	time "time"

	token "github.com/ericbg27/RegistryAPI/token"
	gomock "github.com/golang/mock/gomock"
)

// MockMaker is a mock of Maker interface.
type MockMaker struct {
	ctrl     *gomock.Controller
	recorder *MockMakerMockRecorder
}

// MockMakerMockRecorder is the mock recorder for MockMaker.
type MockMakerMockRecorder struct {
	mock *MockMaker
}

// NewMockMaker creates a new mock instance.
func NewMockMaker(ctrl *gomock.Controller) *MockMaker {
	mock := &MockMaker{ctrl: ctrl}
	mock.recorder = &MockMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMaker) EXPECT() *MockMakerMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockMaker) CreateToken(username string, duration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", username, duration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockMakerMockRecorder) CreateToken(username, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockMaker)(nil).CreateToken), username, duration)
}

// VerifyToken mocks base method.
func (m *MockMaker) VerifyToken(tokenToVerify string) (*token.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", tokenToVerify)
	ret0, _ := ret[0].(*token.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockMakerMockRecorder) VerifyToken(tokenToVerify interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockMaker)(nil).VerifyToken), tokenToVerify)
}