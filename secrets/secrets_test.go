package secrets

import (
	"os"
	"testing"
)

func TestReadSecret(t *testing.T) {
	path := "test_secret.txt"
	expected := "secretvalue"
	os.WriteFile(path, []byte(expected+"\n"), 0644)
	defer os.Remove(path)

	value, err := ReadSecret(path)
	if err != nil {
		t.Fatalf("failed to read secret: %v", err)
	}

	if value != expected {
		t.Errorf("expected '%s', got '%s'", expected, value)
	}
}
