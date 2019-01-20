package main

import (
	"fmt"
	"math/rand"
)

func randNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	for {
		x := randNumber(20, 40)
		fmt.Println("range 20-40:", x)
	}
}
