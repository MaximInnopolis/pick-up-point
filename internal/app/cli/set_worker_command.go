package cli

import (
	"errors"
	"flag"
	"fmt"
	"sync/atomic"
)

type WorkersCommand struct {
	workerCount *int32
}

func (w *WorkersCommand) NewWorkersCommand() WorkersCommand {
	count := int32(2)
	return WorkersCommand{workerCount: &count}
}

func (w WorkersCommand) Name() string {
	return "set-workers"
}

func (w WorkersCommand) Description() string {
	return "Установить количество горутин: использование set-workers [--count=SomeNumber]\n" +
		"--count=SomeNumber: необязательный параметр, количество горутин.\n" +
		"Если параметр SomeNumber не указан, по умолчанию количество горутин равно 2."
}

// Call is a method to set the number of workers
func (w *WorkersCommand) Call(args []string) error {
	var count int

	// Parse flags
	fs := flag.NewFlagSet("set-workers", flag.ContinueOnError)
	fs.IntVar(&count, "count", 2, "use --count=SomeNumber")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if count <= 0 {
		return errors.New("количество горутин должно быть больше нуля")
	}

	atomic.StoreInt32(w.workerCount, int32(count))

	fmt.Printf("Количество горутин установлено на %d\n", *w.workerCount)
	return nil
}

// GetWorkersCount is a method to get the number of workers
func (w *WorkersCommand) GetWorkersCount() int {
	return int(atomic.LoadInt32(w.workerCount))
}
