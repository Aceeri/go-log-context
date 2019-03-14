package logContext

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Logger struct {
	context []string
	debug   bool
}

func NewLogger(context string) Logger {
	return Logger{
		context: []string{context},
	}
}

func (logger Logger) Fork(context string) Logger {
	return Logger{
		context: append(logger.context, context),
	}
}

func (logger *Logger) SetDebug(debug bool) {
	logger.debug = debug
}

func (logger Logger) GetDebug() bool {
	return logger.debug
}

//var redundantWhitespace = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func (logger Logger) Dlog(format string, args ...interface{}) {
	if logger.debug {
		logger.Log(format, args...)
	}
}

func (logger Logger) Log(format string, args ...interface{}) {
	//formatted := strings.TrimSpace(fmt.Sprintf(format, args...))
	//formatted = redundantWhitespace.ReplaceAllString(formatted, " ")
	log.Printf("%s%s\n", logger.Context(), formatted)
}

func (logger Logger) RawLog(format string, args ...interface{}) {
	log.Printf("%s%s\n", logger.Context(), fmt.Sprintf(format, args...))
}

func (logger Logger) Elog(format string, args ...interface{}) {
	logger.Fork("error").Log(format, args...)
}

func (logger Logger) Context() string {
	buffer := bytes.NewBufferString("")
	for _, context := range logger.context {
		buffer.WriteString(fmt.Sprintf("%s: ", context))
	}
	return buffer.String()
}
