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
			str += toHtml(tok.val)
		case tokenComment:
			str += fmt.Sprintf("<!-- %s //-->\n", tok.val)
		case tokenImage:
			str += fmt.Sprintf("<img src=\"/f/%s\">", tok.val)
		}
	}

	str += "</p>"

	return str, nil
}

func toHtml(s string) string {
	var result string
	s = html.EscapeString(s)
Loop:
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ':
			result += "&nbsp;"
		case '\t':
			result += "&nbsp;&nbsp;&nbsp;&nbsp;"
		default:
			result += s[i:]
			break Loop
		}
	}
	return result
}
