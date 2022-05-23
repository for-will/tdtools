package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"testing"
	"text/tabwriter"
)

func TestParseExpr(t *testing.T) {
	expr, _ := parser.ParseExpr("\"123\"")
	ast.Print(nil, expr)

	expr, _ = parser.ParseExpr(`100 +		200*300`)
	ast.Print(nil, expr)
}

func Example_elastic() {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "a\tb\tc")
	fmt.Fprintln(w, "aa\tbb\tcc")
	fmt.Fprintln(w, "aaa\tbbbbbbbbbbbbbbbbb\t") // trailing tab
	fmt.Fprintln(w, "aaaa\tdddd\teeee")
	w.Flush()

	// output:
	// ....a|..b|c
	// ...aa|.bb|cc
	// ..aaa|
	// .aaaa|.dddd|eeee
}

func TestTableWriter(t *testing.T) {
	Example_elastic()
}

func TestPrintQ(t *testing.T) {
	fmt.Printf("say: %q world", "he\"llo")
}