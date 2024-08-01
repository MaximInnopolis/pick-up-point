package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"route/internal/app/models"
	"route/internal/app/module"
)

type KafkaProducer interface {
	SendMessage(message []byte) error
}

type KafkaConsumer interface {
	ReadMessages() error
}

type Command interface {
	Name() string
	Description() string
	Call(args []string) error
}

type CLI struct {
	Module      module.Module
	Producer    KafkaProducer
	Consumer    KafkaConsumer
	commands    map[string]Command
	wg          sync.WaitGroup
	cond        *sync.Cond
	workerCount int
	event       models.Event
	outputMode  string
}

// New is a constructor for CLI
func New(commandMap map[string]Command, outputMode string, consumer KafkaConsumer, producer KafkaProducer) *CLI {
	workersCommand := commandMap["set-workers"].(*WorkersCommand)
	return &CLI{
		commands:    commandMap,
		cond:        sync.NewCond(&sync.Mutex{}),
		workerCount: workersCommand.GetWorkersCount(),
		outputMode:  outputMode,
	}
}

// NewCommands is a function to initialize all commands
func NewCommands(module module.Module) map[string]Command {
	workersCommand := WorkersCommand{}
	workersCommand = workersCommand.NewWorkersCommand()

	return map[string]Command{
		"accept-order":  AcceptOrderCommand{Module: module},
		"return-order":  ReturnOrderCommand{Module: module},
		"issue-order":   IssueOrderCommand{Module: module},
		"list-orders":   ListOrdersCommand{Module: module},
		"accept-return": AcceptReturnCommand{Module: module},
		"list-returns":  ListReturnsCommand{Module: module},
		"set-workers":   &workersCommand,
	}
}

// Run is a method to run CLI
func (c *CLI) Run() error {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine that will perform cleanup when a signal is received
	go func() {
		<-sigChan
		fmt.Println("Получен сигнал остановки. Ожидание завершения всех задач...")
		c.wg.Wait()
		fmt.Println("\"Все задачи завершены. Выход...")
		os.Exit(0)
	}()

	// Cycle for commands reading
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		commandLine, _ := reader.ReadString('\n')
		commandLine = strings.TrimSuffix(commandLine, "\n")

		// Exit from CLI
		if commandLine == "exit" {
			fmt.Println("Exiting...")
			return nil
		}

		// Split the command line into command name and arguments
		fields := strings.Fields(commandLine)
		if len(fields) == 0 {
			continue
		}
		commandName := fields[0]
		args := fields[1:]

		c.executeCommand(commandName, args)
	}
}

func (c *CLI) executeCommand(commandName string, args []string) {
	if commandName == "help" {
		mustHelp(c.commands)
		return
	}

	cmd, ok := c.commands[commandName]
	if !ok {
		fmt.Println("команда не установлена")
		return
	}

	event := models.NewEvent(commandName, args)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Ошибка при сериализации события:", err)
		return
	}

	// Check outputMode before sending message or printing to console
	if c.outputMode == "kafka" {
		if err = c.Producer.SendMessage(eventBytes); err != nil {
			fmt.Println("Ошибка при отправке события в Kafka:", err)
			return
		}
		c.startReadingMessages()
	} else if c.outputMode == "stdout" {
		fmt.Printf("Event: %s\n", string(eventBytes))
	}

	// Execute the command without starting a new consumer for each command
	c.executeCommandWithWorker(commandName, args, cmd)
}

// New method to handle command execution with worker
func (c *CLI) executeCommandWithWorker(commandName string, args []string, cmd Command) {
	// Create a channel to pass error from goroutine
	errChan := make(chan error, 1)

	// Before starting a new goroutine, check the condition
	c.cond.L.Lock()
	for c.workerCount <= 0 {
		c.cond.Wait()
	}
	c.workerCount--
	c.cond.L.Unlock()

	// Generate a unique ID for the goroutine
	goroutineID := uuid.New()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done() // Decrement the counter when the goroutine finishes
		defer func() {
			c.cond.L.Lock()
			c.workerCount++
			c.cond.L.Unlock()
			c.cond.Signal()
		}()
		fmt.Printf("Горутина с ID %s начала выполнение команды %s\n", goroutineID, commandName)
		err := cmd.Call(args)
		if err != nil {
			// Pass the error to the channel
			errChan <- err
			fmt.Printf("Горутина с ID %s не смогла выполнить команду %s из-за ошибки: %v\n", goroutineID, commandName, err)
		} else {
			fmt.Printf("Горутина с ID %s успешно завершила выполнение команды %s\n", goroutineID, commandName)
		}
		close(errChan)
	}()
}

func (c *CLI) startReadingMessages() {
	go func() {
		err := c.Consumer.ReadMessages()
		if err != nil {
			fmt.Printf("Ошибка при чтении сообщений из Kafka: %v\n", err)
		}
	}()
}
