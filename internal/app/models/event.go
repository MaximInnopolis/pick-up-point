package models

import (
	"strings"
	"time"
)

// Event содержит информацию о событии
type Event struct {
	Time    string `json:"time"`
	Command string `json:"command"`
	Args    string `json:"args"`
}

func NewEvent(commandName string, args []string) *Event {
	return &Event{
		Time:    time.Now().Format(time.RFC3339),
		Command: commandName,
		Args:    strings.Join(args, " "),
	}
}
