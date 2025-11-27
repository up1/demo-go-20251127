package main

import (
	"demo"
	"fmt"
)

func main() {
	result := demo.SayHi()
	println(result)
	fmt.Println(result)

	for i := 0; i < 5; i++ {
		fmt.Printf("Count: %d\n", i)
	}
}
