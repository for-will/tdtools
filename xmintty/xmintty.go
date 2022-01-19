package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
)

func main() {

	//log.Printf("filePath: %s", filePath)
	//s := base64.StdEncoding.EncodeToString([]byte("/bin/bash -c \"$(curl -sSL http://172.16.1.220:8081/script/trunk_install.sh)\""))
	//log.Printf("cmd: %s", s)
	//log.Printf("%+v", os.Args)

	if len(os.Args) != 2 {
		return
	}

	url := os.Args[1]
	//url := "cygwin://L2Jpbi9iYXNoIC1jICIkKGN1cmwgLXNTTCBodHRwOi8vMTcyLjE2LjEuMjIwOjgwODEvc2NyaXB0L252MV9pbnN0YWxsLnNoKSI=/"
	//strings.Index(url)
	url = url[9 : len(url)-1]
	b, err := base64.StdEncoding.DecodeString(url)

	if err != nil {
		log.Fatal(err)
	}
	cmd := string(b)
	log.Println(cmd)

	//<-time.After(10 * time.Second)

	filePath, err := filepath.Abs("c:\\cygwin64\\bin\\mintty.exe")
	if err != nil {
		log.Fatal(err)
	}
	proc, err := os.StartProcess(filePath,
		[]string{
			filePath,
			"-e",
			"/bin/bash",
			"-c",
			cmd,
		}, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	proc.Wait()
}
