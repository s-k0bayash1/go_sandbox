package main

import "fmt"
import "golang.org/x/exp/constraints"
import _ "github.com/s-k0bayash1/go_sandbox/generics/approximation_element"

func main() {
	a := "hello golang"
	fmt.Println(genSlice(a))
	b := 1
	fmt.Println(genSlice(b))
	fmt.Println(add(1, 2))
	fmt.Println(add(1.2, 2.4))
}

func genSlice[T any](v T) []T {
	res := make([]T, 1)
	res[0] = v
	return res
}

func add[T Number](a T, b T) T {
	return a + b
}

type Number interface {
	constraints.Integer | constraints.Float
}
