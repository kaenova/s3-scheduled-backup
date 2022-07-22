package pkg

import (
	"log"
	"os"
)

type CustomLogger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

type CustomLoggerI interface {
	Error(v ...any)
	Warning(v ...any)
	Info(v ...any)
}

func NewLogger() CustomLoggerI {
	file, err := os.OpenFile("scheduler_logger.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Cannot init logger")
		os.Exit(1)
	}
	return &CustomLogger{
		InfoLogger:    log.New(file, "[INFO] ", log.Ldate|log.Ltime),
		WarningLogger: log.New(file, "[WARNING] ", log.Ldate|log.Ltime),
		ErrorLogger:   log.New(file, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

func (c *CustomLogger) Error(v ...any) {
	c.ErrorLogger.Println(v...)
}

func (c *CustomLogger) Warning(v ...any) {
	c.WarningLogger.Println(v...)
}

func (c *CustomLogger) Info(v ...any) {
	c.InfoLogger.Println(v...)
}
