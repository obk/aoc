package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	baseURL = "https://adventofcode.com"
	maxBody = 1 << 22 // inputs and puzzle pages are far smaller
)

var errNoSession = errors.New("no session cookie: run `aoc login` first")

// client talks to adventofcode.com, keeping to the site's rules for tools: it
// says who it is in the User-Agent, never fetches an input twice, and doesn't
// ask for a puzzle before it unlocks.
type client struct {
	base      string // baseURL, or a test server
	session   string // the session cookie; inputs need it, puzzle pages don't
	userAgent string
	http      *http.Client
	now       func() time.Time
}

func newClient(session string) *client {
	ua := "github.com/obk/aoc"
	if contact := os.Getenv("AOC_CONTACT"); contact != "" {
		ua += " by " + contact
	}
	return &client{base: baseURL, session: session, userAgent: ua, http: &http.Client{Timeout: 30 * time.Second}, now: time.Now}
}

// get fetches a path of puzzle p, like its page or its input.
func (c *client) get(p Puzzle, path string) ([]byte, error) {
	if wait := p.unlock().Sub(c.now()); wait > 0 {
		return nil, fmt.Errorf("%v unlocks in %v", p, wait.Round(time.Second))
	}
	req, err := http.NewRequest(http.MethodGet, c.base+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	if c.session != "" {
		req.AddCookie(&http.Cookie{Name: "session", Value: c.session})
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxBody {
		return nil, fmt.Errorf("GET %s: the response is over %d bytes", path, maxBody)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		return body, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("%v: not found", p)
	// The site answers 400 to a missing or expired session, and 500 to a malformed one.
	case strings.HasSuffix(path, "/input") && (resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusInternalServerError):
		return nil, errors.New("adventofcode.com didn't accept the session cookie; log in again and run `aoc login` with a fresh one")
	}
	return nil, fmt.Errorf("GET %s: %s", path, resp.Status)
}

// downloadInput saves p's input to path, unless it's already there: an input
// never changes, and the site asks tools not to fetch it again.
func (c *client) downloadInput(p Puzzle, path string) (fetched bool, err error) {
	if exists(path) {
		return false, nil
	}
	if c.session == "" {
		return false, errNoSession
	}
	body, err := c.get(p, p.path()+"/input")
	if err != nil {
		return false, err
	}
	return true, writeNew(path, body)
}
