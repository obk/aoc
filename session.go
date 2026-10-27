package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// sessionPath is where `aoc login` keeps the session cookie, like
// ~/.config/aoc/session or %AppData%\aoc\session.
func sessionPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "aoc", "session"), nil
}

// loadSession returns the session cookie from $AOC_SESSION or from the file
// `aoc login` wrote, or "" if there's none.
func loadSession() (string, error) {
	if s := os.Getenv("AOC_SESSION"); s != "" {
		return cleanSession(s)
	}
	path, err := sessionPath()
	if err != nil {
		return "", nil // no config folder, so nothing saved
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	s, err := cleanSession(string(b))
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return s, nil
}

// cleanSession checks a session cookie's value, with or without "session=" in
// front.
func cleanSession(s string) (string, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "session=")
	if s == "" || len(s) > 512 || strings.Trim(s, "0123456789abcdefABCDEF") != "" {
		return "", errors.New("that doesn't look like a session cookie, which is a long string of hex digits")
	}
	return s, nil
}

// login asks for the session cookie and saves it to path. It reads the
// cookie from in, not the command line, so it stays out of the shell history.
func login(in io.Reader, out io.Writer, path string) error {
	fmt.Fprint(out, "Log in to adventofcode.com, copy the value of the \"session\" cookie\n"+
		"(browser developer tools, under Storage or Application) and paste it here: ")
	line, err := bufio.NewReader(io.LimitReader(in, 1024)).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	s, err := cleanSession(line)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(s+"\n"), 0o600); err != nil {
		return err
	}
	fmt.Fprintf(out, "saved to %s\n", path)
	return nil
}
