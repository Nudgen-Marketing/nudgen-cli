package cmd

import (
	"bytes"
	"testing"
)

func TestRootHelp(t *testing.T) {
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error from help command, got %v", err)
	}

	out := b.String()
	if !contains(out, "Nudgen is an interactive CLI") {
		t.Errorf("Expected help message, got %s", out)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
