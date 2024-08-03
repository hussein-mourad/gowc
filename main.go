package main

import (
	"bufio"
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

func getLinesCount(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		file.Close()
		os.Exit(1)
	}
	defer file.Close()
	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}
	return lines
}

func getWordsCount(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		file.Close()
		os.Exit(1)
	}
	defer file.Close()
	words := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words++
	}
	return words
}

func main() {
	args := os.Args[1:]
	flag := ""
	filename := args[0]
	if len(args) == 2 {
		flag = args[0]
		filename = args[1]
	}

	switch flag {
	case "-c":
		fmt.Println(getBytesCount(filename), filename)
	case "-l":
		fmt.Println(getLinesCount(filename), filename)
	case "-w":
		fmt.Println(getWordsCount(filename), filename)
	default:
		fmt.Println(getLinesCount(filename), getWordsCount(filename), getBytesCount(filename), filename)
	}
}
