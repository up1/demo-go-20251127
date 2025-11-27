package main

var i int = 0

func main() {
	m := make(map[string]int)

	m["Answer"] = 42
	v, found := m["Answer2"]
}
