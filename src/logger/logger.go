package logger

import "log"

// Log message.
func LogMessage(msg string) {
	log.Println(msg)
}

// Log error.
func LogError(err error, msg string) {
	if err != nil {
		log.Printf("msg: %s\nerror: %s", msg, err.Error())
	}
}
