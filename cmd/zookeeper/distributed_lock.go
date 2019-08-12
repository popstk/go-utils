package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"sync"
	"time"
)

const (
	lock = "/myLock"
	prefix = "lock-"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func DoSth(i, t int) {
	fmt.Println(i, " >> Get it - ", t)
	time.Sleep(1 * time.Second)
	fmt.Println(i,">> Release it - ", t)
}

func lockProcess(id int) {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		fmt.Println(id, ">", err)
		return
	}
	defer c.Close()

	data,err := c.Create(lock+"/"+prefix, []byte("hello"), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println(id, ">", err)
		return
	}
	pos := data[len(lock)+1:]

	for t :=1; ;t++ {
		children, _, err := c.Children(lock)
		Must(err)
		sort.Strings(children)

		// 最小的节点取得锁
		if children[0] == pos {
			DoSth(id, t)
			return
		}

		// 总是watch前面的锁
		for i, child := range children {
			if child == pos {
				exists, _, ch, err := c.ExistsW(lock+"/"+ children[i-1])
				if err != nil {
					fmt.Println(i, ">> ", err)
					continue
				}
				if !exists {
					fmt.Println(i, ">> not exists")
					continue
				}

				ev := <- ch
				if ev.Type == zk.EventNodeDeleted {
					break
				}
			}
		}
	}
}


func DistributedLock() {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		fmt.Println(">", err)
	}

	data, err := c.Create(lock, []byte("fuck"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println(">", err)
	} else {
		fmt.Println(data)
	}

	var wg sync.WaitGroup
	for i:=0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			lockProcess(i)
			wg.Done()
		}(i)
	}
	
	wg.Wait()
}

