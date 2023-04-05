package main

import (
	"fmt"
)

type Greeter func(name string)
type Adder func(a, b int) int
type Divider func(a, b float32) float32

type SomeStruct struct {
	x int
	g Greeter
	a Adder
}

// func(Greeter, Adder) func(x int)116
type Assembler func(x int)

// func(Assembler, Divider) func(x int16)
type Assembler2 func(x int16)

type Op func(x int) int

func main() {

	c := DefaultContainer

	Bind[Greeter](&c, func(name string) {
		fmt.Printf("Hello %s\n", name)
	})

	Bind[Adder](&c, func(a, b int) int {
		return a + b
	})

	Bind[Divider](&c, func(a, b float32) float32 {
		return a / b
	})

	BindInject[Assembler](&c, func(g Greeter, a Adder) func(x int) {
		return func(x int) {
			g("Dominik")
			fmt.Println(a(1, 2))
		}
	})

	BindInject[Assembler2](&c, func(ai Assembler, d Divider) func(x int16) {
		return func(x int16) {
			ai(0)
			fmt.Println(d(10, 5))
		}
	})

	BindInject[SomeStruct](&c, func(g Greeter, a Adder) SomeStruct {
		return SomeStruct{
			x: 2,
			g: g,
			a: a,
		}
	})

	BindInject[Op](&c, func(p Properties) func(int) int {
		mode, _ := p.GetString("mode")
		if mode == "quad" {
			return func(x int) int {
				return x * 4
			}
		} else {
			return func(x int) int {
				return x * 2
			}
		}
	})

	x, err := Resolve[Assembler2](c, true)

	if err != nil {
		panic(err)
	}

	x(0)

	y, err := Resolve[SomeStruct](c, true)

	if err != nil {
		panic(err)
	}

	y.g("Mark")

	SetProperty(&c, "mode", "dual")

	z, _ := Resolve[Op](c, true)

	fmt.Println(z(2))

	SetProperty(&c, "mode", "quad")

	z, _ = Resolve[Op](c, true)

	fmt.Println(z(2))
}
