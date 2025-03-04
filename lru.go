package main

import (
	"fmt"
	"sync"
)

type LruNode struct {
	Id    interface{}
	Value interface{}
}

type Lru struct {
	d     []*LruNode          // LRU数据队列
	m     map[interface{}]int // LRU数据邻接表
	c     int                 // LRU数量上限
	mutex *sync.RWMutex       // 读写锁
}

func newLru(cap int) *Lru {
	return &Lru{
		d:     make([]*LruNode, 0, 1024),
		m:     make(map[interface{}]int),
		mutex: &sync.RWMutex{},
		c:     cap,
	}
}

func (l *Lru) get(id string) interface{} {
	index, ok := l.m[id]	
	if !ok {
		fmt.Printf("no Value in lru, id: %s \n", id)
        return nil 
    }
	// 当前访问的节点
	currentNode := l.d[index]

	// 将被访问元素前移一位，被访问元素移到队尾
    length := len(l.d)
    for i := index; i < length - 1; i++ {
        l.d[i] = l.d[i+1]
        l.m[l.d[i].Id] = i
    }
    l.d[length-1] = currentNode
    l.m[id] = length-1
    
    return l.d[index].Value  
}

func (l *Lru) put(node *LruNode) {
	defer l.mutex.Unlock()
	l.mutex.Lock()
	if _, index := l.m[node.Id]; !index {
		l.addNode(node)
	}
}

func (l *Lru) addNode(node *LruNode) {
    id := node.Id 
	l.d = append(l.d, node)
	length := len(l.d)
	l.m[id] = length - 1
	if length > l.c {
		topNode := l.d[0]
		delete(l.m, topNode.Id)
		l.d = l.d[1:length]
		for i, n := range l.d {
			l.m[n.Id] = i
		}
	}
}

func (l *Lru) show() {
	defer l.mutex.Unlock()
	l.mutex.Lock()
	fmt.Println("show data:")
	for _, value := range l.d {
		fmt.Printf("%v ", value)
    }
	fmt.Println()
	fmt.Println("show map:")
	for key, value := range l.m {
        fmt.Printf("%v : %v ", key, value)
    }
	fmt.Println()
}

func main() {
	lru := newLru(5)
	lru.put(&LruNode{ Id: "a", Value: "a"})
	lru.put(&LruNode{ Id: "b", Value: "b"})
	lru.put(&LruNode{ Id: "c", Value: "c"})
	lru.put(&LruNode{ Id: "d", Value: "d"})
	lru.put(&LruNode{ Id: "e", Value: "e"})
    lru.show()
	lru.get("b")
    lru.show()
	lru.put(&LruNode{ Id: "f", Value: "f"})
	lru.show()
}
