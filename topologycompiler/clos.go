package topologycompiler

import "fmt"

func Compile(c chan string) {
	for i := 0; i <= 9; i++ {
		c <- fmt.Sprintf("%d", i)
	}

	close(c) // close channel when done with business logic
}
