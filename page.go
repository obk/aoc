package main

import (
	"html"
	"regexp"
	"strings"
)

var (
	titleRe = regexp.MustCompile(`<h2>--- Day \d+: (.+?) ---</h2>`)
	codeRe  = regexp.MustCompile(`(?s)<pre><code>(.*?)</code></pre>`)
	tagRe   = regexp.MustCompile(`<[^>]*>`)
)

// pageTitle returns the puzzle's name from its page, like "Hoof It", or "".
func pageTitle(page []byte) string {
	m := titleRe.FindSubmatch(page)
	if m == nil {
		return ""
	}
	return strings.Join(strings.Fields(html.UnescapeString(string(m[1]))), " ")
}

// pageExample returns the first code block on the puzzle page, which is
// usually the example input, or "" if there is none. Highlighting inside it,
// like <em>, is dropped.
func pageExample(page []byte) string {
	m := codeRe.FindSubmatch(page)
	if m == nil {
		return ""
	}
	return html.UnescapeString(tagRe.ReplaceAllString(string(m[1]), ""))
}
