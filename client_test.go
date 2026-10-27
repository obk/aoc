package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// puzzlePage is a made-up puzzle page, shaped like the real ones.
const puzzlePage = `<!DOCTYPE html>
<html lang="en-us">
<head><title>Day 3 - Advent of Code 2024</title></head>
<body><main>
<article class="day-desc"><h2>--- Day 3: Pebbles &amp; Pails ---</h2><p>The elves count pebbles.</p>
<p>For example:</p>
<pre><code>3 4
<em>1</em> 2 &lt;- here
</code></pre>
<p>Then use <code>inline code</code> and another block:</p>
<pre><code>not the example</code></pre>
</article>
</main></body>
</html>
`

// fakeSite is a stand-in for adventofcode.com that serves pages by URL path.
// Inputs need the session cookie "c0ffee", as on the real site.
type fakeSite struct {
	hits atomic.Int32
}

// newFakeSite starts a fake site and returns a client for it, logged in, on
// 2026-12-20.
func newFakeSite(t *testing.T, pages map[string]string) (*fakeSite, *client) {
	t.Helper()
	site := &fakeSite{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site.hits.Add(1)
		if ua := r.UserAgent(); !strings.HasPrefix(ua, "github.com/obk/aoc") {
			t.Errorf("User-Agent %q doesn't name the tool", ua)
		}
		if strings.HasSuffix(r.URL.Path, "/input") {
			if ck, err := r.Cookie("session"); err != nil || ck.Value != "c0ffee" {
				http.Error(w, "Puzzle inputs differ by user.  Please log in to get your puzzle input.", http.StatusBadRequest)
				return
			}
		}
		body, ok := pages[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)

	c := newClient("c0ffee")
	c.base, c.http = srv.URL, srv.Client()
	c.now = func() time.Time { return time.Date(2026, 12, 20, 12, 0, 0, 0, time.UTC) }
	return site, c
}

func TestDownloadInput(t *testing.T) {
	site, c := newFakeSite(t, map[string]string{"/2024/day/3/input": "1 2\n3 4\n"})
	path := filepath.Join(t.TempDir(), "2024", "inputs", "3.txt")
	p := Puzzle{2024, 3}

	fetched, err := c.downloadInput(p, path)
	if err != nil || !fetched {
		t.Fatalf("downloadInput = %v, %v", fetched, err)
	}
	if b, _ := os.ReadFile(path); string(b) != "1 2\n3 4\n" {
		t.Errorf("saved %q", b)
	}

	// The second time, it's already there.
	if fetched, err = c.downloadInput(p, path); err != nil || fetched {
		t.Fatalf("second downloadInput = %v, %v", fetched, err)
	}
	if n := site.hits.Load(); n != 1 {
		t.Errorf("%d requests, want 1", n)
	}
}

func TestDownloadInputSession(t *testing.T) {
	site, c := newFakeSite(t, map[string]string{"/2024/day/3/input": "1 2\n"})
	path := filepath.Join(t.TempDir(), "3.txt")

	c.session = ""
	if _, err := c.downloadInput(Puzzle{2024, 3}, path); err != errNoSession {
		t.Errorf("without a session: %v", err)
	}
	if n := site.hits.Load(); n != 0 {
		t.Errorf("asked the site without a session")
	}

	c.session = "bad"
	if _, err := c.downloadInput(Puzzle{2024, 3}, path); err == nil || !strings.Contains(err.Error(), "session cookie") {
		t.Errorf("with a wrong session: %v", err)
	}
	if exists(path) {
		t.Error("saved the error page as the input")
	}
}

func TestGetBeforeUnlock(t *testing.T) {
	site, c := newFakeSite(t, nil)
	c.now = func() time.Time { return time.Date(2026, 12, 2, 4, 59, 0, 0, time.UTC) } // 23:59 EST
	_, err := c.get(Puzzle{2026, 2}, "/2026/day/2")
	if err == nil || err.Error() != "2026 day 2 unlocks in 1m0s" {
		t.Errorf("err = %v", err)
	}
	if site.hits.Load() != 0 {
		t.Error("asked the site for a puzzle that isn't out")
	}
}

func TestGetNotFound(t *testing.T) {
	_, c := newFakeSite(t, nil)
	if _, err := c.get(Puzzle{2025, 13}, "/2025/day/13"); err == nil || err.Error() != "2025 day 13: not found" {
		t.Errorf("err = %v", err)
	}
}
