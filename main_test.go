package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func runFileTest(t *testing.T, cmd *exec.Cmd, expectedOutput string) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}

func TestNumberOfBytes(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-c", "test_data/test.txt")
	expectedOutput := "342190 test_data/test.txt\n"
	runFileTest(t, cmd, expectedOutput)
}

func TestNumberOfLines(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-l", "test_data/test.txt")
	expectedOutput := "7145 test_data/test.txt\n"
	runFileTest(t, cmd, expectedOutput)
}

func TestNumberOfWords(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-w", "test_data/test.txt")
	expectedOutput := "58164 test_data/test.txt\n"
	runFileTest(t, cmd, expectedOutput)
}

func TestNumberOfCharacters(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-m", "test_data/test.txt")
	expectedOutput := "339292 test_data/test.txt\n"
	runFileTest(t, cmd, expectedOutput)
}

func TestDefault(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "test_data/test.txt")
	expectedOutput := "7145 58164 342190 test_data/test.txt\n"
	runFileTest(t, cmd, expectedOutput)
}

func TestLinesFromStdin(t *testing.T) {
	catCmd := exec.Command("cat", "test_data/test.txt")
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
