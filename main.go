package main

import (
	"fmt"
	"github.com/hababisha/ABT/task1"
	"github.com/hababisha/ABT/task2"
)

func main(){
	// fmt.Println("hello world")
	task1.Hello()
	// task2.FreqCount()
	fmt.Println(task2.Palindrome("hello"))
	fmt.Println(task2.Palindrome("omo"))

	fmt.Println(task2.FreqCount("hello hello, Hello test"))
	
}