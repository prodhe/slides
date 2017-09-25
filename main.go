package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/prodhe/slides/parse"
)

func main() {
	flag.Parse()

	var f *os.File
	var err error
	var path string

	if flag.NArg() < 1 || flag.Arg(0) == "-" {
		f = os.Stdin
	} else {
		path = flag.Arg(0)
		f, err = os.Open(path)
		if err != nil {
			fmt.Printf("%s: %v\n", path, err)
			os.Exit(1)
		}
	}
	reader := bufio.NewReader(f)
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("ReadAll error: %v\n", err)
		os.Exit(1)
	}
	f.Close()

	p := parse.NewParser(path, string(input))

	data, err := p.Parse()
	if err != nil {
		fmt.Printf("parse error: %s", err)
		os.Exit(1)
	}

	http.Handle("/", slides(toHtml(data)))
	fmt.Println("Slides are available at http://localhost:3001/")
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
* { border: 0; margin: 0; padding: 0; box-sizing: border-box; }
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
	text-align: left;
	padding: 1em;
}
</style>`

	result = strings.Replace(result, "{{style}}", style, -1)
	result = strings.Replace(result, "{{data}}", data, -1)

	return result
}
