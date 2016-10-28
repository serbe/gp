package main

import (
	"bufio"
	"fmt"
	"os"
)

func writeSlice(slice []string, filename string) error {
	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range slice {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func readSlice(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}

func existsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return createFile(file)
	}
	return true
}

func createFile(file string) bool {
	_, err := os.Create(file)
	if err != nil {
		return false
	}
	return true
}
