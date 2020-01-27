package main

import (
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	newList := []string{"1", "2", "3", "4", "2", "3"}
	oldList := []string{"4", "5"}
	list := removeDuplicates(newList, oldList)
	result := []string{"1", "2", "3"}
	for i := range list {
		if list[i] != result[i] {
			t.Errorf("Got %v error, want %v", list[i], result[i])
		}
	}
}
