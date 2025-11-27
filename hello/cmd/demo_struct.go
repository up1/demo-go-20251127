package main

import "fmt"

type User struct {
	Id   int
	Name string
	Age  int
}

func main() {
	u1 := User{
		Id:   1,
		Name: "Alice",
		Age:  30,
	}

	fmt.Println(u1)
	fmt.Println(u1.Id)
	fmt.Println(u1.Name)
}
