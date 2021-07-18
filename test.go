package main

import (
	"fmt"
	"math/rand"
	"time"
)



func main() {
	//s1 := rand.NewSource(time.Now().UnixNano())
    //r1 := rand.New(s1)
	rand.Seed(time.Now().UnixNano())
	var id int
	id = rand.Intn(50)
	
	fmt.Printf("Type: %T Value: %v\n", id, id)
	fmt.Printf("Type: %T Value: %v\n", id, id)

}