package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"route/internal/app/models"
	"route/internal/app/repository/database"
)

var ErrOrderNotFound = errors.New("order not found")

type Repo struct {
	tm database.TransactionManager
}

func New(tm database.TransactionManager) *Repo {
	return &Repo{tm: tm}
}

// AcceptOrder adds a new order to the database
func (r *Repo) AcceptOrder(order *models.Order, packagingType *models.PackagingType) error {
	return r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)

		// Insert the packaging type and get its ID
		var packagingTypeID int
		err := qe.QueryRow(ctx, "INSERT INTO packaging_types (type) VALUES ($1) RETURNING id",
			string(packagingType.Type)).Scan(&packagingTypeID)
		if err != nil {
			return err
		}

		_, err = qe.Exec(ctx,
			"INSERT INTO orders (id, user_id, deadline, issued_at, hash, packaging_type_id, cost, weight) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			order.OrderID, order.UserID, order.Deadline, order.IssuedAt, order.Hash, packagingTypeID, order.Cost, order.Weight)
		if err != nil {
			return err
		}
		return nil
	})
}

// ReturnOrder removes an order from the database.
func (r *Repo) ReturnOrder(orderID int) error {
	return r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)
		_, err := qe.Exec(ctx, "DELETE FROM orders WHERE id = $1", orderID)
		if err != nil {
			return err
		}
		return nil
	})
}

// IssueOrder updates an order in the database, marking it issued
func (r *Repo) IssueOrder(orderID int, hash string) error {
	return r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)
		_, err := qe.Exec(ctx, "UPDATE orders SET issued_to_user = true, issued_at = NOW(), hash = $1 WHERE id = $2", hash, orderID)
		if err != nil {
			return err
		}
		return nil
	})
}

// ListOrders returns a list of the user's most recent orders from the database
func (r *Repo) ListOrders(userID, lastN int) ([]models.Order, error) {
	var orders []models.Order
	ctx := context.Background()

	qe := r.tm.GetQueryEngine(ctx)
	rows, err := qe.Query(ctx,
		"SELECT id, user_id, deadline, is_returned, is_at_pickup_point, issued_to_user, issued_at, received_from_courier, hash, cost, weight FROM orders WHERE user_id = $1 ORDER BY id DESC LIMIT $2",
		userID, lastN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderID, &order.UserID, &order.Deadline, &order.IsReturned, &order.IsAtPickupPoint,
			&order.IssuedToUser, &order.IssuedAt, &order.ReceivedFromCourier, &order.Hash, &order.Cost, &order.Weight)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// AcceptReturn updates an order in the database, marking it returned
func (r *Repo) AcceptReturn(order models.Order) error {
	return r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)

		// Prepared statement for better performance
		sql := "UPDATE orders SET is_returned = true, hash = $1 WHERE id = $2"
		_, err := qe.Exec(ctx, sql, order.Hash, order.OrderID)
		if err != nil {
			return err
		}
		return nil
	})
}

// ListReturns returns a list of returned orders from the database
func (r *Repo) ListReturns(page, pageSize int) ([]models.Order, error) {
	var orders []models.Order
	err := r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)
		rows, err := qe.Query(ctx,
			"SELECT id, user_id, deadline, is_returned, is_at_pickup_point, issued_to_user, issued_at, received_from_courier, hash, cost, weight FROM orders WHERE is_returned = true ORDER BY id DESC LIMIT $1 OFFSET $2",
			pageSize, (page-1)*pageSize)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var order models.Order
			err = rows.Scan(&order.OrderID, &order.UserID, &order.Deadline, &order.IsReturned, &order.IsAtPickupPoint, &order.IssuedToUser,
				&order.IssuedAt, &order.ReceivedFromCourier, &order.Hash, &order.Cost, &order.Weight)
			if err != nil {
				return err
			}
			orders = append(orders, order)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetAllOrders returns a list of all orders from the database
func (r *Repo) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	ctx := context.Background()
	qe := r.tm.GetQueryEngine(ctx)
	rows, err := qe.Query(ctx,
		"SELECT id, user_id, deadline, is_returned, is_at_pickup_point, issued_to_user, issued_at, received_from_courier, hash, cost, weight FROM orders ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderID, &order.UserID, &order.Deadline, &order.IsReturned, &order.IsAtPickupPoint, &order.IssuedToUser,
			&order.IssuedAt, &order.ReceivedFromCourier, &order.Hash, &order.Cost, &order.Weight)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID returns the order with the given ID from the database
func (r *Repo) GetOrderByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.tm.RunRepeatableRead(context.Background(), func(ctx context.Context) error {
		qe := r.tm.GetQueryEngine(ctx)
		row := qe.QueryRow(ctx,
			"SELECT id, user_id, deadline, is_returned, is_at_pickup_point, issued_to_user, issued_at, received_from_courier, hash, cost, weight FROM orders WHERE id = $1",
			orderID)
		err := row.Scan(&order.OrderID, &order.UserID, &order.Deadline, &order.IsReturned, &order.IsAtPickupPoint, &order.IssuedToUser,
			&order.IssuedAt, &order.ReceivedFromCourier, &order.Hash, &order.Cost, &order.Weight)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrOrderNotFound
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}
