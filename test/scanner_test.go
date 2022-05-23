package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
	"testing"
	"text/tabwriter"
)

func TestScanner(t *testing.T) {
	var src = []byte(`println("hello,world")`)

	var fset = token.NewFileSet()
	var file = fset.AddFile("hello.go", fset.Base(), len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	//tabwriter.NewWriter(os.Stdout, 1).
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Fprintf(w, "\t%v\t%s\t%q\t\n", fset.Position(pos).Column, tok, lit)
	}
	w.Flush()
}
