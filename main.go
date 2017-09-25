package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type TokenType int

const (
	POPEN TokenType = iota
	PCLOSE
	LINE
)

type Token struct {
	Type  TokenType
	Value string
}

func main() {
	ch := make(chan Token)
	out := make(chan string)

	var f *os.File
	f, err := os.Open("slides.txt")
	if err != nil {
		f = os.Stdin
	}

	go parse(ch, out)

	err = scan(f, ch)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f.Close()

	data := <-out

	http.Handle("/", slides(toHtml(data)))
	err = http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func slides(data string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, data)
	})
}

func toHtml(data string) string {
	result := `<!doctype html>
<html>
<head>
<title>Slides</title>
{{style}}
</head>
<body>
{{data}}
</body>
</html>
`

	style := `<style type="text/css">
* { border: 0; margin: 0; padding: 0; }
body { background-color: #ffffea; }
p {
	width: 100vw;
	height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 24pt;
	font-family: monospace;
	line-height: 1.7;
	text-align: center;
}
</style>`

	result = strings.Replace(result, "{{style}}", style, -1)
	result = strings.Replace(result, "{{data}}", data, -1)

	return result
}

func scan(f *os.File, ch chan Token) error {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	ch <- Token{POPEN, ""}

	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "":
			ch <- Token{PCLOSE, ""}
			ch <- Token{POPEN, ""}
		default:
			ch <- Token{LINE, line}

		}
	}

	if scanner.Err() == nil {
		ch <- Token{PCLOSE, ""}
	}

	close(ch)

	return nil
}

func parse(ch chan Token, out chan string) {
	s := ""
	for token := range ch {
		switch token.Type {
		case POPEN:
			s += fmt.Sprintf("<p>\n")
		case PCLOSE:
			s += fmt.Sprintf("</p>\n")
		case LINE:
			s += fmt.Sprintf("\t%s<br>\n", token.Value)
		}
	}
	out <- s
}
