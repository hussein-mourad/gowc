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
	flags map[string]bool
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
	args := parseFlags()

	filesPath := flag.Args()
	reader := os.Stdin

	if len(filesPath) == 0 {
		filesPath := ""
		calculateStats(reader, &filesPath)
		printOutput(args, outputData)
	} else {
		for _, filePath := range filesPath {
			reader, err := openFile(filePath)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			calculateStats(reader, &filePath)
		}

		if len(filesPath) > 1 {
			totalArgs := OutputData{file: "total", lines: total.lines, words: total.words, characters: total.characters, bytes: total.bytes}
			outputData = append(outputData, totalArgs)
		}
		printOutput(args, outputData)
	}
}

func parseFlags() Args {
	printBytes := flag.Bool("c", false, "print the byte counts")
	printChars := flag.Bool("m", false, "print the character counts")
	printWords := flag.Bool("w", false, "print the word counts")
	printLines := flag.Bool("l", false, "print the newline counts")
	flag.Parse()

	return Args{
		flags: map[string]bool{
			"c": *printBytes,
			"m": *printChars,
			"w": *printWords,
			"l": *printLines,
		},
	}
}

func openFile(filePath string) (io.Reader, error) {
	if filePath == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func calculateStats(reader io.Reader, filePath *string) {
	var data OutputData
	var lineBuffer []byte

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

		if args.flags["l"] {
			output = append(output, fmt.Sprintf("%*d", width, data.lines))
		}
		if args.flags["w"] {
			output = append(output, fmt.Sprintf("%*d", width, data.words))
		}
		if args.flags["m"] {
			output = append(output, fmt.Sprintf("%*d", width, data.characters))
		}
		if args.flags["c"] {
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
