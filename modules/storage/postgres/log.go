package postgres

import (
	"fmt"
	"log"
)

func (p *Postgres) Log(messages ...string) {
	logString := fmt.Sprintf("%s: ", p.Signature())

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}
	log.Printf(logString)
}
