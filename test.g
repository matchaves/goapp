package main

import (
	"fmt"
)

func main() {
	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	//rand.Seed(time.Now().UnixNano())
	//var id int
	//id = rand.Intn(50)
	//fmt.Printf("Type: %T Value: %v\n", id, id)
	//fmt.Printf("Type: %T Value: %v\n", id, id)

	//var st2,st1 string
	var s, t string
	s = "fuck"
	t = fmt.Sprintf("There are %s reasons to code!", s)
	fmt.Println(t)

}
