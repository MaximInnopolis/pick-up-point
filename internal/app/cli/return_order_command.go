package cli

import (
	"errors"
	"flag"
	"fmt"

	"route/internal/app/module"
)

const returnOrder = "return-order"

type ReturnOrderCommand struct {
	Module module.Module
}

func (r ReturnOrderCommand) Name() string {
	return returnOrder
}

func (r ReturnOrderCommand) Description() string {
	return "Вернуть заказ курьеру:" +
		" использование return-order --orderID=SomeID \n" +
		"--orderID=SomeID: обязательный параметр, ID заказа."
}

// Call is a method to return order to courier
func (r ReturnOrderCommand) Call(args []string) error {
	var orderID int

	// Parse flags
	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.IntVar(&orderID, "orderID", 0, "use --orderID=SomeID")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if orderID == 0 {
		return errors.New("не указан обязательный параметр orderID")
	}

	err := r.Module.ReturnOrder(orderID)
	if err != nil {
		return err
	}

	fmt.Println("Заказ успешно возвращен курьеру")
	return nil
}
