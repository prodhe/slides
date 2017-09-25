package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenComment   // begin with '#' until '\n'; value is whatever is between
	tokenImage     // begin with '@' until '\n'; value is whatever is between
	tokenText      // plain text
	tokenNewline   // '\n'
	tokenParagraph // equals 2*newline
)

const eof = -1 // used for marking error in rune reading
type Pos int

type token struct {
	typ  tokenType
	pos  Pos
	val  string
	line int
}

func (t token) String() string {
	switch {
	case t.typ == tokenEOF:
		return "EOF"
	case t.typ == tokenError:
		return t.val
	}
	return fmt.Sprintf("%d: %s", t.typ, t.val)
}

// lexer holds the state of the scanner. Yes, this is Rob Pike-inspired deluxe.
type lexer struct {
	name    string  // for nice errors
	input   string  // input to lex
	state   stateFn // next lexing state
	pos     Pos     // current position in input
	start   Pos     // current token start pos
	width   Pos     // width of last rune read
	lastPos Pos     // position of last token read
	line    int     // count of '\n' seen +1 for pretty print
	tokens  chan token
}

type stateFn func(l *lexer) stateFn

// run loops through states until some state function returns nil.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek looks ahead without consuming the rune.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup unreads the latest rune.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit sends away a confirmed token.
func (l *lexer) emit(tok tokenType) {
	properline := l.line - strings.Count(l.input[l.start:l.pos], "\n")
	l.tokens <- token{tok, l.start, l.input[l.start:l.pos], properline}
	l.start = l.pos
}

// ignore ignores what we've read so far and sets the input pos to front.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) nextToken() token {
	token := <-l.tokens
	l.lastPos = l.pos
	return token
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

// create and return a new lexer
func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
		line:   1,
	}

	// start the lexing
	go l.run()

	return l
}

// state functions

func lexText(l *lexer) stateFn {
	l.width = 0
	switch r := l.next(); {
	case r == eof:
		break
	case r == '\n':
		l.backup()
		l.emit(tokenText)
		return lexNewline
	default:
		return lexText
	}
	// reached eof
	if l.pos > l.start {
		l.emit(tokenText)
	}
	l.emit(tokenEOF)
	return nil
}

func lexNewline(l *lexer) stateFn {
	_ = l.next() // we know this will be an \n
	if l.peek() == '\n' {
		// twice! it is a paragraph.
		for {
			if l.next() != '\n' {
				l.backup()
				break
			}
		}
		l.emit(tokenParagraph)
	} else {
		l.emit(tokenNewline)
	}
	return lexText
}
