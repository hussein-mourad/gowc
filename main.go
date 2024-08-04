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

// Args holds the command-line flags.
type Args struct {
	flags map[string]bool
}

// OutputData stores statistics for a single file.
type OutputData struct {
	file       string // file name
	lines      int    // number of lines
	words      int    // number of words
	characters int    // number of characters
	bytes      int    // number of bytes
}

var (
	outputData []OutputData // Stores statistics for all files
	total      OutputData   // Accumulates totals across all files
	args       Args         // Global variable for command-line flags
)

var maxLinesWidth, maxWordsWidth, maxCharsWidth, maxBytesWidth int // Widths for column formatting

func main() {
	args = parseFlags() // Parse command-line flags

	filesPath := flag.Args() // Get file paths from command-line arguments
	reader := os.Stdin       // Default reader is standard input

	if len(filesPath) == 0 {
		// No file paths provided, read from standard input
		filesPath := ""
		calculateStats(reader, &filesPath)
		printOutput(outputData)
	} else {
		// Process each file path provided
		for _, filePath := range filesPath {
			reader, err := openFile(filePath)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			calculateStats(reader, &filePath)
		}

		if len(filesPath) > 1 {
			// Add a total line if more than one file is processed
			totalArgs := OutputData{file: "total", lines: total.lines, words: total.words, characters: total.characters, bytes: total.bytes}
			outputData = append(outputData, totalArgs)
		}
		printOutput(outputData)
	}
}

// parseFlags parses the command-line flags and returns an Args struct.
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

// openFile opens a file or returns stdin if the file path is empty.
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

// calculateStats calculates statistics for a given reader and file path.
func calculateStats(reader io.Reader, filePath *string) {
	var data OutputData
	var lineBuffer []byte

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes) // Read byte by byte

	for scanner.Scan() {
		byte := scanner.Bytes()[0]
		data.bytes++

		if byte == '\n' {
			data.lines++
			data.characters += utf8.RuneCount(lineBuffer) + 1 // Include newline character
			data.words += len(strings.Fields(string(lineBuffer)))
			lineBuffer = nil // Reset line buffer
		} else {
			lineBuffer = append(lineBuffer, byte)
		}
	}

	// Handle the last line if it does not end with a newline
	if len(lineBuffer) > 0 {
		data.characters += utf8.RuneCount(lineBuffer)
		data.words += len(strings.Fields(string(lineBuffer)))
	}

	// Update totals
	total.lines += data.lines
	total.words += data.words
	total.characters += data.characters
	total.bytes += data.bytes

	// Update maximum widths for formatting
	maxLinesWidth = max(maxLinesWidth, len(strconv.Itoa(data.lines)))
	maxWordsWidth = max(maxWordsWidth, len(strconv.Itoa(data.words)))
	maxCharsWidth = max(maxCharsWidth, len(strconv.Itoa(data.characters)))
	maxBytesWidth = max(maxBytesWidth, len(strconv.Itoa(data.bytes)))

	data.file = *filePath
	outputData = append(outputData, data)
}

// getMaxWidth calculates the maximum width needed for each column based on the flags.
func getMaxWdith() int {
	maxWidth := 0
	if args.flags["l"] {
		maxWidth = max(maxWidth, maxLinesWidth)
	}
	if args.flags["w"] {
		maxWidth = max(maxWidth, maxWordsWidth)
	}
	if args.flags["m"] {
		maxWidth = max(maxWidth, maxCharsWidth)
	}
	if args.flags["c"] {
		maxWidth = max(maxWidth, maxBytesWidth)
	}
	if flag.NFlag() == 0 {
		// No flags provided, default to lines, words, and bytes
		maxWidth = max(maxLinesWidth, maxWordsWidth, maxBytesWidth)
	}
	return maxWidth
}

// printOutput formats and prints the collected output data.
func printOutput(outputData []OutputData) {
	width := getMaxWdith() // Get the maximum width for formatting

	for _, data := range outputData {
		output := []string{}

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
			// Default output format if no flags are provided
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
