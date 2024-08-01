//go:generate mockgen -source ./module.go -destination=./mocks/module.go -package=mock_module

package module

import "route/internal/app/models"

// Module is an interface for module
type Module interface {
	AcceptOrder(order *models.Order, packagingType models.PackageType) error
	ReturnOrder(orderID int) error
	IssueOrder(orderID int) error
	ListOrders(userID, lastN int) ([]models.Order, error)
	AcceptReturn(orderID, userID int) error
	ListReturns(page, pageSize int) ([]models.Order, error)
}
