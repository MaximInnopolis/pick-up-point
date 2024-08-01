package cli

import (
	"flag"
	"fmt"

	"route/internal/app/module"
)

const listReturns = "list-returns"

type ListReturnsCommand struct {
	Module module.Module
}

func (l ListReturnsCommand) Name() string {
	return listReturns
}

func (l ListReturnsCommand) Description() string {
	return "Вывести список возвратов пользователя пагинированно: использование list-returns [--page=pageNumber] [--pageSize=pageSizeNumber].\n" +
		"--page=pageNumber: номер страницы (по умолчанию 1)\n" +
		"--pageSize=pageSizeNumber: количество элементов на странице (по умолчанию 5)\n" +
		"Если параметры не указаны, по умолчанию выводятся первые 5 возвратов."
}

// Call is a method to list returns
func (l ListReturnsCommand) Call(args []string) error {
	var page, pageSize int

	// Parse flags
	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.IntVar(&page, "page", 1, "use --page=pageNumber")
	fs.IntVar(&pageSize, "pageSize", 5, "use --pageSize=pageSizeNumber")
	if err := fs.Parse(args); err != nil {
		return err
	}

	list, err := l.Module.ListReturns(page, pageSize)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		fmt.Println("Список возвратов пуст")
		return nil
	}

	for _, order := range list {
		fmt.Printf("OrderID: %v\nUserID: %v\n", order.OrderID, order.UserID)
	}
	return nil
}
