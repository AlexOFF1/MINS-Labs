package observer

import (
	"fmt"
	"log"
)

type LoggerObserver struct{}

func (l *LoggerObserver) Notify(event Event, data interface{}) {
	log.Printf("[EVENT] %s: %+v\n", event, data)
}

type ConsoleObserver struct{}

func (c *ConsoleObserver) Notify(event Event, data interface{}) {
	fmt.Printf("🔔 Уведомление: событие %s\nДанные: %+v\n", event, data)
}
