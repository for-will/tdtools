package main

import (
	"fmt"
	"reflect"
	"robot/client"
)

func main() {
	a := new(struct{})
	b := new(struct{})

	println(a, b, a == b)
	fmt.Println(a, b, a == b)

	c := new(struct{})
	d := new(struct{})

	println(c, d, c == d)
	go func() {
		//_ = c
		//_ = d
	}()
	//fmt.Println(c, d, c == d)

	v := client.ByteOrder
	fmt.Println(123, reflect.TypeOf(v).PkgPath())

	//sync.Cond{}
}
