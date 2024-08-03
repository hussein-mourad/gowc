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
