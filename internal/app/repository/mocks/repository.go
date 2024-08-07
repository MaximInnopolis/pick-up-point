// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go
//
// Generated by this command:
//
//	mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	models "route/internal/app/models"

	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AcceptOrder mocks base method.
func (m *MockRepository) AcceptOrder(order *models.Order, packagingType *models.PackagingType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptOrder", order, packagingType)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptOrder indicates an expected call of AcceptOrder.
func (mr *MockRepositoryMockRecorder) AcceptOrder(order, packagingType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptOrder", reflect.TypeOf((*MockRepository)(nil).AcceptOrder), order, packagingType)
}

// AcceptReturn mocks base method.
func (m *MockRepository) AcceptReturn(order models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptReturn", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptReturn indicates an expected call of AcceptReturn.
func (mr *MockRepositoryMockRecorder) AcceptReturn(order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptReturn", reflect.TypeOf((*MockRepository)(nil).AcceptReturn), order)
}

// GetAllOrders mocks base method.
func (m *MockRepository) GetAllOrders() ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOrders")
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOrders indicates an expected call of GetAllOrders.
func (mr *MockRepositoryMockRecorder) GetAllOrders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOrders", reflect.TypeOf((*MockRepository)(nil).GetAllOrders))
}

// GetOrderByID mocks base method.
func (m *MockRepository) GetOrderByID(orderID int) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", orderID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockRepositoryMockRecorder) GetOrderByID(orderID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockRepository)(nil).GetOrderByID), orderID)
}

// IssueOrder mocks base method.
func (m *MockRepository) IssueOrder(orderID int, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueOrder", orderID, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// IssueOrder indicates an expected call of IssueOrder.
func (mr *MockRepositoryMockRecorder) IssueOrder(orderID, hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueOrder", reflect.TypeOf((*MockRepository)(nil).IssueOrder), orderID, hash)
}

// ListOrders mocks base method.
func (m *MockRepository) ListOrders(userID, lastN int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrders", userID, lastN)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrders indicates an expected call of ListOrders.
func (mr *MockRepositoryMockRecorder) ListOrders(userID, lastN any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrders", reflect.TypeOf((*MockRepository)(nil).ListOrders), userID, lastN)
}

// ListReturns mocks base method.
func (m *MockRepository) ListReturns(page, pageSize int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListReturns", page, pageSize)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListReturns indicates an expected call of ListReturns.
func (mr *MockRepositoryMockRecorder) ListReturns(page, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListReturns", reflect.TypeOf((*MockRepository)(nil).ListReturns), page, pageSize)
}

// ReturnOrder mocks base method.
func (m *MockRepository) ReturnOrder(orderID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrder", orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnOrder indicates an expected call of ReturnOrder.
func (mr *MockRepositoryMockRecorder) ReturnOrder(orderID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrder", reflect.TypeOf((*MockRepository)(nil).ReturnOrder), orderID)
}
