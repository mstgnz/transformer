package main

import (
	"log"
)

func errorHandle(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// generic func for array contains any type
func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// convert to node
func convertToNode(value any) *node {
	knot, _ := value.(*node)
	return knot
}

// convert to any slice
func convertToSlice(value any) []any {
	slc, _ := value.([]any)
	return slc
}
