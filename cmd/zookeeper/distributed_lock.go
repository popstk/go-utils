package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

const lock = "/myLock"

func lockProcess(i int) {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		fmt.Println(i, ">", err)
		return
	}
	
	data,err := c.CreateProtectedEphemeralSequential(lock, []byte("hello"), zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println(i, ">", err)
		return
	}

	fmt.Println(data)
	time.Sleep(5* time.Second)
	c.Close()
}


func DistributedLock() {
	var wg sync.WaitGroup
	for i:=0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			lockProcess(i)
			wg.Done()
		}(i)
	}
	
	wg.Wait()
}

