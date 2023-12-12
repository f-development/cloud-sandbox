package main

import "fmt"

func main() {
	fmt.Println("Int")
	for i := 0; i < 16; i++ {
		fmt.Printf("%d\n", i)
	}
	fmt.Println("Hex")
	for i := 0; i < 16; i++ {
		fmt.Printf("%x\n", i)
	}
}
