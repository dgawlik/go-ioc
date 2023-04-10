package main

import (
	"fmt"
	"math"

	goioc "github.com/dgawlik/go-ioc"
)

type IsPrime func(num int) bool

type Greeter func(name string, age int)

func main() {

	goioc.Bind[IsPrime](func(num int) bool {
		if num < 2 {
			return false
		}
		sq_root := int(math.Sqrt(float64(num)))
		for i := 2; i <= sq_root; i++ {
			if num%i == 0 {
				return false
			}
		}
		return true
	})

	goioc.InjectBind[Greeter](func(isPrime IsPrime) func(name string, age int) {
		return func(name string, age int) {
			statement := "is not"
			if isPrime(age) {
				statement = "is"
			}

			fmt.Printf("Hello %s, your age %s prime.\n", name, statement)
		}
	}, false)

	greeter, _ := goioc.Resolve[Greeter](false)

	greeter("Dominik", 33)
}
