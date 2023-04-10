[![Go Report Card](https://goreportcard.com/badge/github.com/dgawlik/go-ioc)](https://goreportcard.com/report/github.com/dgawlik/go-ioc)
[![Go Reference](https://pkg.go.dev/badge/github.com/dgawlik/go-ioc.svg)](https://pkg.go.dev/github.com/dgawlik/go-ioc)
[![codecov](https://codecov.io/gh/dgawlik/go-ioc/branch/main/graph/badge.svg?token=01D16O3TW8)](https://codecov.io/gh/dgawlik/go-ioc)
[![License: GPL v2](https://img.shields.io/badge/License-GPL_v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)
[![Tag](https://img.shields.io/github/v/tag/dgawlik/go-ioc)](https://github.com/dgawlik/go-ioc/releases/tag/v0.3.0)

# goioc - lightweight ioc container for go

Coming from Java background I was used to Spring to inject dependencies in large projects, so thought I will also use some container in Go. 
There are many other packages serving the same purpose, but this one has api (design) closest to my personal preferences. It also 
is so small that you can read the source code in coffee break.

### Features

* injection by *type definition*
* injection of any type
* dependencies computed on demand
* caching of computed dependencies
* properties attached per container
* minimal codebase


### Installation

To install this package, run the following command in your project directory.

```bash
go get github.com/dgawlik/go-ioc
```

Use it like so:

```go
import "github.com/dgawlik/go-ioc"
```

### Api

```go
func Resolve[T any](forceRebind bool) (T, error) 
```

Returns fully injected value bound to type T. On consecutive calls and forceRebind false, value from
cache is taken. forceRebind forces to create new injected values anyway and overwrite cache.


```go
func Bind[T any](value any) error
```

Associates value with type T. The value is taken as is from the function and put to cache.


```go
func BindInject[T any](value any) error
```

Expects curried function to be provided. Outer function's parameters are values to be injected. During
resolve this function is called to produce **proper** value.

```go
func SetProperty(key string, value any) 
```

Attaches some value to key container-wise. You can think of it as containers metadata. Then you can
inject `Properties` built-in type to recover values.

```go
func NewContainer() *Container 
```

Creates new container with fresh state and no bindings.

```go
func SetContainer(newC *Container)
```

Sets default container to this one.


### Design

The container in fact is currying on steroids. Had you used plain technique of curried functions you would 
have to pass dependencies all the way down the composed objects. In contrast, the container does this for you and alll you 
have to do is to only provide direct  dependencies. The injection is not limited to functions any type can be registered and 
injected (but remember to declare type alias!).

There are no scopes in the container familiar to Spring, but prototype and singleton scopes can easily be emulated. The container
looks up dependencies when you try to retrive value by type. If everything works the intermediate results are stored in cache. Consecutive
calls return same values. Not however, when you pass `forceRebind` true, then cache is omitted in computations and refreshed.


### Examples

Example on general usage

```go
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

	greeter, _ := goioc.InjectResolve[Greeter](func(isPrime IsPrime) func(name string, age int) {
		return func(name string, age int) {
			statement := "is not"
			if isPrime(age) {
				statement = "is"
			}

			fmt.Printf("Hello %s, your age %s prime.\n", name, statement)
		}
	}, false)

	greeter("Dominik", 33)
}

```


Short sample how properties can interact with injections.

```go
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
```

How to refresh dependencies

```go
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
```
