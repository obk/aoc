package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// runDay builds puzzle p's solution and runs it in its year folder, where its
// input path points, then prints how long it ran. With input set, the
// solution gets it as its first argument. The build isn't timed.
func runDay(l layout, p Puzzle, input string, stdout, stderr io.Writer) error {
	src := l.source(p)
	if !exists(src) {
		return fmt.Errorf("there's no %s yet; `aoc new` creates it", l.rel(src))
	}
	tmp, err := os.MkdirTemp("", "aoc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	bin := filepath.Join(tmp, "day")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	build := exec.Command("go", "build", "-o", bin, filepath.Base(src))
	build.Dir, build.Stdout, build.Stderr = l.dir(p), stderr, stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("building %s: %v", l.rel(src), err) // %v: a failed build isn't the solution's exit status
	}

	var args []string
	if input != "" {
		abs, err := filepath.Abs(input) // the solution runs in the year folder
		if err != nil {
			return err
		}
		args = append(args, abs)
	}
	cmd := exec.Command(bin, args...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = l.dir(p), stdout, stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	// The clock starts once the process is up: on Windows, starting a freshly
	// built program waits most of a second for the virus scanner.
	start := time.Now()
	err = cmd.Wait()
	fmt.Fprintf(stderr, "%v ran in %v\n", p, roundDuration(time.Since(start)))
	return err
}

// roundDuration keeps durations readable: 412µs, 12.345ms, 2.345s.
func roundDuration(d time.Duration) time.Duration {
	if d < time.Second {
		return d.Round(time.Microsecond)
	}
	return d.Round(time.Millisecond)
}
