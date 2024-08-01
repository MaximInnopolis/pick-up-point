package cli

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"route/internal/app/models"
	"route/internal/app/module"
)

const acceptOrder = "accept-order"

type AcceptOrderCommand struct {
	Module module.Module
}

func (a AcceptOrderCommand) Name() string {
	return acceptOrder
}

func (a AcceptOrderCommand) Description() string {
	return "Принять заказ от курьера:" +
		" использование accept-order --orderID=SomeID --userID=SomeID --deadline=SomeDate\n" +
		"--orderID=SomeID: обязательный параметр, ID заказа.\n" +
		"--userID=SomeID: обязательный параметр, ID пользователя.\n" +
		"--deadline=SomeDate: обязательный параметр, дата, до которой будет хранится заказ на ПВЗ, указывается в формате RFC3339.\n" +
		"--packagingType=SomeType: обязательный параметр, тип упаковки. Может иметь значения пакет, коробка, пакет\n" +
		"--weight=SomeWeight: обязательный параметр, вес заказа.\n" +
		"--cost=SomeCost: обязательный параметр, стоимость заказа."
}

// Call is a method to accept order from courier
func (a AcceptOrderCommand) Call(args []string) error {
	var orderID, userID int
	var weight, cost float64
	var deadline, packagingType string

	// Parse flags
	fs := flag.NewFlagSet(acceptOrder, flag.ContinueOnError)
	fs.IntVar(&orderID, "orderID", 0, "use --orderID=SomeID")
	fs.IntVar(&userID, "userID", 0, "use --userID=SomeID")
	fs.StringVar(&deadline, "deadline", "", "use --deadline=SomeDate")
	fs.StringVar(&packagingType, "packagingType", "", "use --packagingType=SomeType")
	fs.Float64Var(&weight, "weight", 0, "use --weight=SomeWeight")
	fs.Float64Var(&cost, "cost", 0, "use --cost=SomeCost")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if orderID == 0 {
		return errors.New("не указан обязательный параметр orderID")
	}
	if userID == 0 {
		return errors.New("не указан обязательный параметр userID")
	}
	if deadline == "" {
		return errors.New("не указан обязательный параметр deadline")
	}
	if packagingType == "" {
		return errors.New("не указан обязательный параметр packagingType")
	}
	if weight == 0 {
		return errors.New("не указан обязательный параметр weight")
	}
	if cost == 0 {
		return errors.New("не указан обязательный параметр cost")
	}

	parsedDeadline, err := parseTime(deadline)
	if err != nil {
		return err
	}

	order := models.NewOrder(orderID, userID, parsedDeadline, cost, weight)

	err = a.Module.AcceptOrder(order, models.ToPackageType(packagingType))
	if err != nil {
		return err
	}

	fmt.Printf("Заказ успешно принят командой %s\n", a.Name())
	return nil
}

// parseTime is a helper function to parse time
func parseTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errors.New("неверный формат даты")
	}
	return parsedTime, nil
}
