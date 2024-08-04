package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// countWords counts the number of words in a string
func countWords(s string) int {
	var words int
	inWord := false

	for _, r := range s {
		if unicode.IsSpace(r) {
			if inWord {
				words++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
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

	lines := 0
	words := 0
	characters := 0
	bytes := 0

	for scanner.Scan() {
		line := scanner.Text()
		lines++

		// Count bytes and characters
		lineBytes := len(line) + 1
		lineCharacters := utf8.RuneCountInString(line) + 1

		bytes += lineBytes
		characters += lineCharacters

		// Count words
		words += countWords(line)
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
