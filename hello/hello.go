package demo

// Public function
func SayHi() string {
	r, err := doSth()
	if err != nil {
		return "Error"
	}
	return r
}

// Private function
// fff
func doSth() (result string, err error) {
	return "Hello Go 2025", nil
}
