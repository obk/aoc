package main

import "testing"

func TestPageTitle(t *testing.T) {
	if got := pageTitle([]byte(puzzlePage)); got != "Pebbles & Pails" {
		t.Errorf("pageTitle = %q", got)
	}
	if got := pageTitle([]byte("<h2>Something else</h2>")); got != "" {
		t.Errorf("pageTitle without a title = %q", got)
	}
}

func TestPageExample(t *testing.T) {
	if got, want := pageExample([]byte(puzzlePage)), "3 4\n1 2 <- here\n"; got != want {
		t.Errorf("pageExample = %q, want %q", got, want)
	}
	if got := pageExample([]byte("<p>Only <code>inline</code> code.</p>")); got != "" {
		t.Errorf("pageExample without a block = %q", got)
	}
}
