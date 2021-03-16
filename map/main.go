package main

import (
	"fmt"
	"strconv"
)

func main() {
	test2()
}

func test() {
	type name struct {
		first string
		last  string
	}
	type human struct {
		name name
		age  int8
	}
	humanMap := make(map[name]*human)
	a := &human{
		name: name{
			first: "taro",
			last:  "suzuki",
		},
		age: 20,
	}
	humanMap[name{
		first: "taro",
		last:  "suzuki",
	}] = a
	fmt.Println(humanMap[name{
		first: "taro",
		last:  "suzuki",
	}])
}

func test2() {
	a := map[int64]map[int64]string{1: {}, 2: {}, 3: {}}
	b := []map[int64]int64{{1: 1, 2: 2, 3: 3}, {4: 4, 5: 5, 6: 6}, {7: 7, 8: 8, 9: 9}}
	notSynchronized := make(map[int64]map[int64]string)
	for _, int64s := range b {
		for i, i2 := range int64s {
			if _, ok := a[i]; ok {
				notSynchronized[i][i2] = strconv.FormatInt(i2, 10)
			}
		}
	}
	fmt.Println(a)
}
