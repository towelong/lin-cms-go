// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/towelong/lin-cms-go/internal/service (interfaces: IUserService)

// Package mockservice is a generated GoMock package.
package mockservice

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/towelong/lin-cms-go/internal/domain/model"
	vo "github.com/towelong/lin-cms-go/internal/domain/vo"
)

// MockIUserService is a mock of IUserService interface.
type MockIUserService struct {
	ctrl     *gomock.Controller
	recorder *MockIUserServiceMockRecorder
}

// MockIUserServiceMockRecorder is the mock recorder for MockIUserService.
type MockIUserServiceMockRecorder struct {
	mock *MockIUserService
}

// NewMockIUserService creates a new mock instance.
func NewMockIUserService(ctrl *gomock.Controller) *MockIUserService {
	mock := &MockIUserService{ctrl: ctrl}
	mock.recorder = &MockIUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserService) EXPECT() *MockIUserServiceMockRecorder {
	return m.recorder
}

// GetRootUserId mocks base method.
func (m *MockIUserService) GetRootUserId() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRootUserId")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetRootUserId indicates an expected call of GetRootUserId.
func (mr *MockIUserServiceMockRecorder) GetRootUserId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRootUserId", reflect.TypeOf((*MockIUserService)(nil).GetRootUserId))
}

// GetUserById mocks base method.
func (m *MockIUserService) GetUserById(arg0 int) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockIUserServiceMockRecorder) GetUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockIUserService)(nil).GetUserById), arg0)
}

// GetUserPageByGroupId mocks base method.
func (m *MockIUserService) GetUserPageByGroupId(arg0, arg1, arg2 int) (*vo.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPageByGroupId", arg0, arg1, arg2)
	ret0, _ := ret[0].(*vo.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPageByGroupId indicates an expected call of GetUserPageByGroupId.
func (mr *MockIUserServiceMockRecorder) GetUserPageByGroupId(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPageByGroupId", reflect.TypeOf((*MockIUserService)(nil).GetUserPageByGroupId), arg0, arg1, arg2)
}

// IsAdmin mocks base method.
func (m *MockIUserService) IsAdmin(arg0 int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAdmin", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAdmin indicates an expected call of IsAdmin.
func (mr *MockIUserServiceMockRecorder) IsAdmin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAdmin", reflect.TypeOf((*MockIUserService)(nil).IsAdmin), arg0)
}

// VerifyUser mocks base method.
func (m *MockIUserService) VerifyUser(arg0, arg1 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyUser indicates an expected call of VerifyUser.
func (mr *MockIUserServiceMockRecorder) VerifyUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUser", reflect.TypeOf((*MockIUserService)(nil).VerifyUser), arg0, arg1)
}
