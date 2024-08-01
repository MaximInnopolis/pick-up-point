package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"route/internal/app/module"
)

const issueOrder = "issue-order"

type IssueOrderCommand struct {
	Module module.Module
}

func (i IssueOrderCommand) Name() string {
	return issueOrder
}

func (i IssueOrderCommand) Description() string {
	return "Выдать заказ пользователю:" +
		" использование issue-order --orderIDs=ID1,ID2,ID3,...\n" +
		"--orderIDs=ID1,ID2,ID3,...: обязательный параметр, ID заказов, которые необходимо выдать пользователю, разделенные запятой."
}

// Call is a method to issue order to client
func (i IssueOrderCommand) Call(args []string) error {
	var orderIDs string

	// Parse flags
	fs := flag.NewFlagSet(issueOrder, flag.ContinueOnError)
	fs.StringVar(&orderIDs, "orderIDs", "", "use --orderIDs=ID1,ID2,ID3,...")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if orderIDs == "" {
		return errors.New("не указан обязательный параметр orderIDs")
	}

	// Split the orderIDs string into individual IDs
	orderIDStrs := strings.Split(orderIDs, ",")

	// Convert each ID to an integer and issue the order
	for _, orderIDStr := range orderIDStrs {
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			return fmt.Errorf("не удалось преобразовать ID заказа в число: %v", err)
		}

		err = i.Module.IssueOrder(orderID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Заказы успешно выданы")
	return nil
}
