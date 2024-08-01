package module

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"route/internal/app/models"
	mockrepository "route/internal/app/repository/mocks"
)

type OrderMatcher struct {
	expected *models.Order
}

func (m *OrderMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*models.Order)
	if !ok {
		return false
	}

	// Сравнение всех полей или нужных полей структуры, кроме даты
	return actual.OrderID == m.expected.OrderID &&
		actual.UserID == m.expected.UserID &&
		actual.Weight == m.expected.Weight &&
		actual.Cost == m.expected.Cost &&
		actual.IssuedToUser == m.expected.IssuedToUser &&
		actual.ReceivedFromCourier == m.expected.ReceivedFromCourier
}

func (m *OrderMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func EqOrder(expected *models.Order) gomock.Matcher {
	return &OrderMatcher{expected: expected}
}

func TestCheckPackagingType(t *testing.T) {
	t.Parallel()

	// arrange
	tests := []struct {
		name          string
		packagingType models.PackageType
		weight        float64
		expectedError string
		expectedType  *models.PackagingType
	}{
		{
			name:          "invalid packaging type",
			packagingType: "invalid",
			weight:        5,
			expectedError: "недопустимый тип упаковки: invalid",
		},
		{
			name:          "valid package type with acceptable weight",
			packagingType: models.Package,
			weight:        5,
			expectedType:  models.NewPackagingType(models.Package, models.PackageCost),
		},
		{
			name:          "valid package type with excessive weight",
			packagingType: models.Package,
			weight:        15,
			expectedError: fmt.Sprintf("вес заказа превышает допустимый для пакета: %f", 15.0),
		},
		{
			name:          "valid box type with acceptable weight",
			packagingType: models.Box,
			weight:        20, // Within the box weight limit
			expectedType:  models.NewPackagingType(models.Box, models.BoxCost),
		},
		{
			name:          "valid box type with excessive weight",
			packagingType: models.Box,
			weight:        40, // Exceeds the box weight limit
			expectedError: fmt.Sprintf("вес заказа превышает допустимый для коробки: %f", 40.0),
		},
		{
			name:          "valid film type with any weight",
			packagingType: models.Film,
			weight:        100, // Arbitrary high weight, should still pass for film
			expectedType:  models.NewPackagingType(models.Film, models.FilmCost),
		},
		{
			name:          "excessively high weight for any packaging",
			packagingType: models.Package,
			weight:        1000, // Excessively high weight, should fail for package
			expectedError: fmt.Sprintf("вес заказа превышает допустимый для пакета: %f", 1000.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			pt, err := checkPackagingType(tt.packagingType, tt.weight)

			// assert
			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if pt.Type != tt.expectedType.Type || pt.AdditionalCost != tt.expectedType.AdditionalCost {
					t.Errorf("Expected packaging type %v, got %v", tt.expectedType, pt)
				}
			}
		})
	}
}

func TestModule_AcceptOrder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepository.NewMockRepository(ctrl)
	module := New(mockRepo)

	order := &models.Order{
		OrderID:  1,
		UserID:   1,
		Deadline: time.Now().Add(24 * time.Hour), // Future deadline
		Weight:   5,
		Cost:     100,
	}

	expectedOrder := &models.Order{
		OrderID: order.OrderID,
		UserID:  order.UserID,
		Weight:  order.Weight,
		Cost:    order.Cost,
	}

	t.Run("order already exists", func(t *testing.T) {
		t.Parallel()
		mockRepo.EXPECT().GetOrderByID(order.OrderID).Return(&models.Order{}, nil)

		err := module.AcceptOrder(order, models.Package)
		assert.EqualError(t, err, "заказ с ID 1 уже существует")
	})

	t.Run("deadline in the past", func(t *testing.T) {
		t.Parallel()
		pastOrder := *order
		pastOrder.Deadline = time.Now().Add(-24 * time.Hour) // Past deadline

		mockRepo.EXPECT().GetOrderByID(pastOrder.OrderID).Return(nil, pgx.ErrNoRows)

		err := module.AcceptOrder(&pastOrder, models.Package)
		assert.EqualError(t, err, "срок хранения не может быть в прошлом")
	})

	t.Run("invalid packaging type", func(t *testing.T) {
		t.Parallel()
		mockRepo.EXPECT().GetOrderByID(order.OrderID).Return(nil, pgx.ErrNoRows)

		err := module.AcceptOrder(order, "invalid")
		assert.EqualError(t, err, "недопустимый тип упаковки: invalid")
	})

	t.Run("successful order acceptance", func(t *testing.T) {
		t.Parallel()
		mockRepo.EXPECT().GetOrderByID(order.OrderID).Return(nil, pgx.ErrNoRows)
		mockRepo.EXPECT().AcceptOrder(EqOrder(expectedOrder), gomock.Any()).Return(nil)

		err := module.AcceptOrder(order, models.Package)
		assert.NoError(t, err)
	})

	t.Run("repository error on GetOrderByID", func(t *testing.T) {
		t.Parallel()
		mockRepo.EXPECT().GetOrderByID(order.OrderID).Return(nil, errors.New("database error"))

		err := module.AcceptOrder(order, models.Package)
		assert.EqualError(t, err, "database error")
	})

	t.Run("repository error on AcceptOrder", func(t *testing.T) {
		t.Parallel()
		mockRepo.EXPECT().GetOrderByID(order.OrderID).Return(nil, pgx.ErrNoRows)
		mockRepo.EXPECT().AcceptOrder(EqOrder(expectedOrder), gomock.Any()).Return("database error")

		err := module.AcceptOrder(order, models.Package)
		assert.EqualError(t, err, "database error")
	})
}

func TestModule_ReturnOrder(t *testing.T) {
	t.Parallel()

	// arrange
	ctrl := gomock.NewController(t)
	mockRepo := mockrepository.NewMockRepository(ctrl)
	mod := New(mockRepo)

	currentTime := time.Now()
	pastTime := currentTime.Add(-24 * time.Hour)
	futureTime := currentTime.Add(24 * time.Hour)

	tests := []struct {
		name          string
		orderID       int
		setupMocks    func()
		expectedError string
	}{
		{
			name:    "order not found",
			orderID: 1,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(1).Return(nil, pgx.ErrNoRows)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не найден", 1),
		},
		{
			name:    "order already issued to user",
			orderID: 2,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(2).Return(&models.Order{OrderID: 2, IssuedToUser: true}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d уже был выдан клиенту", 2),
		},
		{
			name:    "order not expired",
			orderID: 3,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(3).Return(&models.Order{OrderID: 3, Deadline: futureTime}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d еще не просрочен", 3),
		},
		{
			name:    "successful order return",
			orderID: 4,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(4).Return(&models.Order{OrderID: 4, IssuedToUser: false, Deadline: pastTime}, nil)
				mockRepo.EXPECT().ReturnOrder(4).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMocks()

			// act
			err := mod.ReturnOrder(tt.orderID)

			// assert
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			}
		})
	}
}

func TestModule_IssueOrder(t *testing.T) {
	t.Parallel()

	// arrange
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockRepository(ctrl)
	mod := New(mockRepo)

	currentTime := time.Now()
	futureTime := currentTime.Add(24 * time.Hour)

	tests := []struct {
		name          string
		orderID       int
		setupMocks    func()
		expectedError string
	}{
		{
			name:    "order not found",
			orderID: 1,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(1).Return(nil, pgx.ErrNoRows)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не найден", 1),
		},
		{
			name:    "order already issued to user",
			orderID: 2,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(2).Return(&models.Order{OrderID: 2, IssuedToUser: true}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d уже был выдан клиенту", 2),
		},
		{
			name:    "order not received from courier",
			orderID: 3,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(3).Return(&models.Order{OrderID: 3, ReceivedFromCourier: false}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не был получен курьером", 3),
		},
		{
			name:    "successful order issue",
			orderID: 4,
			setupMocks: func() {
				mockRepo.EXPECT().GetOrderByID(4).Return(&models.Order{OrderID: 4, IssuedToUser: false, ReceivedFromCourier: true, Deadline: futureTime}, nil)
				mockRepo.EXPECT().IssueOrder(4, gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMocks()

			// act
			err := mod.IssueOrder(tt.orderID)

			// assert
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			}
		})
	}
}

func TestModule_ListOrders(t *testing.T) {
	t.Parallel()

	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()

		// arrange
		var (
			userID = 1
			lastN  = 5
		)

		ctrl := gomock.NewController(t)

		mockRepo := mockrepository.NewMockRepository(ctrl)
		mod := New(mockRepo)

		mockRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return([]models.Order{{OrderID: 1}, {OrderID: 2}}, nil)

		//act
		orders, err := mod.ListOrders(userID, lastN)

		// assert
		require.NoError(t, err)
		assert.Len(t, orders, 2)
	})

	t.Run("error listing orders", func(t *testing.T) {
		t.Parallel()

		// arrange
		var (
			userID = 1
			lastN  = 5
		)

		ctrl := gomock.NewController(t)

		mockRepo := mockrepository.NewMockRepository(ctrl)
		mod := New(mockRepo)

		mockRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database error"))

		//act
		_, err := mod.ListOrders(userID, lastN)

		// assert
		require.EqualError(t, err, "database error")
	})
}

func TestModule_ListReturns(t *testing.T) {
	t.Parallel()

	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()

		// arrange
		var (
			page     = 1
			pageSize = 5
		)

		ctrl := gomock.NewController(t)

		mockRepo := mockrepository.NewMockRepository(ctrl)
		mod := New(mockRepo)

		mockRepo.EXPECT().ListReturns(gomock.Any(), gomock.Any()).Return([]models.Order{{OrderID: 1}, {OrderID: 2}}, nil)

		//act
		orders, err := mod.ListReturns(page, pageSize)

		// assert
		require.NoError(t, err)
		assert.Len(t, orders, 2)
	})

	t.Run("error listing orders", func(t *testing.T) {
		t.Parallel()

		// arrange
		var (
			page     = 1
			pageSize = 5
		)

		ctrl := gomock.NewController(t)

		mockRepo := mockrepository.NewMockRepository(ctrl)
		mod := New(mockRepo)

		mockRepo.EXPECT().ListReturns(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database error"))

		//act
		_, err := mod.ListReturns(page, pageSize)

		// assert
		require.EqualError(t, err, "database error")
	})
}

func TestModule_AcceptReturn(t *testing.T) {
	t.Parallel()

	// arrange
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockRepository(ctrl)
	mod := New(mockRepo)

	currentTime := time.Now()

	tests := []struct {
		name          string
		orderID       int
		userID        int
		setupMocks    func()
		expectedError string
	}{
		{
			name:    "order not found",
			orderID: 1,
			userID:  1,
			setupMocks: func() {
				mockRepo.EXPECT().GetAllOrders().Return([]models.Order{}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не найден", 1),
		},
		{
			name:    "order already returned",
			orderID: 2,
			userID:  1,
			setupMocks: func() {
				mockRepo.EXPECT().GetAllOrders().Return([]models.Order{
					{OrderID: 2, UserID: 1, IsReturned: true}}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d уже был возвращен", 2),
		},
		{
			name:    "order not issued to user",
			orderID: 3,
			userID:  1,
			setupMocks: func() {
				mockRepo.EXPECT().GetAllOrders().Return([]models.Order{
					{OrderID: 3, UserID: 1, IssuedToUser: false}}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не был выдан клиенту", 3),
		},
		{
			name:    "order issued more than two days ago",
			orderID: 4,
			userID:  1,
			setupMocks: func() {
				mockRepo.EXPECT().GetAllOrders().Return([]models.Order{
					{OrderID: 4, UserID: 1, IssuedToUser: true, IssuedAt: currentTime.Add(-49 * time.Hour)}}, nil)
			},
			expectedError: fmt.Sprintf("заказ с ID %d не может быть возвращен, так как прошло более двух дней с момента его выдачи", 4),
		},
		{
			name:    "successful order return",
			orderID: 5,
			userID:  1,
			setupMocks: func() {
				mockRepo.EXPECT().GetAllOrders().Return([]models.Order{
					{OrderID: 5, UserID: 1, IssuedToUser: true, IssuedAt: currentTime}}, nil)
				mockRepo.EXPECT().AcceptReturn(gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setupMocks()

			// act
			err := mod.AcceptReturn(tt.orderID, tt.userID)

			// assert
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			}
		})
	}
}
