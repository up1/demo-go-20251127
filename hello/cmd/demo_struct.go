package main

import "fmt"

type User struct {
	Id   int
	Name string
	Age  int
}

type UserV2 struct {
	User
	Salary int
}

func (u UserV2) String() string {
	return fmt.Sprintf("UserV2(Id: %d, Name: %s, Age: %d, Salary: %d)", u.Id, u.Name, u.Age, u.Salary)
}

func (u User) doSth() string {
	u.Age += 1
	return fmt.Sprintf("User %s is %d years old.", u.Name, u.Age)
}

func main() {
	u2 := UserV2{
		User:   User{Id: 2, Name: "Bob", Age: 25},
		Salary: 50000,
	}

	fmt.Println(u2)
	fmt.Println(u2.Id)
	fmt.Println(u2.Name)
	fmt.Println(u2.Age)
	fmt.Println(u2.Salary)
	u2.doSth()
}
