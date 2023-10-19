package transformer

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/mstgnz/transformer/node"
)

// ErrorHandle
// wrap error
func ErrorHandle(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// Contains
// generic func for array contains any type
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// ConvertToNode
// convert to Node
func ConvertToNode(value any) *node.Node {
	if knot, ok := value.(*node.Node); ok {
		return knot
	}
	return nil
}

// ConvertToSlice
// convert to any slice
func ConvertToSlice(value any) []any {
	slc, _ := value.([]any)
	return slc
}

// StripSpaces
// remove all spaces
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

// PrintErr
// with format
func PrintErr(format string, args ...any) {
	_, err := fmt.Fprintf(os.Stderr, color.RedString("error: "+format+"\n"), args...)
	if err != nil {
		log.Println("error printing failed")
	}
}
