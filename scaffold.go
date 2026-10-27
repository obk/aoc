package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"text/template"
)

//go:embed day.go.tmpl
var defaultTemplate string

// dayData is what a template can use: {{.Year}}, {{.Day}} and {{.Title}}.
type dayData struct {
	Year, Day int
	Title     string // "" if the page didn't have one
}

// newDay sets up puzzle p: the source file from the template, the example
// from the puzzle page and the input. It skips what's already there, and only
// asks for the page when it needs it.
func newDay(l layout, c *client, p Puzzle, out io.Writer) error {
	src, example := l.source(p), l.example(p)
	if !exists(src) || !exists(example) {
		page, err := c.get(p, p.path())
		if err != nil {
			return err
		}
		if !exists(src) {
			code, err := renderDay(l, dayData{p.Year, p.Day, pageTitle(page)})
			if err != nil {
				return err
			}
			if err := writeNew(src, code); err != nil {
				return err
			}
			fmt.Fprintf(out, "created %s\n", l.rel(src))
		}
		if ex := pageExample(page); ex != "" && !exists(example) {
			if err := writeNew(example, []byte(ex)); err != nil {
				return err
			}
			fmt.Fprintf(out, "saved the first code block on the page to %s; check that it's the example\n", l.rel(example))
		}
	}
	if c.session == "" {
		fmt.Fprintln(out, "no session cookie, so no input yet: run `aoc login`, then `aoc input`")
		return nil
	}
	return fetchInput(l, c, p, out)
}

// fetchInput downloads p's input, unless it's already there.
func fetchInput(l layout, c *client, p Puzzle, out io.Writer) error {
	path := l.input(p)
	fetched, err := c.downloadInput(p, path)
	switch {
	case err != nil:
		return err
	case fetched:
		fmt.Fprintf(out, "downloaded the input to %s\n", l.rel(path))
	default:
		fmt.Fprintf(out, "the input is already in %s\n", l.rel(path))
	}
	return nil
}

// renderDay fills in the repository's template, or the default one, and
// formats the result like gofmt.
func renderDay(l layout, data dayData) ([]byte, error) {
	text := defaultTemplate
	custom, err := os.ReadFile(l.template())
	switch {
	case err == nil:
		text = string(custom)
	case !errors.Is(err, fs.ErrNotExist):
		return nil, err
	}
	t, err := template.New("day").Parse(text)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("the template doesn't make valid Go: %w", err)
	}
	return src, nil
}
