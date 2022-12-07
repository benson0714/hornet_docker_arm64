package workerpool

import (
	"sync"
	"testing"
)

func TestSignalBarrier(t *testing.T) {
	pool := New(func(task Task) {
		println(task.Param(0).(int))
		task.Return(nil)
	}, WorkerCount(10), QueueSize(2000))
	pool.Start()

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)

		go func(i int) {
			<-pool.Submit(i)

			wg.Done()
		}(i)
	}

	pool.SubmitBarrier()

	for i := 0; i < 200; i++ {
		wg.Add(1)

		go func(i int) {
			<-pool.Submit(i)

			wg.Done()
		}(i)
	}

	wg.Wait()

}

func Benchmark(b *testing.B) {
	pool := New(func(task Task) {
		task.Return(task.Param(0))
	}, WorkerCount(10), QueueSize(2000))
	pool.Start()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)

		go func(i int) {
			<-pool.Submit(i)

			wg.Done()
		}(i)
	}

	wg.Wait()
}
