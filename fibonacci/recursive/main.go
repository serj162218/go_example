package main

import "fmt"

var arr []int

func main() {
	n := 100
	arr = make([]int, n+1)
	fmt.Println(fib(n))
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	if arr[n] != 0 {
		return arr[n]
	}
	arr[n] = fib(n-1) + fib(n-2)
	return arr[n]
}
