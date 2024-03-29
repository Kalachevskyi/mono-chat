// Code generated by MockGen. DO NOT EDIT.
// Source: ./transaction.go

// Package usecases is a generated GoMock package.
package usecases_test

import (
	reflect "reflect"
	time "time"

	model "github.com/Kalachevskyi/mono-chat/app/model"
	gomock "github.com/golang/mock/gomock"
)

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Error mocks base method
func (m *MockLogger) Error(args ...interface{}) {
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error
func (mr *MockLoggerMockRecorder) Error(args ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), args...)
}

// MockMonoRepo is a mock of MonoRepo interface
type MockMonoRepo struct {
	ctrl     *gomock.Controller
	recorder *MockMonoRepoMockRecorder
}

// MockMonoRepoMockRecorder is the mock recorder for MockMonoRepo
type MockMonoRepoMockRecorder struct {
	mock *MockMonoRepo
}

// NewMockMonoRepo creates a new mock instance
func NewMockMonoRepo(ctrl *gomock.Controller) *MockMonoRepo {
	mock := &MockMonoRepo{ctrl: ctrl}
	mock.recorder = &MockMonoRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMonoRepo) EXPECT() *MockMonoRepoMockRecorder {
	return m.recorder
}

// GetTransactions mocks base method
func (m *MockMonoRepo) GetTransactions(token, account string, from, to time.Time) ([]model.Transaction, error) {
	ret := m.ctrl.Call(m, "GetTransactions", token, account, from, to)
	ret0, _ := ret[0].([]model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions
func (mr *MockMonoRepoMockRecorder) GetTransactions(token, account, from, to interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockMonoRepo)(nil).GetTransactions), token, account, from, to)
}
