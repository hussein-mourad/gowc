package main

import (
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

func TestStdin(t *testing.T) {
	cmd := exec.Command("cat", "test.txt", "|", "go", "run", "main.go", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "7145\n"

	if string(output) != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, string(output))
	}
}
