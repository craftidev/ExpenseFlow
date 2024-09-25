package utils

import (
    "log"
    "errors"
    "fmt"
)

func LogError(format string, args ...interface{}) error {
    errMessage := fmt.Sprintf(format, args...)
    log.Printf("Error: %s", errMessage)
    return errors.New(errMessage)
}
