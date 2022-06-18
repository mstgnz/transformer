package main

import (
	"log"
	"strings"
	"unicode"
)

// errorHandle
func errorHandle(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// contains generic func for array contains any type
func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// convertToNode convert to node
func convertToNode(value any) *node {
	knot, _ := value.(*node)
	return knot
}

// convertToSlice convert to any slice
func convertToSlice(value any) []any {
	slc, _ := value.([]any)
	return slc
}

// stripSpaces
func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
