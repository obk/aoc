package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// est is the time zone Advent of Code runs on: puzzles unlock at midnight
// UTC-5, and December has no daylight saving time.
var est = time.FixedZone("EST", -5*60*60)

const firstYear = 2015

// Puzzle is one day of one event.
type Puzzle struct {
	Year, Day int
}

func (p Puzzle) String() string { return fmt.Sprintf("%d day %d", p.Year, p.Day) }

// path is the puzzle page's URL path; the input is at path + "/input".
func (p Puzzle) path() string { return fmt.Sprintf("/%d/day/%d", p.Year, p.Day) }

// unlock is when the puzzle comes out.
func (p Puzzle) unlock() time.Time {
	return time.Date(p.Year, time.December, p.Day, 0, 0, 0, 0, est)
}

// parsePuzzle picks the puzzle from the command's numbers: none for today's,
// a day alone, or a year and a day. A day alone is of folderYear when aoc runs
// in a year folder, else of the latest event.
func parsePuzzle(nums []string, now time.Time, folderYear int) (Puzzle, error) {
	now = now.In(est)
	latest := now.Year() // the latest event that has started
	if now.Month() != time.December {
		latest--
	}
	n := make([]int, len(nums))
	for i, s := range nums {
		var err error
		if n[i], err = strconv.Atoi(s); err != nil {
			return Puzzle{}, fmt.Errorf("%q is not a year or a day", s)
		}
	}

	var p Puzzle
	switch len(n) {
	case 0:
		if now.Month() != time.December || now.Day() > 25 {
			return Puzzle{}, errors.New("there's no puzzle today; say which one, like `2025 3`")
		}
		p = Puzzle{now.Year(), now.Day()}
	case 1:
		p = Puzzle{latest, n[0]}
		if folderYear != 0 {
			p.Year = folderYear
		}
	case 2:
		p = Puzzle{n[0], n[1]}
	default:
		return Puzzle{}, errors.New("give at most a year and a day, like `2025 3`")
	}
	if p.Year < firstYear || p.Year > latest {
		return Puzzle{}, fmt.Errorf("year %d is outside %d..%d", p.Year, firstYear, latest)
	}
	// Events since 2025 have 12 days; the site answers "not found" for the rest.
	if p.Day < 1 || p.Day > 25 {
		return Puzzle{}, fmt.Errorf("day %d is outside 1..25", p.Day)
	}
	return p, nil
}

// layout is where a repository keeps its files: YYYY/D.go, with the inputs and
// examples in YYYY/inputs/.
type layout struct{ root string }

// findLayout makes dir the root, or its parent when dir is a year folder, so
// aoc works from both. year is the folder's year, or 0.
func findLayout(dir string) (l layout, year int) {
	if n, err := strconv.Atoi(filepath.Base(dir)); err == nil && n >= firstYear && n < 10_000 {
		return layout{filepath.Dir(dir)}, n
	}
	return layout{dir}, 0
}

func (l layout) dir(p Puzzle) string    { return filepath.Join(l.root, strconv.Itoa(p.Year)) }
func (l layout) source(p Puzzle) string { return filepath.Join(l.dir(p), fmt.Sprintf("%d.go", p.Day)) }
func (l layout) input(p Puzzle) string {
	return filepath.Join(l.dir(p), "inputs", fmt.Sprintf("%d.txt", p.Day))
}
func (l layout) example(p Puzzle) string {
	return filepath.Join(l.dir(p), "inputs", fmt.Sprintf("%d.example.txt", p.Day))
}

// template is the repository's own template for new days, if it has one.
func (l layout) template() string { return filepath.Join(l.root, ".aoc", "day.go.tmpl") }

// rel shortens path for messages, like "2024/inputs/3.txt".
func (l layout) rel(path string) string {
	if r, err := filepath.Rel(l.root, path); err == nil {
		return r
	}
	return path
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// writeNew creates the file at path with data, and its folder if needed. It
// never replaces a file that's already there.
func writeNew(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(path) // don't leave half a file behind
	}
	return err
}
