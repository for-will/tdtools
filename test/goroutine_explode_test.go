package main

import (
	"fmt"
	"reflect"
	"robot/GameMsg"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_goroutineExplode(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"run goroutineExplode"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goroutineExplode()
			printMem(t)
		})
	}
}

func Benchmark_goroutineExplore(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		goroutineExplode()

	}

}

func TestSelectWriteChannel(t *testing.T) {

	ch := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for true {
			select {
			case ch <- struct{}{}:
			default:
				runtime.Goexit()
			}
		}
	}()

	//close(ch)
	wg.Wait()
}

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

func TestIsClosed(t *testing.T) {
	ch := make(chan struct{}, 1)
	//ch <- struct{}{}
	t.Log(IsClosed(ch))
}

func IsClosed(ch chan struct{}) bool {

	select {
	case <-ch:
		return true
	default:

	}
	return false

	//atomic.CompareAndSwapUint32()
}

func TestSched(t *testing.T) {

	//debug.SetMaxThreads(0)
	runtime.GOMAXPROCS(1)

	var a int32 = 1

	go func() {
		for {
			//atomic.CompareAndSwapInt32()
			if atomic.LoadInt32(&a) == 100 {
				break
			}
			atomic.AddInt32(&a, 1)
			runtime.Gosched()
		}
	}()

	for {
		v := atomic.LoadInt32(&a)

		fmt.Println(v)
		if v == 100 {
			break
		}
		runtime.Gosched()
	}
}

func TestSelectSend(t *testing.T) {

	ch := make(chan int, 1)

	ch <- 1

	select {
	case ch <- 2:
	default:

	}

	t.Log(<-ch)

	select {
	case a := <-ch:
		t.Log(a)
	default:
		t.Log("default")
	}
	//t.Log(<-ch)
}

func TestReflectType(t *testing.T) {
	var l = []interface{}{
		(*GameMsg.AccountCheck)(nil),
		(*GameMsg.UnlockHeroTalentPageRs)(nil),
	}

	for _, v := range l {
		typ := reflect.TypeOf(v).Elem()
		t.Log(typ.Name())
	}
}
