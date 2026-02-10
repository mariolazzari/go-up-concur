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
