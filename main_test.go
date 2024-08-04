package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestNumberOfBytes(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-c", "test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "342190 test.txt\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestNumberOfLines(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-l", "test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "7145 test.txt\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestNumberOfWords(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-w", "test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "58164 test.txt\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestNumberOfCharacters(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-m", "test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "339292 test.txt\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestDefault(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "7145 58164 342190 test.txt\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestLinesFromStdin(t *testing.T) {
	catCmd := exec.Command("cat", "test.txt")
	goCmd := exec.Command("go", "run", "main.go", "-l")

	var catOut bytes.Buffer
	var goOut bytes.Buffer

	catCmd.Stdout = &catOut
	goCmd.Stdin = &catOut
	goCmd.Stdout = &goOut

	// Start `cat` command
	if err := catCmd.Run(); err != nil {
		t.Fatalf("Failed to run `cat` command: %v", err)
	}

	// Start `go run` command
	if err := goCmd.Run(); err != nil {
		t.Fatalf("Failed to run `go run` command: %v", err)
	}

	// Check the output
	output := goOut.String()
	expectedOutput := "7145\n"

	if output != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, output)
	}
}
