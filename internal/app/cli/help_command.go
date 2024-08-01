package cli

import (
	"fmt"
)

func mustHelp(commands map[string]Command) {
	fmt.Println("Список команд:")
	for _, cmd := range commands {
		fmt.Printf("\n%s - %s\n", cmd.Name(), cmd.Description())
	}
}
