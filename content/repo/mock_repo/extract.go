// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/urandom/readeef/content/repo (interfaces: Extract)

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	gomock "github.com/golang/mock/gomock"
	content "github.com/urandom/readeef/content"
	reflect "reflect"
)

// MockExtract is a mock of Extract interface
type MockExtract struct {
	ctrl     *gomock.Controller
	recorder *MockExtractMockRecorder
}

// MockExtractMockRecorder is the mock recorder for MockExtract
type MockExtractMockRecorder struct {
	mock *MockExtract
}

// NewMockExtract creates a new mock instance
func NewMockExtract(ctrl *gomock.Controller) *MockExtract {
	mock := &MockExtract{ctrl: ctrl}
	mock.recorder = &MockExtractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExtract) EXPECT() *MockExtractMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockExtract) Get(arg0 content.Article) (content.Extract, error) {
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(content.Extract)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockExtractMockRecorder) Get(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockExtract)(nil).Get), arg0)
}

// Update mocks base method
func (m *MockExtract) Update(arg0 content.Extract) error {
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockExtractMockRecorder) Update(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockExtract)(nil).Update), arg0)
}
