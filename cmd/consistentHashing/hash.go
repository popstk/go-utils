package main

import (
	"crypto/md5"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Node struct {
	Id      string
	Address string
}

type ConsistentHashing struct {
	mutex    sync.RWMutex
	nodes    map[int]Node
	replicas int
}

func NewConsistentHashing(nodes []Node, replicas int) *ConsistentHashing {
	ch := &ConsistentHashing{nodes: make(map[int]Node), replicas: replicas}
	for _, node := range nodes {
		ch.AddNode(node)
	}
	return ch
}

// 每个节点，生成replicas个虚拟节点
func (ch *ConsistentHashing) AddNode(node Node) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	for i := 0; i < ch.replicas; i++ {
		k := hash(node.Id + "_" + strconv.Itoa(i))
		ch.nodes[k] = node
	}
}

func (ch *ConsistentHashing) RemoveNode(node Node) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	for i := 0; i < ch.replicas; i++ {
		k := hash(node.Id + "_" + strconv.Itoa(i))
		delete(ch.nodes, k)
	}
}

func (ch *ConsistentHashing) GetNode(outerKey string) Node {
	key := hash(outerKey)
	nodeKey := ch.findNearestNodeKeyClockwise(key)
	return ch.nodes[nodeKey]
}

// 顺时针查找最近的节点，否则返回第一个节点
func (ch *ConsistentHashing) findNearestNodeKeyClockwise(key int) int {
	ch.mutex.RLock()
	sortKeys := sortKeys(ch.nodes)
	ch.mutex.RUnlock()
	for _, k := range sortKeys {
		if key <= k {
			return k
		}
	}
	return sortKeys[0]
}

func sortKeys(m map[int]Node) []int {
	var sortedKeys []int
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func hash(key string) int {
	md5Sum := md5.Sum([]byte(key))
	return int(crc32.ChecksumIEEE(md5Sum[:]))
}