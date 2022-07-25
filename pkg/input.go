package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

// A validator that takes a string input and a valid input if returns true.
type InputValidator func(string) bool

// Return a valid string contrained by InputValidator functions
func InputString(text string, validator ...InputValidator) string {
	inputValidator := func(s string) bool { return true }
	if len(validator) > 0 {
		inputValidator = validator[0]
	}

	var input string
	var done = false
	for !done {
		fmt.Println(text)
		fmt.Scanln(&input)
		trimmedInput := strings.TrimSpace(input)
		if len(trimmedInput) != 0 && inputValidator(trimmedInput) {
			break
		}
	}
	return strings.TrimSpace(input)
}

// Return a valid string contrained by InputValidator functions
func InputInt(text string, validator ...InputValidator) int {
	inputValidator := func(s string) bool { return true }
	if len(validator) > 0 {
		inputValidator = validator[0]
	}

	var input string
	var intInput int
	var err error
	var done = false
	for !done {
		fmt.Println(text)
		fmt.Scanln(&input)
		trimmedInput := strings.TrimSpace(input)
		intInput, err = strconv.Atoi(trimmedInput)
		if err == nil && len(trimmedInput) != 0 && inputValidator(trimmedInput) {
			break
		}
	}
	return intInput
}

// Return a boolean defined by InputValidator
func InputBool(text string, validator InputValidator) bool {
	var input string
	var done = false
	for !done {
		fmt.Println(text)
		fmt.Scanln(&input)
		if len(strings.TrimSpace(input)) != 0 {
			break
		}
	}

	return validator(strings.TrimSpace(input))
}
