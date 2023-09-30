package diff

import (
	"context"
	"testing"
)

func TestCalculate_Must_Have_Patch(t *testing.T) {
	diffs, err := Calculate(context.Background(), ".", "266555338f78735f5182a8b60025ba861df85edc", "bba8505866a7172e1269f5a8cd0ceba1258c7880")
	if err != nil {
		t.Fatal(err)
	}

	for _, diff := range diffs {
		if diff.Patch == nil {
			t.Fatalf("%s", diff.Src.Path)
		}
	}
}
