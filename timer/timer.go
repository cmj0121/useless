// ~ Copyright (C) 2020-2020 cmj <cmj@cmj.tw>. All right reserved. ~
package main

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"log"
	"html/template"
)

func Help() {
	os.Stderr.WriteString("usage: timer [OPTION]\n")
	os.Stderr.WriteString("\n")
	os.Stderr.WriteString("option\n")
	os.Stderr.WriteString("  -h, --help           show this message\n")
	os.Stderr.WriteString("  -b, --bind ADDRESS   bind to specified address and port (default :8888)\n")
	os.Exit(1)
}

func html_index(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL)

	query := r.URL.Query()
	title := query.Get("m")
	if title == "" {
		title = "開啟網頁"
	}

	timestamp := query.Get("t")

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, struct {
		Title string
		Timestamp string
	} {
		Title: title,
		Timestamp: timestamp,
	})
}

func html_static(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL)

	switch r.URL.Path {
	case "/main.dart.js":
		file, _ := ioutil.ReadFile(r.URL.Path[1:])
		w.Write(file)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s", r))
			Help()
		}
	}()

	bind := ":8888"
	debug := false
	for idx := 1; idx < len(os.Args); idx ++ {
		switch args := os.Args[idx]; args {
		case "-d", "--debug":
			debug = true
		case "-h", "--help":
			Help()
		case "-b", "--bind":
			idx ++
			if idx >= len(os.Args) {
				err := fmt.Errorf("-b, --bind need ADDRESS")
				panic(err)
			}
			bind = os.Args[idx]
		default:
			Help()
		}
	}

	if debug {
		http.HandleFunc("/main.dart.js", html_static)
	}
	http.HandleFunc("/", html_index)
	log.Printf("run server on %s", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
