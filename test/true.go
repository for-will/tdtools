package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	//a := new(struct{})
	//b := new(struct{})
	//
	//println(a, b, a == b)
	//fmt.Println(a)
	//
	//c := new(struct{})
	//d := new(struct{})
	//
	//println(c, d, c == d)
	//defer func() {
	//	var m runtime.MemStats
	//	runtime.ReadMemStats(&m)
	//	fmt.Printf("%+v", m.HeapAlloc)
	//}()
	//var a int32 = 123
	//println(&a)
	//NewInt32(a)
	//dumpHeapAlloc()
	//dumpHeapAlloc()
	//println(NewInt32(a))
	//dumpHeapAlloc()
	//println(NewInt32(a))
	//rsp := &GameMsg.SyncPlayerBase{
	//	Lv: NewInt32(100),
	//}
	//println(rsp)
	//dumpHeapAlloc()

	//println(&a)
	//println(ShowMemEscape())

	//var p Person
	//var p2 = p
	//fmt.Println(p2)
	//var arr = [...]int{1, 2, 3, 4}
	//var s = "hello\tworld\n"
	//fmt.Printf("%q", s)
	//
	//var cd = GameMsg.TASK_CONDITION_COST_GOLD
	//fmt.Println(cd)
	//cd += 1

	//syscall.Exec()
	filePath, err := filepath.Abs("c:\\cygwin64\\bin\\mintty.exe")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("filePath: %s", filePath)
	proc, err := os.StartProcess(filePath, []string{filePath,
		"-e",
		"/bin/bash",
		"-c",
		"/bin/bash -c \"$(curl -sSL http://172.16.1.220:8081/script/nv1_install.sh)\"",
		//"-i",
		//"\"$(curl -sSL http://172.16.1.220:8081/script/nv1_install.sh)\"",
	},
		&os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	proc.Release()

}

func ShowMemEscape() interface{} {
	var a int32 = 123
	println(&a)
	NewInt32(a)
	dumpHeapAlloc()
	dumpHeapAlloc()
	println(NewInt32(a))
	dumpHeapAlloc()
	println(NewInt32(a))
	//rsp := &GameMsg.SyncPlayerBase{
	//	Lv: NewInt32(100),
	//}
	//println(rsp)
	dumpHeapAlloc()
	p := NewArr()
	dumpHeapAlloc()
	NewArr()
	dumpHeapAlloc()
	return p
}

func NewInt32(v int32) *int32 {
	return &v
}

func NewInt8(v int8) *int8 {
	return &v
}

func NewArr() interface{} {
	return new([128]int8)
}

var m runtime.MemStats
var lastHeapAlloc int64 = 0

func dumpHeapAlloc() {
	runtime.ReadMemStats(&m)
	println(m.HeapAlloc)
}

type DoNotCopy struct{}

func (*DoNotCopy) Lock()   {}
func (*DoNotCopy) Unlock() {}

type Person struct {
	DoNotCopy
	Name string
}

//
//type DoNotCompare [0]func()
//
//type NoUnkeyedLiterals struct{}
