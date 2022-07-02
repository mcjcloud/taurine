package util

import (
	"testing"
)

func TestImportGraph(t *testing.T) {
	g := NewImportGraph("a")

	g.Add("a", "b")
	g.Add("b", "c")
	g.Add("a", "c")

	// test g.Nodes
	if a, ok := g.Nodes["a"]; !ok {
		t.Error("graph missing node 'a'")
	} else if a.Path != "a" {
		t.Error("node a's path does not match key")
	} else if len(a.Children) != 2 {
		t.Error("expected node a to have 2 children")
	}

	if b, ok := g.Nodes["b"]; !ok {
		t.Error("graph missing node 'b'")
	} else if b.Path != "b" {
		t.Error("node b's path does not match key")
	} else if len(b.Children) != 1 {
		t.Error("expected node b to have 1 children")
	}

	if c, ok := g.Nodes["c"]; !ok {
		t.Error("graph missing node 'c'")
	} else if c.Path != "c" {
		t.Error("node c's path does not match key")
	} else if len(c.Children) != 0 {
		t.Error("expected node b to have 0 children")
	}

	// test g.Matrix
	if a, ok := g.Matrix["a"]; !ok || a[0] != "b" || a[1] != "c" {
		t.Errorf("matrix entry for 'a' should be [b, c] but was %s", a)
	}
	if b, ok := g.Matrix["b"]; !ok || b[0] != "c" {
		t.Errorf("matrix entry for 'b' should be [c] but was %s", b)
	}
	if c, ok := g.Matrix["c"]; !ok || len(c) != 0 {
		t.Errorf("matrix entry for 'c' should be [] but was %s", c)
	}
}
