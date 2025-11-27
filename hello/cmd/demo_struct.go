package main

import "fmt"

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) doSth() string {
	return fmt.Sprintf("User %s is %d years old.", u.Name, u.Age)
}

func main() {
	u1 := User{Id: 1, Name: "Alice", Age: 30}

	fmt.Println(u1)
	fmt.Println(u1.doSth())
}
