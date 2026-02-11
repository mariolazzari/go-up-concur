# Up and Running with Concurrency in Go

## Understanding concurrency

### Sync example

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func main() {
	start := time.Now()
	for i := 1; i < 50; i++ {

		client := &http.Client{}
		request, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			fmt.Print(err.Error())
		}
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Content-Type", "application/json")
		response, err := client.Do(request)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(response.Body)

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Print(err.Error())
		}
		var responseObject Response
		err = json.Unmarshal(bodyBytes, &responseObject)
		if err != nil {
			return
		}
		fmt.Println("\n", responseObject.Joke)

	}
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}
```

## First goroutine

### Not using

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	doSomething()
	doSomethingElse()

	time.Sleep(time.Second * 5)

	fmt.Println("\n\nI guess I'm done")
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func doSomething() {
	time.Sleep(time.Second * 2)
	fmt.Println("\nI've done something")
}

func doSomethingElse() {
	time.Sleep(time.Second * 2)
	fmt.Println("I've done something else")
}
```

### Using

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	go doSomething()
	go doSomethingElse()

	fmt.Println("\n\nI guess I'm done")
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func doSomething() {
	time.Sleep(time.Second * 2)
	fmt.Println("\nI've done something")
}

func doSomethingElse() {
	time.Sleep(time.Second * 2)
	fmt.Println("I've done something else")
}
```

### Sleep

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	go doSomething()
	go doSomethingElse()

	time.Sleep(time.Second * 5) // Don't do this in production.  VERY inefficient and unpredictable.

	fmt.Println("\n\nI guess I'm done")
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func doSomething() {
	time.Sleep(time.Second * 2)
	fmt.Println("\nI've done something")
}

func doSomethingElse() {
	time.Sleep(time.Second * 2)
	fmt.Println("I've done something else")
}
```

## Blocking code

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	doTheFirstThing()
	doTheSecondThing()
	fmt.Println("No more blocking.  I'm done")
}

func doTheFirstThing() {
	fmt.Println("FirstThing 'blocking' for 2 seconds")
	time.Sleep(time.Second * 2)
}

func doTheSecondThing() {
	fmt.Println("SecondThing 'blocking' for 3 seconds")
	time.Sleep(time.Second * 3)
}
```

### Using WaitGroup

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg = sync.WaitGroup{}
	wg.Add(2)
	start := time.Now()
	go doSomething(&wg)
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
```

## Using channels

### Syntax

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	// ch <- "hello from main" //This won't work since it blocks
	go sendMe(ch)
	for i := 1; i < 2; i++ { // This for loop just reads the channel as messages come in.
		fmt.Println(<-ch)
	}
}

func sendMe(ch chan<- string) {
	time.Sleep(time.Second * 2)
	ch <- "SendMe is done"
}
```

### Using channels

```go
package main

import (
	"fmt"
	"time"
)

var ch = make(chan string)

func main() {
	start := time.Now()
	go doSomething()
	go doSomethingElse()

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	fmt.Println("I guess I'm done")
	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func doSomething() {
	time.Sleep(time.Second * 2)
	fmt.Println("\nI've done something")
	ch <- "doSomething finished"
}

func doSomethingElse() {
	time.Sleep(time.Second * 2)
	fmt.Println("I've done something else")
	ch <- "doSomethingElse finished"
}
```

#### Buffered channels

From [Youtube](https://www.youtube.com/watch?v=LvgVSSpwND8&t=563s) "Concurrency in Go" tutorial by Jake Wright

```go
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
```

#### Select statement

Uses a select / case statement to monitor 2 channels and print whichever is active

```go
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go func() {
		for {
			time.Sleep(time.Second)
			c1 <- "Sending every 1 second"

		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 4)
			c2 <- "Sending every 4 sec"

		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 10)
			c3 <- "We're done"
		}
	}()

	for { // infinite for loop  This is the operator - listening for activity on all channels.
		// This is a clever way to get around the blocking nature when you try to read a channel
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg + " Something cool happened")
		case msg := <-c3:
			fmt.Println(msg)
			os.Exit(0)
		}
	}
}
```

## I/O vs CPU bound

### CPU bound

#### Sequential

```go
// This is NON concurrent code
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	runtime.GOMAXPROCS(16) // Extra processors don't help w sequential tasks

	start := time.Now()
	counta()
	countb()
	countc()
	countd()
	counte()
	countf()
	countg()
	counth()

	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}
func counta() {
	fmt.Println("AAAA is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("AAAA is done  ")

}
func countb() {
	fmt.Println("BBBB is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("BBBB is done")

}
func countc() {
	fmt.Println("CCCC is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("CCCC is done    ")

}
func countd() {
	fmt.Println("DDDD is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("DDDD is done     ")

}
func counte() {
	fmt.Println("EEEE is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("EEEE is done   ")

}
func countf() {
	fmt.Println("FFFF is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("FFFF is done     ")

}
func countg() {
	fmt.Println("GGGG is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("GGGG is done     ")

}
func counth() {
	fmt.Println("HHHH is starting  ")
	for i := 1; i < 10_000_000_000; i++ {
	}

	fmt.Println("HHHH is done     ")

}
```

#### I/O bound

```go
package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(16) // Extra processors don't help with sequential tasks
	fmt.Println(runtime.NumCPU())

	links := []string{
		"http://hashnode.com",
		"http://dev.to",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://medium.com",
		"http://github.com",
		"http://techcrunch.com",
		"http://techrepublic.com",
	}

	start := time.Now()

	for _, link := range links {
		checkLink(link)
	}

	elapsed := time.Since(start)
	fmt.Printf("Processes took %s", elapsed)
}

func checkLink(link string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "is not responding!")

		return
	}
	fmt.Println(link, "is LIVE!")
}
```

## Race conditions

### Mutex
