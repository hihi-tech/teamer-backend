package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const logDirectory = "logs/"

func NewLogger(slug string, name string) *log.Logger {
	filename := path.Join(logDirectory, fmt.Sprintf("%s.log", slug))
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	stream := io.MultiWriter(os.Stdout, logFile)

	return log.New(stream, name, log.Ldate|log.Ltime|log.Lshortfile)
}
