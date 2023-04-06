package main

import (
	"fmt"

	goioc "github.com/dgawlik/go-ioc"
)

type Fn func(elem int) bool

type Filter func(arr []int) []int

func main() {

	goioc.Bind[Fn](func(el int) bool {
		if el%2 == 0 {
			return true
		} else {
			return false
		}
	})

	goioc.BindInject[Filter](func(f Fn) func([]int) []int {

		return func(arr []int) []int {
			var newArr []int

			for _, e := range arr {
				if f(e) {
					newArr = append(newArr, e)
				}
			}

			return newArr
		}

	})

	arr := [10]int{1, 2, 3, 4, 5, 7, 8, 9}

	filter, _ := goioc.Resolve[Filter](true)

	fmt.Println(filter(arr[:]))

	goioc.Bind[Fn](func(el int) bool {
		if el%2 == 1 {
			return true
		} else {
			return false
		}
	})

	filter, _ = goioc.Resolve[Filter](true)

	fmt.Println(filter(arr[:]))

}
