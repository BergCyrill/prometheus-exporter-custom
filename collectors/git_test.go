package collectors

import (
	"testing"
)

func TestInjectTokenIntoURL(t *testing.T) {
	url := "https://github.com/org/repo.git"
	token := "ghp_testtoken"
	expected := "https://" + token + "@github.com/org/repo.git"
	result := injectTokenIntoURL(url, token)

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}
