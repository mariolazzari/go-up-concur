package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	wg              sync.WaitGroup
	widgetInventory int32 = 1000 //Package-level variable will avoid all the pointers
)

func main() {
	fmt.Println("Starting inventory count = ", widgetInventory)
	wg.Add(2)
	go makeSales()
	go newPurchases()
	wg.Wait()
	fmt.Println("Ending inventory count = ", widgetInventory)
}

func makeSales() { // 1000000 widgets sold
	for range 300000 {
		atomic.AddInt32(&widgetInventory, -100)
	}
	wg.Done()
}

func newPurchases() { // 1000000 widgets purchased
	for range 300000 {
		atomic.AddInt32(&widgetInventory, 100)
	}
	wg.Done()
}
