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
	file       string
	lines      int
	words      int
	characters int
	bytes      int
}

var outputData []OutputData

var total OutputData

var maxLinesWidth, maxWordsWidth, maxCharsWidth, maxBytesWidth int

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
		calculateStats(&filesPath, reader)
		printOutput(args, outputData)
	} else {
		for _, filePath := range filesPath {
			calculateStats(&filePath, reader)
		}

		if len(filesPath) > 1 {
			totalArgs := OutputData{file: "total", lines: total.lines, words: total.words, characters: total.characters, bytes: total.bytes}
			outputData = append(outputData, totalArgs)
		}
		printOutput(args, outputData)
	}
}

func calculateStats(filePath *string, reader io.Reader) {
	var data OutputData
	var lineBuffer []byte

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

	// Calculate the total
	total.lines += data.lines
	total.words += data.words
	total.characters += data.characters
	total.bytes += data.bytes

	// Calculate Longest width number with help in formatting the output
	maxLinesWidth = max(maxLinesWidth, len(strconv.Itoa(data.lines)))
	maxWordsWidth = max(maxWordsWidth, len(strconv.Itoa(data.words)))
	maxCharsWidth = max(maxCharsWidth, len(strconv.Itoa(data.characters)))
	maxBytesWidth = max(maxBytesWidth, len(strconv.Itoa(data.bytes)))

	data.file = *filePath
	outputData = append(outputData, data)
}

func printOutput(args Args, outputData []OutputData) {
	width := max(maxLinesWidth, maxWordsWidth, maxCharsWidth, maxBytesWidth)

	for _, data := range outputData {
		output := make([]string, 0)

		if args.printLines {
			output = append(output, fmt.Sprintf("%*d", width, data.lines))
		}
		if args.printWords {
			output = append(output, fmt.Sprintf("%*d", width, data.words))
		}
		if args.printChars {
			output = append(output, fmt.Sprintf("%*d", width, data.characters))
		}
		if args.printBytes {
			output = append(output, fmt.Sprintf("%*d", width, data.bytes))
		}
		if flag.NFlag() == 0 {
			output = append(output, fmt.Sprintf("%*d", width, data.lines))
			output = append(output, fmt.Sprintf("%*d", width, data.words))
			output = append(output, fmt.Sprintf("%*d", width, data.bytes))
		}
		if data.file != "" {
			output = append(output, data.file)
		}
		fmt.Println(strings.Join(output, " "))
	}
}
