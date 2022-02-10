package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) != 2 {
		return
	}

	url := os.Args[1]
	url = url[len("cygwin://") : len(url)-1]
	b, err := base64.StdEncoding.DecodeString(url)

	if err != nil {
		log.Fatal(err)
	}
	cmd := string(b)
	//log.Println(cmd)

	filePath, err := filepath.Abs("c:\\cygwin64\\bin\\mintty.exe")
	if err != nil {
		log.Fatal(err)
	}
	proc, err := os.StartProcess(filePath,
		[]string{
			filePath,
			"-t","TDGM",
			"-i", "/tdgm.ico",
			"-e", "/bin/bash", "-c", cmd,
		}, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	proc.Wait()
}
