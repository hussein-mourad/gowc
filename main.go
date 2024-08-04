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

type Args struct {
	printBytes bool
	printChars bool
	printWords bool
	printLines bool
}

func main() {
	printBytes := flag.Bool("c", false, "print the byte counts")
	printChars := flag.Bool("m", false, "print the character counts")
	printWords := flag.Bool("w", false, "print the word counts")
	printLines := flag.Bool("l", false, "print the newline counts")

	flag.Parse()

	args := Args{*printBytes, *printChars, *printWords, *printLines}

	filesPath := flag.Args()
	reader := os.Stdin

	if len(filesPath) == 0 {
		filesPath := ""
		printFileOutput(&filesPath, reader, args)
	}

	for _, filePath := range filesPath {
		printFileOutput(&filePath, reader, args)
	}
}

func printFileOutput(filePath *string, reader io.Reader, args Args) {
	if *filePath != "" {
		f, err := os.Open(*filePath)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		reader = f
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

	if args.printLines {
		output = append(output, strconv.Itoa(lines))
	}
	if args.printWords {
		output = append(output, strconv.Itoa(words))
	}
	if args.printChars {
		output = append(output, strconv.Itoa(characters))
	}
	if args.printBytes {
		output = append(output, strconv.Itoa(bytes))
	}
	if flag.NFlag() == 0 {
		output = append(output, strconv.Itoa(lines))
		output = append(output, strconv.Itoa(words))
		output = append(output, strconv.Itoa(bytes))
	}
	if *filePath != "" {
		output = append(output, *filePath)
	}
	fmt.Println(strings.Join(output, " "))
}
