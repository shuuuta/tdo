package utils

import "testing"

func TestHashPath(t *testing.T) {
	p1 := "/foo/bar"
	h1, _ := HashPath(p1)
	if p1 == h1 {
		t.Fatalf("expect not equal path and hash %q", h1)
	}

	p2 := "/foo/bar/"
	h2, _ := HashPath(p2)

	if h1 != h2 {
		t.Fatalf("expect equal %q and %q", h1, h2)
	}

	if len(h1) != 12 {
		t.Fatalf("expect hash length 12, got %d: %v", len(h1), h1)
	}
}

func TestHashPathGotEmptyPath(t *testing.T) {
	h1, err := HashPath("")
	if err != nil {
		t.Fatal(err)
	}
	if h1 == "" {
		t.Fatalf("return empty text")
	}
}
