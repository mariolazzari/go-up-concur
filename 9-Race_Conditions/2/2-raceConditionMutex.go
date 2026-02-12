package main

import (
	"fmt"
	"sync"
)

var (
	wg              sync.WaitGroup
	mutex                 = sync.Mutex{}
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
		mutex.Lock()
		widgetInventory -= 100
		mutex.Unlock()
	}
	wg.Done()
}

func newPurchases() { // 1000000 widgets purchased
	for range 300000 {
		mutex.Lock()
		widgetInventory += 100
		mutex.Unlock()
	}
	wg.Done()
}
