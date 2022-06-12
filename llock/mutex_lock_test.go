package llock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	// 逻辑中使用的某个变量
	count int
	// 与变量对应的使用互斥锁
	countGuard sync.Mutex
)

func Increment(useLock bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if useLock {
		countGuard.Lock()
		defer countGuard.Unlock()
	}
	count++
}

func TestMutex(t *testing.T) {
	// 并发去更新
	var wg sync.WaitGroup //定义一个同步等待的组
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go Increment(false, &wg)
	}
	wg.Wait()
	// 多运行几次，这里的结果不会是100
	fmt.Printf("没有使用锁的时候，并发更新的结果%d\n", count)

	count = 0
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go Increment(true, &wg)
	}
	wg.Wait()
	fmt.Printf("使用锁的时候，并发更新的结果%d\n", count)
}

func TestRWMutex(t *testing.T) {
	var rwMutex sync.RWMutex

	var wg sync.WaitGroup //定义一个同步等待的组
	wg.Add(1)
	// 协程1加读锁
	go func() {
		rwMutex.RLock()
		fmt.Println("协程1加读锁")
		defer rwMutex.RUnlock()
		time.Sleep(5 * time.Second)
		wg.Done()
	}()

	// 等1下，确保协程1已经加锁
	time.Sleep(1 * time.Second)
	wg.Add(1)
	// 协程2加读锁
	go func() {
		if rwMutex.TryRLock() {
			fmt.Println("协程2加读锁成功")
			defer rwMutex.RUnlock()
			time.Sleep(5 * time.Second)
		} else {
			t.Error("协程2加读锁失败")
		}
		wg.Done()
	}()

	// 等1下，确保协程2已经加锁
	// 协程3加写锁
	time.Sleep(1 * time.Second)
	wg.Add(1)
	go func() {
		if rwMutex.TryLock() {
			t.Error("协程3加写锁成功")
			defer rwMutex.Unlock()
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("协程3加写锁失败")
		}
		wg.Done()
	}()
	wg.Wait()
}
