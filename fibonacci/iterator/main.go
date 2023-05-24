package main

import "fmt"

func main() {
	n := 100
	if n < 2 {
		fmt.Println(1)
		return
	}
	a := 1
	b := 1
	sum := 0
	for i := 2; i <= n; i++ {
		sum = a + b
		b = sum - b
		a = sum
	}
	fmt.Println(sum)
}
