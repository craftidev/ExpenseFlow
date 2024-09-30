package utils

import (
    "unicode/utf8"
)


func ValidateStringInput(input string) error {
    if !utf8.ValidString(input) {
        return LogError("invalid UTF-8 string input")
    }

    if stringContainsNullByte(input) {
        return LogError("string input contains null byte")
    }

    return nil
}

func stringContainsNullByte(input string) bool {
    for _, r := range input {
        if r == '\x00' {
            return true
        }
    }
    return false
}
