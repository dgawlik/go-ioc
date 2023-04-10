package main

import (
	"fmt"

	goioc "github.com/dgawlik/go-ioc"
)

type SomeFunc func(int, string)

type SomeOtherFunc func() bool

func main() {

	err := goioc.Bind[SomeOtherFunc](func() string {
		return "Hello World"
	})

	fmt.Println(err)

	err = goioc.InjectBind[SomeFunc](func(fn SomeOtherFunc) func(string) {
		return func(s string) {
			fmt.Println("Hello World")
		}
	}, true)

	fmt.Println(err)
}
