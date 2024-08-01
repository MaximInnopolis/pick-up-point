package cli

import (
	"errors"
	"flag"
	"fmt"

	"route/internal/app/module"
)

const listOrders = "list-orders"

type ListOrdersCommand struct {
	Module module.Module
}

func (l ListOrdersCommand) Name() string {
	return listOrders
}

func (l ListOrdersCommand) Description() string {
	return "Вывести список заказов пользователя: использование list-orders --userID=ID [--lastN=Number].\n" +
		"--userID=ID: обязательный параметр, ID пользователя.\n" +
		"--lastN=SomeNumber: опциональный параметр, получить последние N заказов пользователя.\n" +
		"Если параметр SomeNumber не указан, по умолчанию выводятся последние 5 заказов."
}

// Call is a method to list orders
func (l ListOrdersCommand) Call(args []string) error {
	var lastN, userID int

	// Parse flags
	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.IntVar(&userID, "userID", 0, "use --userID=SomeID")
	fs.IntVar(&lastN, "lastN", 5, "use --lastN=SomeNumber")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if userID == 0 {
		return errors.New("не указан обязательный параметр userID")
	}

	list, err := l.Module.ListOrders(userID, lastN)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		fmt.Println("По ID пользователя заказов в ПВЗ не найдено")
		return nil
	}

	for _, order := range list {
		fmt.Printf("OrderID: %v\nUserID: %v\n\n", order.OrderID, order.UserID)
	}
	return nil
}
