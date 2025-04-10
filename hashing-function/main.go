package main

import "fmt"

func main() {
	input := "hello world"
	h := hashGenerator(input)
	fmt.Printf("Input: %s\nHash: %s\n", input, h)

	input2 := "hello world!"
	h2 := hashGenerator(input2)
	fmt.Printf("Input: %s\nHash: %s\n", input2, h2)
}
