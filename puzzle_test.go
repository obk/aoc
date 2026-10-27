package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParsePuzzle(t *testing.T) {
	december := time.Date(2026, 12, 5, 3, 0, 0, 0, time.UTC) // still December 4 in EST
	october := time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		nums   []string
		now    time.Time
		folder int    // the year folder aoc runs in, or 0
		want   Puzzle // zero for an error
	}{
		{nil, december, 0, Puzzle{2026, 4}},
		{[]string{"2"}, december, 0, Puzzle{2026, 2}},
		{[]string{"2015", "25"}, december, 0, Puzzle{2015, 25}},
		{nil, october, 0, Puzzle{}},                     // no puzzle today
		{[]string{"7"}, october, 0, Puzzle{2025, 7}},    // the latest event is 2025's
		{[]string{"7"}, october, 2024, Puzzle{2024, 7}}, // in the 2024 folder
		{[]string{"2023", "7"}, october, 2024, Puzzle{2023, 7}},
		{[]string{"7"}, october, 2026, Puzzle{}},        // a folder for an event that hasn't started
		{[]string{"2026", "1"}, october, 0, Puzzle{}},   // not started yet
		{[]string{"2014", "1"}, october, 0, Puzzle{}},   // before the first event
		{[]string{"2024", "26"}, october, 0, Puzzle{}},  // no such day
		{[]string{"2024", "0"}, october, 0, Puzzle{}},   // no such day
		{[]string{"x"}, october, 0, Puzzle{}},           // not a number
		{[]string{"1", "2", "3"}, october, 0, Puzzle{}}, // too many
	}
	for _, tt := range tests {
		got, err := parsePuzzle(tt.nums, tt.now, tt.folder)
		if tt.want == (Puzzle{}) {
			if err == nil {
				t.Errorf("parsePuzzle(%q) at %v in %d = %v, want an error", tt.nums, tt.now, tt.folder, got)
			}
			continue
		}
		if err != nil || got != tt.want {
			t.Errorf("parsePuzzle(%q) at %v in %d = %v, %v; want %v", tt.nums, tt.now, tt.folder, got, err, tt.want)
		}
	}
}

func TestFindLayout(t *testing.T) {
	root := filepath.Join("home", "advent-of-code")
	for dir, wantYear := range map[string]int{root: 0, filepath.Join(root, "2024"): 2024} {
		if l, year := findLayout(dir); l.root != root || year != wantYear {
			t.Errorf("findLayout(%q) = %q, %d; want %q, %d", dir, l.root, year, root, wantYear)
		}
	}
	l := layout{root}
	p := Puzzle{2024, 3}
	if got, want := l.rel(l.input(p)), filepath.Join("2024", "inputs", "3.txt"); got != want {
		t.Errorf("input = %q, want %q", got, want)
	}
}

func TestWriteNew(t *testing.T) {
	path := filepath.Join(t.TempDir(), "2024", "3.go")
	if err := writeNew(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := writeNew(path, []byte("second")); !errors.Is(err, fs.ErrExist) {
		t.Errorf("overwriting: err = %v, want fs.ErrExist", err)
	}
	if b, _ := os.ReadFile(path); string(b) != "first" {
		t.Errorf("file holds %q", b)
	}
}
