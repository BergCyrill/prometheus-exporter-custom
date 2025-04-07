package collectors

import (
	"os/exec"
	"testing"
)

func TestDockerCommandMock(t *testing.T) {
	// Just a basic test that Command can be intercepted
	cmd := exec.Command("echo", "docker pull test-image")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("error executing mock command: %v", err)
	}

	expected := "docker pull test-image\n"
	if string(out) != expected {
		t.Errorf("expected '%s', got '%s'", expected, out)
	}
}
