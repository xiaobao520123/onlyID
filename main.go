package main

import (
	"fmt"
	"sync"
)

func main() {
	h, err := NewHost(1)
	if err != nil {
		fmt.Printf("create host failed, err: %v\n", err)
		return
	}
	wg := sync.WaitGroup{}
	ids := make([]ID, 0)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			ids = append(ids, h.Generate())
		}()
	}
	for i, id := range ids {
		fmt.Printf("id[%v]: %v\n", i, id.ToString())
		fmt.Printf("%d | %d | %d\n", id.Timestamp(), id.NodeID(), id.SeqID())

	}
}
