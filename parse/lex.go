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
	tokenComment        // '#'
	tokenImage          // '@'
	tokenText           // plain text
	tokenNewline        // '\n' or '.' followed by a '\n'
	tokenParagraphDelim // 2 * '\n'
	tokenSpace          // ' '
	tokenTab            // '\t'
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
	case len(t.val) > 20:
		return fmt.Sprintf("%.20q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
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
	for l.state = lexChar; l.state != nil; {
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
	l.tokens <- token{tok, l.start, l.input[l.start:l.pos], l.line}
	l.start = l.pos
}

// ignore ignores what we've read so far and sets the input pos to front.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// nextToken returns next token from the channel.
// This is only used by the parser.
func (l *lexer) nextToken() token {
	token := <-l.tokens
	l.lastPos = l.pos
	return token
}

// drain will empty the channel, so the lexer routine can finish.
// This is intended for the parser.
func (l *lexer) drain() {
	for _ = range l.tokens {
	}
}

// errorf returns a tokenized error with value containing the error message.
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

// lexChar is the default state.
func lexChar(l *lexer) stateFn {
	l.width = 0
	switch r := l.next(); {
	case r == eof:
		break
	case r == '\n':
		return lexNewline
	case r == '#':
		return lexComment
	case r == '@':
		return lexImage
	case r == ' ':
		l.emit(tokenSpace)
		return lexChar
	case r == '\t':
		l.emit(tokenTab)
		return lexChar
	case r == '.':
		return lexDot
	case r == '\\':
		l.ignore()
		return lexText
	default:
		// everything else is just text
		return lexText
	}
	// reached eof: send away whatever is left of the input.
	if l.pos > l.start {
		l.emit(tokenText)
	}
	l.emit(tokenEOF)

	// returning nil will end the state machine
	return nil
}

// lexText absorbs unrecognized characters up until, but not including, a newline
func lexText(l *lexer) stateFn {
	for !isEOL(l.peek()) {
		l.next()
	}
	l.emit(tokenText)
	return lexChar
}

// lexNewline returns either a tokenNewline or tokenParagraphDelim, based on number of newlines found.
// One newline is already consumed.
func lexNewline(l *lexer) stateFn {
	if l.accept("\n") {
		// twice! it is a paragraph. absorp more newlines if possible.
		l.acceptRun("\n")
		l.emit(tokenParagraphDelim)
	} else {
		l.emit(tokenNewline)
	}
	return lexChar
}

// lexComment returns a tokenComment with value of the comment stripped from # marker and optional whitespace between.
func lexComment(l *lexer) stateFn {
	// ignore the comment marker and optional whitespace
	l.acceptRun(" \t")
	l.ignore()

	for !isEOL(l.peek()) {
		l.next()
	}
	l.emit(tokenComment)

	l.acceptRun("\n")
	l.ignore()

	return lexChar
}

// lexImage reads a text value from the @ sign until '\n' and returns it as a tokenImage.
func lexImage(l *lexer) stateFn {
	// ignore the @
	l.ignore()

	for !isEOL(l.peek()) {
		l.next()
	}
	l.emit(tokenImage)

	return lexChar
}

// lexDot returns a tokenNewline if '.' is followed by a '\n'.
func lexDot(l *lexer) stateFn {
	if l.peek() == '\n' {
		l.next()
		l.emit(tokenNewline)
		return lexChar
	}

	return lexText
}

// helper functions

// isEOL returns true if rune is EOF or newline.
func isEOL(r rune) bool {
	return r == eof || r == '\n'
}
