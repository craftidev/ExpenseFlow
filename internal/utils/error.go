package utils

import (
	"errors"
	"fmt"
    "log"
	"runtime"
)


func LogError(format string, args ...interface{}) error {
    errMessage := fmt.Sprintf(format, args...)

    // Log the error
    pc, fn, line, _ := runtime.Caller(1)
    log.Printf("[error] %s:%d\n\t[call] %s:\n\t[msg] %s\n", fn, line, runtime.FuncForPC(pc).Name(), errMessage)

    return errors.New(errMessage)
}
