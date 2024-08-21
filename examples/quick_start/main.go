package hi

import "fmt"

func main() {
	s := MyStruct{
		Answer: 99,
	}
	err := Validate(s)
	fmt.Println(err) // answer is wrong
}
