package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderDefault(t *testing.T) {
	l := layout{t.TempDir()}
	src, err := renderDay(l, dayData{2024, 3, "Pebbles & Pails"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"// Day 3: Pebbles & Pails\n", `inputPath     = "inputs/3.txt"`, "func day3(path string) error {", "if err := day3(path); err != nil {"} {
		if !bytes.Contains(src, []byte(want)) {
			t.Errorf("missing %q in:\n%s", want, src)
		}
	}
	untitled, err := renderDay(l, dayData{2024, 3, ""})
	if err != nil || !bytes.HasPrefix(untitled, []byte("// Day 3\n")) {
		t.Errorf("without a title: %v\n%s", err, untitled)
	}
}

func TestRenderCustom(t *testing.T) {
	l := layout{t.TempDir()}
	if err := writeNew(l.template(), []byte("package main\n// {{.Year}}/{{.Day}}\nfunc main()   {}\n")); err != nil {
		t.Fatal(err)
	}
	src, err := renderDay(l, dayData{2024, 3, ""})
	if want := "package main\n\n// 2024/3\nfunc main() {}\n"; err != nil || string(src) != want {
		t.Errorf("renderDay = %q, %v; want %q (formatted)", src, err, want)
	}

	if err := os.WriteFile(l.template(), []byte("not go {{.Day}}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := renderDay(l, dayData{2024, 3, ""}); err == nil {
		t.Error("a template that isn't Go should fail")
	}
}

// TestNewAndRun sets up a day from a fake site with the default template,
// then builds and runs it, on the input and on the example.
func TestNewAndRun(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("needs the go command")
	}
	site, c := newFakeSite(t, map[string]string{"/2024/day/3": puzzlePage, "/2024/day/3/input": "1 2\n3 4\n"})
	l := layout{t.TempDir()}
	p := Puzzle{2024, 3}

	var out bytes.Buffer
	if err := newDay(l, c, p, &out); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{l.source(p), l.example(p), l.input(p)} {
		if !exists(path) {
			t.Errorf("%s is missing; output:\n%s", path, out.String())
		}
	}

	// Again: nothing new to fetch or write.
	out.Reset()
	if err := newDay(l, c, p, &out); err != nil {
		t.Fatal(err)
	}
	if n := site.hits.Load(); n != 2 {
		t.Errorf("%d requests in all, want 2 (page and input)", n)
	}
	if !strings.Contains(out.String(), "already") || strings.Contains(out.String(), "created") {
		t.Errorf("second run said:\n%s", out.String())
	}

	for _, input := range []string{"", l.example(p)} {
		var stdout, stderr bytes.Buffer
		if err := runDay(l, p, input, &stdout, &stderr); err != nil {
			t.Fatalf("runDay(%q): %v\n%s", input, err, stderr.String())
		}
		if stdout.String() != "Part 1: 0\nPart 2: 0\n" {
			t.Errorf("runDay(%q) printed %q", input, stdout.String())
		}
		if !strings.Contains(stderr.String(), "2024 day 3 ran in ") {
			t.Errorf("runDay(%q) didn't time it: %q", input, stderr.String())
		}
	}
}

func TestRunFailures(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("needs the go command")
	}
	l := layout{t.TempDir()}
	exits := Puzzle{2024, 1}
	if err := writeNew(l.source(exits), []byte("package main\n\nimport \"os\"\n\nfunc main() { os.Exit(3) }\n")); err != nil {
		t.Fatal(err)
	}
	var exitErr *exec.ExitError
	if err := runDay(l, exits, "", &bytes.Buffer{}, &bytes.Buffer{}); !errors.As(err, &exitErr) || exitErr.ExitCode() != 3 {
		t.Errorf("a solution that exits with 3: %v", err)
	}

	broken := Puzzle{2024, 2}
	if err := writeNew(l.source(broken), []byte("package main\n\nfunc main() { x }\n")); err != nil {
		t.Fatal(err)
	}
	err := runDay(l, broken, "", &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || errors.As(err, &exitErr) || !strings.Contains(err.Error(), "building "+filepath.Join("2024", "2.go")) {
		t.Errorf("a solution that doesn't build: %v", err)
	}

	if err := runDay(l, Puzzle{2024, 5}, "", &bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("running a day that doesn't exist should fail")
	}
}
