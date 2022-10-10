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

func TestCap(t *testing.T) {
	var l1 []int32
	for i := 0; i < 10; i++ {
		l1 = append(l1, int32(i))
	}

	t.Log(len(l1), cap(l1))

	l2 := l1[5:]
	l2 = append(l2, 200)
	l1 = append(l1, 100)
	var l3 []int32
	l3 = append(l3, l1...)
	copy(l3, l2)

	//t.Log(cap(l2), len(l2))
	t.Log(l1)
	t.Log(l2)
	t.Log(l3)
}