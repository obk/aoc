// Command aoc fetches Advent of Code puzzles and inputs and runs Go solutions.
// It works in a repository laid out like github.com/obk/advent-of-code: one
// file per day in YYYY/D.go, with the inputs in YYYY/inputs/D.txt.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

// version is set with -ldflags "-X main.version=v0.1.0", or comes from the
// module version when installed with go install.
var version = "dev"

const usage = `aoc fetches Advent of Code inputs and runs Go solutions laid out as
YYYY/D.go, with the inputs in YYYY/inputs/. Run it in the repository or a year folder.

Usage:
  aoc login                          save your session cookie (inputs need it)
  aoc new   [[year] day]             create YYYY/D.go, save the example, download the input
  aoc input [[year] day]             download the input to YYYY/inputs/D.txt
  aoc run   [[year] day] [-e|-i file]  build and run a day, and time it
                                     (-e: on the example, -i: on another input)
  aoc version

Without a year and day it's today's puzzle, in December. A day alone is from the
year folder you're in, or else from the latest event.
`

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	err := command(os.Args[1], os.Args[2:])
	var exitErr *exec.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exitErr):
		os.Exit(exitErr.ExitCode()) // the solution has already said what went wrong
	default:
		fmt.Fprintln(os.Stderr, "aoc:", err)
		os.Exit(1)
	}
}

func command(name string, args []string) error {
	switch name {
	case "version", "-version", "--version":
		fmt.Println("aoc", version)
		return nil
	case "help", "-h", "-help", "--help":
		fmt.Print(usage)
		return nil
	case "login":
		path, err := sessionPath()
		if err != nil {
			return err
		}
		return login(os.Stdin, os.Stdout, path)
	case "new", "input", "run":
	default:
		return fmt.Errorf("unknown command %q; see `aoc help`", name)
	}

	nums, input, example := args, "", false
	if name == "run" {
		var err error
		if nums, input, example, err = splitRunArgs(args); err != nil {
			return err
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	l, folderYear := findLayout(cwd)
	p, err := parsePuzzle(nums, time.Now(), folderYear)
	if err != nil {
		return err
	}

	if name == "run" {
		if example {
			if input = l.example(p); !exists(input) {
				return fmt.Errorf("there's no example for %v; `aoc new` saves it to %s", p, l.rel(input))
			}
		}
		return runDay(l, p, input, os.Stdout, os.Stderr)
	}
	session, err := loadSession()
	if err != nil {
		return err
	}
	c := newClient(session)
	if name == "new" {
		return newDay(l, c, p, os.Stdout)
	}
	return fetchInput(l, c, p, os.Stdout)
}

// splitRunArgs separates the year and day from run's flags, which may come
// before or after them.
func splitRunArgs(args []string) (nums []string, input string, example bool, err error) {
	for i := 0; i < len(args); i++ {
		switch a := args[i]; a {
		case "-e", "-example":
			example = true
		case "-i", "-input":
			if i+1 == len(args) {
				return nil, "", false, errors.New("-i needs a file")
			}
			i++
			input = args[i]
		default:
			if strings.HasPrefix(a, "-") {
				return nil, "", false, fmt.Errorf("unknown flag %s; see `aoc help`", a)
			}
			nums = append(nums, a)
		}
	}
	if example && input != "" {
		return nil, "", false, errors.New("use -e or -i, not both")
	}
	return nums, input, example, nil
}
