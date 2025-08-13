package utils

import (
	"log"
	"os"
)

// Logger wraps the standard log.Logger
var Logger *log.Logger

func InitLogger() {
	// Create a new logger instance
	Logger = log.New(os.Stdout, "CHAT_APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info prints an informational message
func Info(msg string) {
	Logger.SetPrefix("INFO: ")
	Logger.Println(msg)
}

// Error prints an error message
func Error(msg string) {
	Logger.SetPrefix("ERROR: ")
	Logger.Println(msg)
}

// Debug prints a debug message
func Debug(msg string) {
	Logger.SetPrefix("DEBUG: ")
	Logger.Println(msg)
}
