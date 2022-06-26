package transformer

import (
	"log"
	"strings"
	"unicode"
)

// ErrorHandle wrap error
func ErrorHandle(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// Contains generic func for array contains any type
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// ConvertToNode convert to Node
func ConvertToNode(value any) *Node {
	knot, _ := value.(*Node)
	return knot
}

// ConvertToSlice convert to any slice
func ConvertToSlice(value any) []any {
	slc, _ := value.([]any)
	return slc
}

// StripSpaces remove all spaces
func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
