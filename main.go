package main

import (
	"fmt"
	"os"
)

var programName = os.Args[0]

func getBytesCount(filename string) int {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return len(data)
}

func main() {
	args := os.Args[1:]
	// flag := args[0]
	filename := args[1]
	fmt.Println(getBytesCount(filename), filename)
}
