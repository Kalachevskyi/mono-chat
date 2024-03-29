// Code generated by MockGen. DO NOT EDIT.
// Source: ./mapping.go

// Package usecases is a generated GoMock package.
package usecases_test

import (
	"reflect"

	"github.com/Kalachevskyi/mono-chat/app/model"
	"github.com/golang/mock/gomock"
)

// MockMappingRepo is a mock of MappingRepo interface
type MockMappingRepo struct {
	ctrl     *gomock.Controller
	recorder *MockMappingRepoMockRecorder
}

// MockMappingRepoMockRecorder is the mock recorder for MockMappingRepo
type MockMappingRepoMockRecorder struct {
	mock *MockMappingRepo
}

// NewMockMappingRepo creates a new mock instance
func NewMockMappingRepo(ctrl *gomock.Controller) *MockMappingRepo {
	mock := &MockMappingRepo{ctrl: ctrl}
	mock.recorder = &MockMappingRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMappingRepo) EXPECT() *MockMappingRepoMockRecorder {
	return m.recorder
}

// Set mocks base method
func (m *MockMappingRepo) Set(key string, val map[string]model.CategoryMapping) error {
	ret := m.ctrl.Call(m, "Set", key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockMappingRepoMockRecorder) Set(key, val interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockMappingRepo)(nil).Set), key, val)
}

// CheckUser mocks base method
func (m *MockMappingRepo) Get(key string) (map[string]model.CategoryMapping, error) {
	ret := m.ctrl.Call(m, "CheckUser", key)
	ret0, _ := ret[0].(map[string]model.CategoryMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser
func (mr *MockMappingRepoMockRecorder) Get(key interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockMappingRepo)(nil).Get), key)
}
