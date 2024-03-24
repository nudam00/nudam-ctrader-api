package logger

import "log"

// Logs message.
func LogMessage(msg string) {
	log.Println(msg)
}

// Logs error.
func LogError(err error, msg string) {
	if err != nil {
		log.Printf("msg: %s\nerror: %s", msg, err.Error())
	}
}
