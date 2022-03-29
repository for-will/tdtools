package main

import (
	"io/ioutil"
	"log"
	"market/suituphandler/internal"
)

func main() {

	err := ioutil.WriteFile("robot/client/msgid_map.go", internal.GenMsgIdMap(), 0664)
	if err != nil {
		log.Fatalf("write: robot/client/msgid_map.go error: %v", err)
	}

}
