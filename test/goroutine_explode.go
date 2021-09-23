package main

import (
	"strconv"
	"sync"
)

func goroutineExplode() {

	//defer func() {
	//	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	//}()

	ch := make(chan string, 100)

	//limit := make(chan struct{}, 100)
	wg := &sync.WaitGroup{}

	go func() {
		for range ch {

		}
		//buf := bytes.NewBuffer(make([]byte, 0, 65536))
		//for s := range ch {
		//	if buf.Len()+len(s) > buf.Cap() {
		//		fmt.Println(buf.Len(), buf.Cap())
		//		buf.Reset()
		//	}
		//	buf.WriteString(s)
		//}
	}()

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		//limit <- struct{}{}
		go func(idx int) {
			s := strconv.FormatInt(int64(idx), 10) + ", "
			ch <- s
			//<-time.After(time.Millisecond*time.Duration(i)%50 + 30*time.Second)
			wg.Done()
			//<-limit
		}(i)
	}

	wg.Wait()
	//<-time.After(time.Second * 2)
	/*	p := make([]runtime.MemProfileRecord, 64, 64)
		n, ok := runtime.MemProfile(p, true)
		fmt.Println(n, ok)

		var total int64
		for i := 0; i < n; i++ {
			total += p[i].AllocBytes
			//fmt.Println(p[i].AllocBytes)
		}
		fmt.Println(total / 1024)*/

	close(ch)
}
