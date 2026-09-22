package lib_test

import (
	"testing"

	"github.com/3l4l5/blog/src/lib"
)

func TestDivideTopDir(t *testing.T) {
	topDir, theOther := lib.DivideTopDir("aaa/bbb/ccc/ddd.md")

	topDirWant := "aaa"
	if topDir != topDirWant {
		t.Fatalf("got %s, want %s", topDir, topDirWant)
	}

	baseDirWant := "bbb/ccc/ddd.md"
	if theOther != baseDirWant {
		t.Fatalf("got %s, want %s", theOther, baseDirWant)
	}
}
