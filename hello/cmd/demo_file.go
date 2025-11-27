package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("non_existent_file.txt")
	if err != nil {
		fmt.Println("Error", err)
	}
	defer f.Close()
	// TODO: Add your code here

}
