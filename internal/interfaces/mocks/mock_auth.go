// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lenarlenar/gomart/internal/interfaces (interfaces: AuthService)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_auth.go . AuthService
//

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	reflect "reflect"

	models "github.com/lenarlenar/gomart/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
	isgomock struct{}
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthService) Login(user models.UserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceMockRecorder) Login(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), user)
}

// Register mocks base method.
func (m *MockAuthService) Register(user models.UserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockAuthServiceMockRecorder) Register(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), user)
}
