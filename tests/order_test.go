//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route/internal/app/models"
	"route/internal/app/repository/postgresql"
	"testing"
	"time"
)

func TestAcceptOrder(t *testing.T) {

	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	order := &models.Order{
		OrderID:  1,
		UserID:   1,
		Deadline: time.Now().Add(24 * time.Hour),
		Cost:     100,
		Weight:   5,
	}
	packagingType := &models.PackagingType{
		Type:           models.ToPackageType("пакет"),
		AdditionalCost: 5,
	}

	// act
	err := repo.AcceptOrder(order, packagingType)

	// assert
	require.NoError(t, err, "AcceptOrder should not error")

	var orderID int
	err = db.DB.GetQueryEngine(context.Background()).QueryRow(context.Background(),
		"SELECT id FROM orders WHERE id = $1", order.OrderID).Scan(&orderID)

	require.NoError(t, err, "Querying inserted order should not error")
	assert.Equal(t, order.OrderID, orderID, "Expected order ID to match")

	var packagingTypeID int

	err = db.DB.GetQueryEngine(context.Background()).QueryRow(context.Background(),
		"SELECT id FROM packaging_types WHERE type = $1", string(packagingType.Type)).Scan(&packagingTypeID)
	require.NoError(t, err, "Querying inserted packaging type should not error")
}

func TestReturnOrder(t *testing.T) {
	// arrange

	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	order := &models.Order{
		OrderID:  2,
		UserID:   1,
		Deadline: time.Now().Add(24 * time.Hour),
		Cost:     100,
		Weight:   5,
	}
	// Insert test order
	_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
		"INSERT INTO orders (id, user_id, deadline, cost, weight) VALUES ($1, $2, $3, $4, $5)",
		order.OrderID, order.UserID, order.Deadline, order.Cost, order.Weight)
	require.NoError(t, err, "Inserting test order should not error")

	// Act
	err = repo.ReturnOrder(order.OrderID)
	require.NoError(t, err, "ReturnOrder should not error")

	// Assert
	var exists bool
	err = db.DB.GetQueryEngine(context.Background()).QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)", order.OrderID).Scan(&exists)

	require.NoError(t, err, "Querying for order existence should not error")
	assert.False(t, exists, "Order should not exist after being returned")
}

func TestIssueOrder(t *testing.T) {
	// Setup Test Environment
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	order := &models.Order{
		OrderID:  3,
		UserID:   1,
		Deadline: time.Now().Add(24 * time.Hour),
		Cost:     100,
		Weight:   5,
	}
	// Insert test order
	_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
		"INSERT INTO orders (id, user_id, deadline, cost, weight) VALUES ($1, $2, $3, $4, $5)",
		order.OrderID, order.UserID, order.Deadline, order.Cost, order.Weight)

	require.NoError(t, err, "Inserting test order should not error")

	newHash := "newHashValue"

	// Act
	err = repo.IssueOrder(order.OrderID, newHash)
	require.NoError(t, err, "IssueOrder should not error")

	// Assert
	var issuedToUser bool
	var hash string
	err = db.DB.GetQueryEngine(context.Background()).QueryRow(context.Background(),
		"SELECT issued_to_user, hash FROM orders WHERE id = $1", order.OrderID).Scan(&issuedToUser, &hash)

	require.NoError(t, err, "Querying updated order should not error")
	assert.True(t, issuedToUser, "Order should be marked as issued to user")
	assert.Equal(t, newHash, hash, "Hash should match the new hash value")
}

func TestListOrders(t *testing.T) {
	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	userID := 1
	now := time.Now()
	ordersToInsert := []models.Order{
		{OrderID: 4, UserID: userID, Deadline: now.Add(24 * time.Hour), Cost: 100, Weight: 5},
		{OrderID: 5, UserID: userID, Deadline: now.Add(48 * time.Hour), Cost: 200, Weight: 10},
	}
	for _, order := range ordersToInsert {
		_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
			"INSERT INTO orders (id, user_id, deadline, cost, weight) VALUES ($1, $2, $3, $4, $5)",
			order.OrderID, order.UserID, order.Deadline, order.Cost, order.Weight)
		require.NoError(t, err, "Inserting test order should not error")
	}

	// Act
	retrievedOrders, err := repo.ListOrders(userID, 2)
	require.NoError(t, err, "ListOrders should not error")

	// Assert
	require.Len(t, retrievedOrders, len(ordersToInsert), "The number of retrieved orders should match the number of inserted orders")
	for i, order := range retrievedOrders {
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].OrderID, order.OrderID, "OrderID should match")
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].Cost, order.Cost, "Cost should match")
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].Weight, order.Weight, "Weight should match")
	}
}

func TestAcceptReturn(t *testing.T) {
	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	order := &models.Order{
		OrderID:  6,
		UserID:   1,
		Deadline: time.Now().Add(24 * time.Hour),
		Cost:     100,
		Weight:   5,
		Hash:     "testHash",
	}
	// Insert test order
	_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
		"INSERT INTO orders (id, user_id, deadline, cost, weight, hash) VALUES ($1, $2, $3, $4, $5, $6)",
		order.OrderID, order.UserID, order.Deadline, order.Cost, order.Weight, order.Hash)
	require.NoError(t, err, "Inserting test order should not error")

	// Act
	err = repo.AcceptReturn(*order)
	require.NoError(t, err, "AcceptReturn should not error")

	// Assert
	var isReturned bool
	err = db.DB.GetQueryEngine(context.Background()).QueryRow(context.Background(),
		"SELECT is_returned FROM orders WHERE id = $1", order.OrderID).Scan(&isReturned)

	require.NoError(t, err, "Querying updated order should not error")
	assert.True(t, isReturned, "Order should be marked as returned")
}

func TestListReturns(t *testing.T) {
	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	now := time.Now()
	// Insert returned orders
	returnedOrders := []models.Order{
		{OrderID: 7, UserID: 1, Deadline: now, IsReturned: true, Cost: 100, Weight: 5},
		{OrderID: 8, UserID: 1, Deadline: now.Add(-24 * time.Hour), IsReturned: true, Cost: 200, Weight: 10},
	}
	// Insert a non-returned order
	_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
		"INSERT INTO orders (id, user_id, deadline, is_returned, cost, weight) VALUES ($1, $2, $3, $4, $5, $6)",
		3, 1, now.Add(-48*time.Hour), false, 300, 15)
	require.NoError(t, err, "Inserting non-returned order should not error")

	for _, order := range returnedOrders {
		_, err = db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
			"INSERT INTO orders (id, user_id, deadline, is_returned, cost, weight) VALUES ($1, $2, $3, $4, $5, $6)",
			order.OrderID, order.UserID, order.Deadline, order.IsReturned, order.Cost, order.Weight)
		require.NoError(t, err, "Inserting returned order should not error")
	}

	// Act
	retrievedOrders, err := repo.ListReturns(1, 10)
	require.NoError(t, err, "ListReturns should not error")

	// Assert
	require.Len(t, retrievedOrders, len(returnedOrders), "The number of retrieved orders should match the number of inserted returned orders")
	for i, order := range retrievedOrders {
		assert.Equal(t, returnedOrders[i].OrderID, order.OrderID, "OrderID should match")
		assert.True(t, order.IsReturned, "Order should be marked as returned")
	}
}

func TestGetAllOrders(t *testing.T) {
	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	now := time.Now()
	ordersToInsert := []models.Order{
		{OrderID: 9, UserID: 1, Deadline: now.Add(24 * time.Hour), Cost: 100, Weight: 5},
		{OrderID: 10, UserID: 1, Deadline: now.Add(48 * time.Hour), Cost: 200, Weight: 10},
	}
	for _, order := range ordersToInsert {
		_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
			"INSERT INTO orders (id, user_id, deadline, cost, weight) VALUES ($1, $2, $3, $4, $5)",
			order.OrderID, order.UserID, order.Deadline, order.Cost, order.Weight)
		require.NoError(t, err, "Inserting test order should not error")
	}

	// Act
	retrievedOrders, err := repo.GetAllOrders()
	require.NoError(t, err, "GetAllOrders should not error")

	// Assert
	require.Len(t, retrievedOrders, len(ordersToInsert), "The number of retrieved orders should match the number of inserted orders")
	// Assuming the orders are returned in descending order by OrderID
	for i, order := range retrievedOrders {
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].OrderID, order.OrderID, "OrderID should match")
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].Cost, order.Cost, "Cost should match")
		assert.Equal(t, ordersToInsert[len(ordersToInsert)-1-i].Weight, order.Weight, "Weight should match")
	}
}

func TestGetOrderByID(t *testing.T) {
	// arrange
	db.SetUp(t)
	defer db.TearDown(t)

	repo := postgresql.New(db.DB)

	// Prepare Test Data
	testOrder := &models.Order{
		OrderID:             11,
		UserID:              1,
		Deadline:            time.Now().Add(24 * time.Hour),
		Cost:                100,
		Weight:              5,
		IssuedToUser:        false,
		IsAtPickupPoint:     false,
		ReceivedFromCourier: true,
		IsReturned:          false,
		IssuedAt:            time.Time{},
		Hash:                "testHash",
	}
	// Insert test order
	_, err := db.DB.GetQueryEngine(context.Background()).Exec(context.Background(),
		"INSERT INTO orders (id, user_id, deadline, cost, weight, issued_to_user, is_at_pickup_point, received_from_courier, is_returned, issued_at, hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		testOrder.OrderID, testOrder.UserID, testOrder.Deadline, testOrder.Cost, testOrder.Weight, testOrder.IssuedToUser, testOrder.IsAtPickupPoint, testOrder.ReceivedFromCourier, testOrder.IsReturned, testOrder.IssuedAt, testOrder.Hash)
	require.NoError(t, err, "Inserting test order should not error")

	// Act
	retrievedOrder, err := repo.GetOrderByID(testOrder.OrderID)
	require.NoError(t, err, "GetOrderByID should not error")

	// Assert
	assert.Equal(t, testOrder.OrderID, retrievedOrder.OrderID, "OrderID should match")
	assert.Equal(t, testOrder.UserID, retrievedOrder.UserID, "UserID should match")
	assert.Equal(t, testOrder.Cost, retrievedOrder.Cost, "Cost should match")
	assert.Equal(t, testOrder.Weight, retrievedOrder.Weight, "Weight should match")
	assert.Equal(t, testOrder.IssuedToUser, retrievedOrder.IssuedToUser, "IssuedToUser should match")
	assert.Equal(t, testOrder.IsAtPickupPoint, retrievedOrder.IsAtPickupPoint, "IsAtPickupPoint should match")
	assert.Equal(t, testOrder.ReceivedFromCourier, retrievedOrder.ReceivedFromCourier, "ReceivedFromCourier should match")
	assert.Equal(t, testOrder.IsReturned, retrievedOrder.IsReturned, "IsReturned should match")
	assert.Equal(t, testOrder.Hash, retrievedOrder.Hash, "Hash should match")
}
