package main

import (
	"fmt"
)

func main() {

	c := make(chan string, 3) // channel doesn't block until full ("buffered" channel)
	c <- "Hello "
	c <- "Earth "
	c <- "from Mars"
	//c <- "from Venus"

	msg := <-c
	fmt.Print(msg)

	msg = <-c // Notice we used = NOT := because msg is already declared
	fmt.Print(msg)

	msg = <-c // Notice we used = NOT := because msg is already declared
	fmt.Println(msg)

}
