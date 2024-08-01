package module

import (
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"route/internal/app/models"
	"route/internal/app/repository"
	"route/internal/app/repository/postgresql"
	"route/pkg/hash"
)

type IMCache[K comparable, V any] interface {
	Set(key K, value V, now time.Time)
	Get(key K) (V, bool)
	Delete(key K)
}

type OrderModule struct {
	repo                repository.Repository
	cache               IMCache[int, models.Order]
	issuedOrdersCounter prometheus.Counter
}

func New(repo repository.Repository, cache IMCache[int, models.Order]) *OrderModule {
	return &OrderModule{
		repo:  repo,
		cache: cache,
	}
}

func (m OrderModule) AcceptOrder(order *models.Order, packagingType models.PackageType) error {
	// Check if the order is already in cache
	_, ok := m.cache.Get(order.OrderID)
	if ok {
		// If the order is in cache, return an error
		return fmt.Errorf("заказ с ID %d уже существует", order.OrderID)
	}

	foundOrder, err := m.repo.GetOrderByID(order.OrderID)
	if err != nil && !errors.Is(err, postgresql.ErrOrderNotFound) {
		return err
	}

	if foundOrder != nil {
		// Order with orderID already exists, return an error
		return fmt.Errorf("заказ с ID %d уже существует", order.OrderID)
	}

	// Check that deadline is not in the past
	if order.Deadline.Before(time.Now()) {
		return errors.New("срок хранения не может быть в прошлом")
	}

	// Check the packaging type and get the packaging type struct
	pt, err := checkPackagingType(packagingType, order.Weight)
	if err != nil {
		return err
	}

	// Sum total cost by adding cost to its additional cost
	totalCost := pt.AdditionalCost + order.Cost

	// Create a new order and increase cost by additional cost
	modifiedOrder := models.NewOrder(order.OrderID, order.UserID, order.Deadline, totalCost, order.Weight)

	res := m.repo.AcceptOrder(modifiedOrder, pt)
	if res == nil {
		// Set the order to cache
		m.cache.Set(order.OrderID, *modifiedOrder, time.Now())
	}
	return res
}

func (m OrderModule) ReturnOrder(orderID int) error {
	// Check if the order is already in cache
	cachedOrder, found := m.cache.Get(orderID)
	if found {
		cond := processReturnOrderCondition(&cachedOrder)
		if cond != nil {
			return cond
		}
		res := m.repo.ReturnOrder(orderID)
		// If the order is successfully returned, delete it from cache
		m.cache.Delete(orderID)
		return res
	}

	order, err := m.repo.GetOrderByID(orderID)
	if err != nil && !errors.Is(err, postgresql.ErrOrderNotFound) {
		return err
	}

	if order == nil {
		return fmt.Errorf("заказ с ID %d не найден", orderID)
	}
	cond := processReturnOrderCondition(order)
	if cond != nil {
		return cond
	}

	return m.repo.ReturnOrder(orderID)
}

func processReturnOrderCondition(order *models.Order) error {
	// If order is already returned, return an error
	if order.IssuedToUser {
		return fmt.Errorf("заказ с ID %d уже был выдан клиенту", order.OrderID)
	}

	// If order is not expired, return an error
	if order.Deadline.After(time.Now()) {
		return fmt.Errorf("заказ с ID %d еще не просрочен", order.OrderID)
	}

	return nil
}

func (m OrderModule) IssueOrder(orderID int) error {

	// Attempt to retrieve the order from the cache
	cachedOrder, found := m.cache.Get(orderID)
	if found {
		return processIssueCondition(&cachedOrder)
	}

	order, err := m.repo.GetOrderByID(orderID)
	if err != nil && !errors.Is(err, postgresql.ErrOrderNotFound) {
		return err
	}

	if order == nil {
		return fmt.Errorf("заказ с ID %d не найден", orderID)
	}

	cond := processIssueCondition(order)
	if cond != nil {
		return cond
	}

	err = m.repo.IssueOrder(orderID, hash.GenerateHash())
	if err != nil {
		return err
	}

	// Get the updated order from the database
	updatedOrder, err := m.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if updatedOrder == nil {
		return fmt.Errorf("не удалось получить обновленные данные заказа с ID %d", orderID)
	}

	// Set the updated order to cache
	m.cache.Set(orderID, *updatedOrder, time.Now())

	// Increment counter
	m.issuedOrdersCounter.Inc()

	return nil
}

func processIssueCondition(order *models.Order) error {
	// If order is already issued to user, return an error
	if order.IssuedToUser {
		return fmt.Errorf("заказ с ID %d уже был выдан клиенту", order.OrderID)
	}

	// If order is not received from courier, return an error
	if !order.ReceivedFromCourier {
		return fmt.Errorf("заказ с ID %d не был получен курьером", order.OrderID)
	}

	// If order is expired, return an error
	if order.Deadline.Before(time.Now()) {
		return fmt.Errorf("заказ с ID %d просрочен", order.OrderID)
	}

	return nil
}

func (m OrderModule) ListOrders(userID, lastN int) ([]models.Order, error) {
	return m.repo.ListOrders(userID, lastN)
}

func (m OrderModule) AcceptReturn(orderID, userID int) error {
	// Check if the order is already in cache
	cachedOrder, found := m.cache.Get(orderID)
	if found && cachedOrder.UserID == userID {
		return processAcceptReturnCondition(&cachedOrder)
	}

	order, err := m.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order == nil || order.UserID != userID {
		return fmt.Errorf("заказ с ID %d не найден", orderID)
	}

	err = processAcceptReturnCondition(order)
	if err != nil {
		return err
	}
	order.Hash = hash.GenerateHash()

	err = m.repo.AcceptReturn(*order)
	if err != nil {
		return err
	}

	// Get the updated order from the database
	updatedOrder, err := m.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if updatedOrder == nil {
		return fmt.Errorf("не удалось получить обновленные данные заказа с ID %d", orderID)
	}

	// Set the updated order to cache
	m.cache.Set(orderID, *updatedOrder, time.Now())

	return nil
}

func processAcceptReturnCondition(order *models.Order) error {
	if order.IsReturned {
		return fmt.Errorf("заказ с ID %d уже был возвращен", order.OrderID)
	}

	if !order.IssuedToUser {
		return fmt.Errorf("заказ с ID %d не был выдан клиенту", order.OrderID)
	}

	if time.Since(order.IssuedAt) > 2*24*time.Hour {
		return fmt.Errorf("заказ с ID %d не может быть возвращен, так как прошло более двух дней с момента его выдачи", order.OrderID)
	}

	return nil
}

func (m OrderModule) ListReturns(page, pageSize int) ([]models.Order, error) {
	return m.repo.ListReturns(page, pageSize)
}

func checkPackagingType(packagingType models.PackageType, weight float64) (*models.PackagingType, error) {
	// Check if the provided packaging type is allowed
	switch packagingType {
	case models.Package:
		if weight >= models.PackageWeightLimit {
			return nil, fmt.Errorf("вес заказа превышает допустимый для пакета: %f", weight)
		}
		return models.NewPackagingType(models.Package, models.PackageCost), nil

	case models.Box:
		if weight >= models.BoxWeightLimit {
			return nil, fmt.Errorf("вес заказа превышает допустимый для коробки: %f", weight)
		}
		return models.NewPackagingType(models.Box, models.BoxCost), nil

	case models.Film:
		return models.NewPackagingType(models.Film, models.FilmCost), nil

	default:
		// If the packaging type is not allowed, return an error
		return nil, fmt.Errorf("недопустимый тип упаковки: %s", packagingType)
	}
}
