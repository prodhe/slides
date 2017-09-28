package parse

import (
	"fmt"
	"html"
)

type Parser struct {
	l *lexer
}

func NewParser(name, input string) *Parser {
	p := &Parser{l: lex(name, input)}
	return p
}

func (p *Parser) Parse() (string, error) {
	str := "<section><div>\n"
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
			str += fmt.Sprintf("\n</div></section>\n\n<section><div>\n")
		case tokenNewline:
			str += fmt.Sprintf("<br>\n")
		case tokenText:
			str += fmt.Sprintf("<span>%s</span>", html.EscapeString(tok.val))
		case tokenSpace:
			str += "<span>&nbsp;</span>"
		case tokenTab:
			str += "<span>&nbsp;&nbsp;&nbsp;&nbsp;</span>"
		case tokenComment:
			str += fmt.Sprintf("<!-- %s //-->\n", tok.val)
		case tokenImage:
			str += fmt.Sprintf("<img src=\"/f/%s\">", tok.val)
		}
	}

	str += "</div></section>\n"

	p.l.drain()

	return str, nil
}
