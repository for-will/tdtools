package main

import (
	"fmt"
	"robot/client"
	"sync"
)

func main() {

	//client.UpdateDb()
	//fmt.Println(f())
	//
	newRobot()
	//Benchmark()
}

func newRobot() {
	r := &client.Robot{
		MsgHandler: client.DefaultMsgHandler,
		Account:    client.RobotAccount,
		Password:   client.RobotPassword,
	}

	r.Start()
	client.Log.Sync()
}

func Benchmark() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		r := &client.Robot{
			MsgHandler: client.DefaultMsgHandler,
			Account:    fmt.Sprintf("Test%d", i),
			Password:   "123456",
		}

		wg.Add(1)
		go func() {
			r.Start()
			wg.Done()
		}()
	}
	wg.Wait()
}

func f() (r int) {
	t := 5
	defer func() {
		r = r + t
	}()
	return 1
}
