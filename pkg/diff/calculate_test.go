package diff

import (
	"context"
	"testing"
)

func TestCalculate_Must_Have_Lines(t *testing.T) {
	diffs, err := Calculate(context.Background(), ".", "4ce756ca74a9d94637312910d51c570888447157", "edf7a69bf3fa8277416895dcdf1b4dd9afaf8b81")
	if err != nil {
		t.Fatal(err)
	}

	for _, diff := range diffs {
		if diff.Lines == nil {
			t.Fatalf("%s", diff.Src.Path)
		}
	}
}
