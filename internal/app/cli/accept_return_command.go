package cli

import (
	"errors"
	"flag"
	"fmt"

	"route/internal/app/module"
)

const acceptReturn = "accept-return"

type AcceptReturnCommand struct {
	Module module.Module
}

func (a AcceptReturnCommand) Name() string {
	return acceptReturn
}

func (a AcceptReturnCommand) Description() string {
	return "Принять возврат от пользователя:" +
		" использование accept-return --userID=SomeID --orderID=SomeID\n" +
		"--userID=SomeID: обязательный параметр, ID пользователя.\n" +
		"--orderID=SomeID: обязательный параметр, ID заказа."
}

// Call is a method to accept return from client
func (a AcceptReturnCommand) Call(args []string) error {
	var orderID, userID int

	// Parse flags
	fs := flag.NewFlagSet(acceptReturn, flag.ContinueOnError)
	fs.IntVar(&orderID, "orderID", 0, "use --orderID=SomeID")
	fs.IntVar(&userID, "userID", 0, "use --userID=SomeID")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if orderID == 0 {
		return errors.New("не указан обязательный параметр orderID")
	}
	if userID == 0 {
		return errors.New("не указан обязательный параметр userID")
	}

	err := a.Module.AcceptReturn(orderID, userID)
	if err != nil {
		return err
	}

	fmt.Println("Возврат успешно принят")
	return nil
}
