package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

var files = []string{"test_data/test1.txt", "test_data/test2.txt", "test_data/test3.txt"}

func runFileTest(t *testing.T, args []string, expectedOutputs []string) {
	cmdArgs := append([]string{"run", "main.go"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	for _, expectedOutput := range expectedOutputs {
		if !strings.Contains(outputStr, expectedOutput) {
			t.Errorf("Expected %s but got %s", expectedOutput, outputStr)
		}
	}
}

func TestNumberOfBytes(t *testing.T) {
	args := append([]string{"-c"}, files...)
	expectedOutputs := []string{
		"5 test_data/test1.txt\n",
		"12 test_data/test2.txt\n",
		"335045 test_data/test3.txt\n",
		"335062 total\n",
	}
	runFileTest(t, args, expectedOutputs)
}

func TestNumberOfLines(t *testing.T) {
	args := append([]string{"-l"}, files...)
	expectedOutputs := []string{
		"1 test_data/test1.txt\n",
		"2 test_data/test2.txt\n",
		"7145 test_data/test3.txt\n",
		"7148 total\n",
	}
	runFileTest(t, args, expectedOutputs)
}

func TestNumberOfWords(t *testing.T) {
	args := append([]string{"-w"}, files...)
	expectedOutputs := []string{
		"1 test_data/test1.txt\n",
		"2 test_data/test2.txt\n",
		"58164 test_data/test3.txt\n",
		"58167 total\n",
	}
	runFileTest(t, args, expectedOutputs)
}

func TestNumberOfCharacters(t *testing.T) {
	args := []string{"-m", "test_data/test3.txt"}
	expectedOutputs := []string{
		"332147 test_data/test3.txt\n",
	}
	runFileTest(t, args, expectedOutputs)
}

func TestDefault(t *testing.T) {
	args := append([]string{}, files...)
	expectedOutputs := []string{
		"     1      1      5 test_data/test1.txt\n",
		"     2      2     12 test_data/test2.txt\n",
		"  7145  58164 335045 test_data/test3.txt\n",
		"  7148  58167 335062 total\n",
	}
	runFileTest(t, args, expectedOutputs)
}

func TestLinesFromStdin(t *testing.T) {
	catCmd := exec.Command("cat", "test_data/test3.txt")
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

	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected %s but got %s", expectedOutput, output)
	}
}
