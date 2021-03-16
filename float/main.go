package main

import (
	"fmt"
	"math"
)

func main() {
	b := 1.2345678903
	fmt.Println(b)
	a := math.Pow10(10)
	bil := int64(b * a)
	fmt.Println(bil)
	converted := float64(bil) / math.Pow10(10)
	fmt.Println(converted)
}
