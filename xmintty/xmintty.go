package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	filePath, err := filepath.Abs("c:\\cygwin64\\bin\\mintty.exe")
	if err != nil {
		log.Fatal(err)
	}

	c := exec.Command(filePath,
		"-t", "TDGM",
		"-i", "/tdgm.ico",
		"/bin/sh",
		//"-lc", "exec zsh -lc 'date;pwd;ps;who;read'",
		"-lc", fmt.Sprintf("exec bash -lc '%s'", cmd),
	)
	c.Run()
}

func StartProcess(processPath string, cmd string) {
	proc, err := os.StartProcess(processPath,
		[]string{
			processPath,
			"-t", "TDGM",
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
