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

type OutputData struct {
	lines      int
	words      int
	characters int
	bytes      int
}

var total OutputData

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
		calculateStats(&filesPath, reader, args)
	} else {
		for _, filePath := range filesPath {
			calculateStats(&filePath, reader, args)
		}
		if len(filesPath) > 0 {
			str := "total"
			printOutput(&str, args, total)
		}
	}
}

func calculateStats(filePath *string, reader io.Reader, args Args) {
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

	var data OutputData

	var lineBuffer []byte

	for scanner.Scan() {
		byte := scanner.Bytes()[0]
		data.bytes++

		if byte == '\n' {
			data.lines++
			lineCharacters := utf8.RuneCount(lineBuffer) + 1 // Add the end of line character
			data.characters += lineCharacters
			data.words += len(strings.Fields(string(lineBuffer)))

			lineBuffer = nil // Reset line buffer for the next line
		} else {
			lineBuffer = append(lineBuffer, byte)
		}
	}

	// handle the last line if it does not end with a newline
	if len(lineBuffer) > 0 {
		lineCharacters := utf8.RuneCount(lineBuffer)
		data.characters += lineCharacters
		data.words += len(strings.Fields(string(lineBuffer)))
	}

	total.lines += data.lines
	total.words += data.words
	total.characters += data.characters
	total.bytes += data.bytes

	printOutput(filePath, args, data)
}

func printOutput(filePath *string, args Args, data OutputData) {
	output := make([]string, 0)

	if args.printLines {
		output = append(output, strconv.Itoa(data.lines))
	}
	if args.printWords {
		output = append(output, strconv.Itoa(data.words))
	}
	if args.printChars {
		output = append(output, strconv.Itoa(data.characters))
	}
	if args.printBytes {
		output = append(output, strconv.Itoa(data.bytes))
	}
	if flag.NFlag() == 0 {
		output = append(output, strconv.Itoa(data.lines))
		output = append(output, strconv.Itoa(data.words))
		output = append(output, strconv.Itoa(data.bytes))
	}
	if *filePath != "" {
		output = append(output, *filePath)
	}
	fmt.Println(strings.Join(output, " "))
}
