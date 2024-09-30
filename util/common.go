package util

import (
	"math/rand"
	"reflect"
)

// ExistsInSlide check if item exists in slice
func Contains[T comparable](arr []T, item T) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, item) {
			return true
		}
	}
	return false
}

// UniqueSlide remove duplicate items in slice
func RemoveDuplicates(input []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueSlide := []string{}

	for _, item := range input {
		// Check if the item is unique
		if !uniqueMap[item] {
			uniqueSlide = append(uniqueSlide, item)
			uniqueMap[item] = true
		}
	}

	return uniqueSlide
}

func Intersect[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	for _, item := range a {
		m[item] = true
	}

	var result []T
	for _, item := range b {
		if m[item] {
			result = append(result, item)
		}
	}
	return result
}

func RandomString(length int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
