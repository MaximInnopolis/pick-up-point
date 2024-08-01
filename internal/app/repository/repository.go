//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import (
	"route/internal/app/models"
)

type Repository interface {
	AcceptOrder(order *models.Order, packagingType *models.PackagingType) error
	ReturnOrder(orderID int) error
	IssueOrder(orderID int, hash string) error
	ListOrders(userID, lastN int) ([]models.Order, error)
	AcceptReturn(order models.Order) error
	ListReturns(page, pageSize int) ([]models.Order, error)

	GetAllOrders() ([]models.Order, error)
	GetOrderByID(orderID int) (*models.Order, error)
}
