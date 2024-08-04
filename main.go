package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var programName = os.Args[0]

func getBytesCount(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		file.Close()
		os.Exit(1)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return len(data)
}

func getCharactersCount(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		file.Close()
		os.Exit(1)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	content := string(data)

	charCount := utf8.RuneCountInString(content)
	return charCount
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
	reader := os.Stdin

	if filename != "" {
		r, err := os.Open(filename)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		reader = r
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	lines := 0
	words := 0
	characters := 0
	bytes := 0
	var lineBuffer []byte

	for scanner.Scan() {
		byte := scanner.Bytes()[0]
		bytes++

		if byte == '\n' {
			lines++
			lineCharacters := utf8.RuneCount(lineBuffer) + 1 // Add the end of line character
			characters += lineCharacters
			words += len(strings.Fields(string(lineBuffer)))

			lineBuffer = nil // Reset line buffer for the next line
		} else {
			lineBuffer = append(lineBuffer, byte)
		}
	}

	// handle the last line if it does not end with a newline
	if len(lineBuffer) > 0 {
		lineCharacters := utf8.RuneCount(lineBuffer)
		characters += lineCharacters
		words += len(strings.Fields(string(lineBuffer)))
	}

	output := make([]string, 0)

	if *printLines {
		output = append(output, strconv.Itoa(lines))
	}
	if *printWords {
		output = append(output, strconv.Itoa(words))
	}
	if *printChars {
		output = append(output, strconv.Itoa(characters))
	}
	if *printBytes {
		output = append(output, strconv.Itoa(bytes))
	}
	if flag.NFlag() == 0 {
		output = append(output, strconv.Itoa(lines))
		output = append(output, strconv.Itoa(words))
		output = append(output, strconv.Itoa(bytes))
	}
	if filename != "" {
		output = append(output, filename)
	}
	fmt.Println(strings.Join(output, " "))
}
