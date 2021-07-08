package tfinstall

import (
	"context"
	"testing"
)

func TestListVersions(t *testing.T) {
	c, err := ListVersions(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(c) == 0 {
		t.Fatal("didn't find any versions when we expected to")
	}
}
