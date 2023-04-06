package main

import (
	"fmt"

	goioc "github.com/dgawlik/go-ioc"
)

type Operation func(x int) int

func double(x int) int {
	return x * 2
}

func quad(x int) int {
	return x * 4
}

func main() {

	goioc.BindInject[Operation](func(props goioc.Properties) func(x int) int {
		v, _ := props.String("mode")
		if v == "double" {
			return double
		} else {
			return quad
		}
	})

	goioc.SetProperty("mode", "double")

	op, _ := goioc.Resolve[Operation](true)

	fmt.Printf("Operation double: %d -> %d\n", 2, op(2))

	goioc.SetProperty("mode", "quad")

	op, _ = goioc.Resolve[Operation](true)

	fmt.Printf("Operation quad: %d -> %d\n", 2, op(2))

}
