package parse

import "fmt"

type Parser struct {
	l *lexer
}

func NewParser(name, input string) *Parser {
	p := &Parser{l: lex(name, input)}
	return p
}

func (p *Parser) Parse() (string, error) {
	str := "<p>\n"
	for {
		tok := p.l.nextToken()
		//fmt.Printf("%s\n", tok)
		switch tok.typ {
		case tokenEOF:
			str += "</p>\n"
			return str, nil
		case tokenError:
			return "", fmt.Errorf("%s:%d: %s",
				p.l.name,
				tok.line,
				tok.val,
			)
		case tokenParagraphDelim:
			str += fmt.Sprintf("\n</p>\n<p>\n")
		case tokenNewline:
			str += fmt.Sprintf("<br>\n")
		case tokenText:
			str += tok.val
		case tokenComment:
			str += fmt.Sprintf("<!-- %s //-->\n", tok.val)
		}
	}

	str += "</p>"

	return str, nil
}