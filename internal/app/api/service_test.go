package service

import (
	"context"
	"errors"
	"testing"

	"route/internal/app/models"
	mockmodule "route/internal/app/module/mocks"
	order "route/pkg/api/proto/order/v1/order/v1"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderService_AcceptOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		orderRequest   *order.OrderRequest
		mockSetup      func()
		expectedResult *order.OrderResponse
		expectedError  error
	}{
		{
			name: "success",
			orderRequest: &order.OrderRequest{
				OrderId:       1,
				UserId:        2,
				Weight:        3.5,
				PackagingType: "коробка",
			},
			mockSetup: func() {
				mockModule.EXPECT().
					AcceptOrder(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedResult: &order.OrderResponse{Status: "success"},
			expectedError:  nil,
		},
		{
			name: "module error",
			orderRequest: &order.OrderRequest{
				OrderId:       1,
				UserId:        2,
				Weight:        3.5,
				PackagingType: "ящик",
			},
			mockSetup: func() {
				mockModule.EXPECT().
					AcceptOrder(gomock.Any(), gomock.Any()).
					Return(errors.New("internal error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("rpc error: code = Internal desc = internal error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.mockSetup()
			result, err := orderService.AcceptOrder(context.Background(), tc.orderRequest)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestOrderService_IssueOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		orderRequest   *order.OrderRequest
		setupMock      func()
		expectedResult *order.OrderResponse
		expectedError  string
	}{
		{
			name: "successful issue",
			orderRequest: &order.OrderRequest{
				OrderId: 1,
			},
			setupMock: func() {
				mockModule.EXPECT().
					IssueOrder(gomock.Eq(1)).
					Return(nil).
					Times(1)
			},
			expectedResult: &order.OrderResponse{Status: "success"},
			expectedError:  "",
		},
		{
			name: "issue error",
			orderRequest: &order.OrderRequest{
				OrderId: 2,
			},
			setupMock: func() {
				mockModule.EXPECT().
					IssueOrder(gomock.Eq(2)).
					Return(errors.New("issue error")).
					Times(1)
			},
			expectedResult: nil,
			expectedError:  "rpc error: code = Internal desc = issue error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.setupMock()
			resp, err := orderService.IssueOrder(context.Background(), tc.orderRequest)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, resp)
			}
		})
	}
}

func TestOrderService_ReturnOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		orderRequest   *order.OrderRequest
		setupMock      func()
		expectedError  string
		expectedResult *order.OrderResponse
	}{
		{
			name: "Return order success",
			orderRequest: &order.OrderRequest{
				OrderId: 1,
			},
			setupMock: func() {
				mockModule.EXPECT().ReturnOrder(1).Return(nil)
			},
			expectedError:  "",
			expectedResult: &order.OrderResponse{Status: "success"},
		},
		{
			name: "Return order not found error",
			orderRequest: &order.OrderRequest{
				OrderId: 2,
			},
			setupMock: func() {
				mockModule.EXPECT().ReturnOrder(2).Return(errors.New("заказ с ID 2 не найден"))
			},
			expectedError:  "rpc error: code = Internal desc = заказ с ID 2 не найден",
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.setupMock()
			resp, err := orderService.ReturnOrder(context.Background(), tc.orderRequest)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, resp)
			}
		})
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		listRequest    *order.ListOrdersRequest
		setupMock      func()
		expectedError  string
		expectedResult *order.ListResponse
	}{
		{
			name: "List orders success",
			listRequest: &order.ListOrdersRequest{
				UserId: 1,
				LastN:  2,
			},
			setupMock: func() {
				mockModule.EXPECT().ListOrders(1, 2).Return([]models.Order{
					{OrderID: 1, UserID: 1, Weight: 5},
					{OrderID: 2, UserID: 1, Weight: 10},
				}, nil)
			},
			expectedError: "",
			expectedResult: &order.ListResponse{
				Orders: []*order.OrderInfo{
					{OrderId: 1, UserId: 1, Status: "Success", Weight: 5},
					{OrderId: 2, UserId: 1, Status: "Success", Weight: 10},
				},
			},
		},
		{
			name: "List orders internal error",
			listRequest: &order.ListOrdersRequest{
				UserId: 2,
				LastN:  3,
			},
			setupMock: func() {
				mockModule.EXPECT().ListOrders(2, 3).Return(nil, errors.New("internal error"))
			},
			expectedError:  "rpc error: code = Internal desc = internal error",
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.setupMock()
			resp, err := orderService.ListOrders(context.Background(), tc.listRequest)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, resp)
			}
		})
	}
}

func TestOrderService_AcceptReturn(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		setupMock      func()
		orderRequest   *order.OrderRequest
		expectedResult *order.OrderResponse
		expectedError  string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockModule.EXPECT().AcceptReturn(gomock.Any(), gomock.Any()).Return(nil)
			},
			orderRequest: &order.OrderRequest{
				OrderId: 1,
				UserId:  1,
			},
			expectedResult: &order.OrderResponse{Status: "success"},
		},
		{
			name: "Error",
			setupMock: func() {
				mockModule.EXPECT().AcceptReturn(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))
			},
			orderRequest: &order.OrderRequest{
				OrderId: 2,
				UserId:  2,
			},
			expectedError: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.setupMock()
			resp, err := orderService.AcceptReturn(context.Background(), tc.orderRequest)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, resp)
			}
		})
	}
}

func TestOrderService_ListReturns(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModule := mockmodule.NewMockModule(ctrl)
	orderService := New(mockModule)

	testCases := []struct {
		name           string
		setupMock      func()
		listRequest    *order.ListReturnsRequest
		expectedResult *order.ListResponse
		expectedError  string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockModule.EXPECT().ListReturns(gomock.Any(), gomock.Any()).Return([]models.Order{
					{OrderID: 1, UserID: 1, Weight: 5},
					{OrderID: 2, UserID: 2, Weight: 10},
				}, nil)
			},
			listRequest: &order.ListReturnsRequest{
				Page:     1,
				PageSize: 2,
			},
			expectedResult: &order.ListResponse{
				Orders: []*order.OrderInfo{
					{OrderId: 1, UserId: 1, Status: "Success", Weight: 5},
					{OrderId: 2, UserId: 2, Status: "Success", Weight: 10},
				},
			},
		},
		{
			name: "Error",
			setupMock: func() {
				mockModule.EXPECT().ListReturns(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
			},
			listRequest: &order.ListReturnsRequest{
				Page:     1,
				PageSize: 2,
			},
			expectedError: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.setupMock()
			resp, err := orderService.ListReturns(context.Background(), tc.listRequest)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, resp)
			}
		})
	}
}
