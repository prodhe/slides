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

	http.Handle("/f/", http.StripPrefix("/f/", http.FileServer(http.Dir("./"))))
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
	result := HTML_TMPL

	result = strings.Replace(result, "{{style}}", STYLESHEET, -1)
	result = strings.Replace(result, "{{javascript}}", JAVASCRIPT, -1)

	result = strings.Replace(result, "{{data}}", data, -1)

	return result
}
