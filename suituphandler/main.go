package main

import (
	"flag"
	"go/format"
	"io/ioutil"
	"market/suituphandler/internal"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()
	var dir string
	if len(args) != 0 {
		dir = filepath.Dir(args[0])
	} else {
		dir = "."
	}
	out := internal.GenHandlerWrap(dir)
	if formatted, err := format.Source([]byte(out)); err == nil {
		out = string(formatted)
	}

	err := ioutil.WriteFile("handler_wrap.go", []byte(out), 0644)
	if err != nil {
		panic(err)
	}
}
