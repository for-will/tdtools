package client

import (
	"log"
	"strconv"
	"strings"
	"testing"
)

func MessageDump() {
	log.Printf("%s [%s|%s] ", GREEN("<"), YELLOW(strconv.Itoa(32)), GREEN("S2C_SeasonRs"))
}

func GREEN(s string) string {
	builder := strings.Builder{}
	builder.WriteString("\u001B[37m")
	builder.WriteString(s)
	builder.WriteString("\u001B[0m")
	return builder.String()
}

func RGREEN(s string) string {
	builder := strings.Builder{}
	builder.WriteString("\u001B[7;32m")
	builder.WriteString(s)
	builder.WriteString("\u001B[0m")
	return builder.String()
}

func YELLOW(s string)string  {
	builder := strings.Builder{}
	builder.WriteString("\u001B[33m")
	builder.WriteString(s)
	builder.WriteString("\u001B[0m")
	return builder.String()
}

func RYELLOW(s string)string  {
	builder := strings.Builder{}
	builder.WriteString("\u001B[7;33m")
	builder.WriteString(s)
	builder.WriteString("\u001B[0m")
	return builder.String()
}


func TestLogRcvMsg(t *testing.T) {
	MessageDump()
}
