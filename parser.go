package core

import (
	"errors"
	"strings"

	"github.com/sinloss/stk"
)

type quote rune
type state int

const (
	qut state = iota
	spl
	esc
)

// ArgParser the arguments parser type
type ArgParser struct {
	symbols  string
	escaper  rune
	boundary int
	lenient  bool
}

// NewArgParser create an arguments parser; any of the `splitors` inside a pair of identical `quotes`
// is treated as normal character; the `escaper` can escape a following rune; we can have more than
// one pair of `quotes` as well as `splitors`, yet only a pair of **identical** `quotes` are a integral
// scope of quote and all characters inside this scope will be treated as text except for the
// `escaper`; if the `lenient` is true, this parser will not report any error if any of the scopes is not
// finished( *meaning that something does not end properly* ) after the parsing process is over.
func NewArgParser(quotes string, splitors string, escaper rune, lenient bool) *ArgParser {
	return &ArgParser{
		symbols:  quotes + splitors,
		boundary: len([]rune(quotes)),
		escaper:  escaper,
		lenient:  lenient}
}

// GeneralArgParser the general purpose `lenient` arguments parser with `"` and `'` as its `quotes`, `\t` as its `splitors`,
// `\` as its escaper.
func GeneralArgParser() *ArgParser {
	return NewArgParser("\"'", " \t", '\\', true)
}

// check the the type of the given rune
func (ap *ArgParser) check(c rune) interface{} {
	if c == ap.escaper {
		return esc
	}

	i := strings.IndexRune(ap.symbols, c)
	if i == -1 {
		return nil
	} else if i < ap.boundary {
		return qut
	}
	return spl
}

// Parse parse the given text
func (ap *ArgParser) Parse(text string) ([]string, error) {
	var args []string
	scopes := stk.NewStack(false)
	var buf string

	// final a scope
	fin := func() {
		scopes.Pop()
		if buf != "" {
			args = append(args, buf)
			buf = ""
		}
	}

	// accumulate the given rune to the buffer
	acc := func(c rune) {
		buf += string(c)
	}

	// try to escape the given rune in the given scope
	escape := func(scope interface{}, c rune) bool {
		if scope == esc {
			acc(c)
			scopes.Pop()
			return true
		}
		return false
	}

	// the parsing begins
	scopes.Push(spl)
	for _, c := range []rune(text) {
		t := ap.check(c)
		scope := scopes.Peek()

		switch t {
		case esc:
			if !escape(scope, c) { // try to escape current escape mark
				scopes.Push(esc) // escape the next rune
			}
		case spl:
			switch scope.(type) {
			case quote: // in a scope of quote
				acc(c) // splitor is treated as a normal rune
			case state:
				if escape(scope, c) { // try to escape current splitor
					break
				}
				if scope == spl { // already in a scope of spl
					fin() // then this scope is done
				}
				scopes.Push(spl) // and the next spl scope must be well prepared
			}
		case qut:
			switch scp := scope.(type) {
			case quote: // already in a scope of quote
				if c == rune(scp) { // and c is exactly the same quote mark
					fin() // then this scope is done
					scopes.Push(spl)
				} else {
					acc(c) // otherwise it should be accumulated
				}
			case state:
				if !escape(scope, c) { // try to escape current quote mark
					if scope == spl {
						fin()
					}
					scopes.Push(quote(c)) // otherwise a new scope of quote should be started
				}
			}
		default:
			acc(c)
		}
	}

	// all the runes are consumed, yet there are still things need attendance
	if scopes.Size() != 0 {
		if ap.lenient {
			fin()
		} else {
			unfinished := scopes.Pop()
			if scopes.Size() == 0 && unfinished == spl {
				fin()
			} else if _, ok := unfinished.(quote); ok {
				return nil, errors.New("arg parser: unfinished quote for the input")
			} else if unfinished == esc {
				return nil, errors.New("arg parser: unfinished escape for the input")
			}
		}
	}

	return args, nil
}
