# aoc

A command-line helper for solving [Advent of Code](https://adventofcode.com) in Go. It sets up each day from a template, downloads your input once, and runs and times your solution. It's a single binary with no dependencies.

```
$ aoc new 2025 3
created 2025/3.go
saved the first code block on the page to 2025/inputs/3.example.txt; check that it's the example
downloaded the input to 2025/inputs/3.txt

$ aoc run 3 -e
Part 1: 0
Part 2: 0
2025 day 3 ran in 2.104ms
```

It expects one file per day, with the inputs next to them, which is the layout of [obk/advent-of-code](https://github.com/obk/advent-of-code):

```text
2025/
├── 3.go
└── inputs/
    ├── 3.txt
    └── 3.example.txt
```

Each day is a standalone `package main`, so `go run 3.go` works without aoc too.

## Install

```sh
go install github.com/obk/aoc@latest
```

## Usage

```sh
aoc login                   # save your session cookie, which downloading inputs needs
aoc new                     # in December: set up today's puzzle
aoc new 2024 7              # or any other one
aoc input 2024 7            # download just the input
aoc run 2024 7              # build and run it on the input, and time it
aoc run 7 -e                # on the example
aoc run 7 -i other.txt      # on another file
```

Run it in the repository or inside a year folder. In a year folder, a day alone means that year's day; elsewhere it means the latest event's day.

**`aoc new`** creates `YYYY/D.go` from the template, with the puzzle's title in its header comment. It saves the first code block on the puzzle page as `YYYY/inputs/D.example.txt`. That block is usually the example, but check it. Then it downloads the input. Files that already exist are never overwritten, so it's safe to run again.

**`aoc run`** builds the solution, runs it in its year folder (where its `inputs/D.txt` path points), and prints how long it ran on stderr. The build isn't timed, and neither is the virus scan that Windows runs before a freshly built program starts. With `-e` or `-i`, the solution gets the file's path as its first argument.

**Your session cookie** is the value of the `session` cookie on adventofcode.com: find it in your browser's developer tools, under Storage or Application. `aoc login` saves it in your config folder (`~/.config/aoc/session`, or `%AppData%\aoc\session` on Windows). You can also set `AOC_SESSION` instead.

## The template

The default template has a bounded input reader, a `dayN(path string) error` function that prints both parts, and a `main` that takes another input path as an optional argument. To use your own, put it in `.aoc/day.go.tmpl` at the root of your repository. It's a Go [text/template](https://pkg.go.dev/text/template) with `{{.Year}}`, `{{.Day}}` and `{{.Title}}`, and the result is formatted like gofmt.

## Being nice to adventofcode.com

The site [asks tools](https://www.reddit.com/r/adventofcode/wiki/faqs/automation) to cache inputs, identify themselves and not make needless requests. aoc:

- downloads an input only once, and never again while the file exists
- fetches the puzzle page only when `aoc new` still needs the title or the example
- doesn't ask for a puzzle before it unlocks (midnight EST)
- sends `User-Agent: github.com/obk/aoc`, followed by `by <contact>` if you set `AOC_CONTACT` to your email or repository, which the site asks for

Don't commit your inputs or the puzzle text: add `*/inputs/` to your `.gitignore`.

## Development

```sh
go test ./...
```

The tests use a fake adventofcode.com and never touch the real one. `TestNewAndRun` sets up a day from the default template, then builds and runs it.

## License

[MIT](LICENSE)
