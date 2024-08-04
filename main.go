package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func getCharactersCount(filename string) int {
	// TODO: Get characters count
	return getBytesCount(filename)
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
	printBytes := flag.Bool("c", false, "print the byte counts")
	printChars := flag.Bool("m", false, "print the character counts")
	printWords := flag.Bool("w", false, "print the word counts")
	printLines := flag.Bool("l", false, "print the newline counts")

	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("Filename is not defined")
		os.Exit(1)
	}

	output := make([]string, 0)

	if *printLines {
		output = append(output, strconv.Itoa(getLinesCount(filename)))
	}

	if *printWords {
		output = append(output, strconv.Itoa(getWordsCount(filename)))
	}

	if *printChars {
		output = append(output, strconv.Itoa(getCharactersCount(filename)))
	}

	if *printBytes {
		output = append(output, strconv.Itoa(getBytesCount(filename)))
	}

	if flag.NFlag() == 0 {
		output = append(output, strconv.Itoa(getLinesCount(filename)))
		output = append(output, strconv.Itoa(getWordsCount(filename)))
		output = append(output, strconv.Itoa(getBytesCount(filename)))
	}

	output = append(output, filename)

	fmt.Println(strings.Join(output, " "))
}
