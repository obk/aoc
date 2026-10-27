package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanSession(t *testing.T) {
	for in, want := range map[string]string{
		"53616c7465645f5f\n":          "53616c7465645f5f",
		"  session=53616C7465645F5F ": "53616C7465645F5F",
		"":                            "",
		"not hex":                     "",
		"session=":                    "",
		strings.Repeat("a", 600):      "",
	} {
		got, err := cleanSession(in)
		if got != want || (want == "") != (err != nil) {
			t.Errorf("cleanSession(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestLogin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aoc", "session")
	var out bytes.Buffer
	if err := login(strings.NewReader("session=c0ffee\n"), &out, path); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(path); string(b) != "c0ffee\n" {
		t.Errorf("saved %q", b)
	}
	if err := login(strings.NewReader("oops\n"), &out, path+"2"); err == nil || exists(path+"2") {
		t.Errorf("saved a bad cookie (err = %v)", err)
	}
}

func TestLoadSessionFromEnv(t *testing.T) {
	t.Setenv("AOC_SESSION", "session=c0ffee")
	if s, err := loadSession(); s != "c0ffee" || err != nil {
		t.Errorf("loadSession = %q, %v", s, err)
	}
}
