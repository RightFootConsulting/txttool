package textp

import (
	"testing"

	_ "github.com/flashlabs/rootpath"
)

func TestJustifyText(t *testing.T) {
	err := JustifyText("testInput.txt", "testOutput.txt", 120)
	if err != nil {
		t.Fatalf("failure: %v", err)
	}
}
