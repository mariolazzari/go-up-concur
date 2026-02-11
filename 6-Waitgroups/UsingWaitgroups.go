package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	var wg = sync.WaitGroup{}
	wg.Add(1)
	go doSomething(&wg)
	wg.Add(1)
	go doSomethingElse(&wg)
	wg.Wait()

	fmt.Println("\n\nI guess I'm done")
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func doSomething(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * 2)
	fmt.Println("\nI've done something")
}

func doSomethingElse(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * 2)
	fmt.Println("I've done something else")
}
