package gin_server

import (
	"fmt"
	"log"
)

func (s *Server) Log(messages ...string) {
	logString := fmt.Sprintf("%s: ", s.Signature())

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}
	log.Printf(logString)
}
