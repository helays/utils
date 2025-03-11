package main

import (
	"fmt"
	"reflect"
)

func main() {
	run()
}

type t struct {
	val any
}

func run() (res *t) {
	ref := reflect.TypeOf(res)
	fmt.Println("ccc", "1", ref.String(), ref.Kind() == reflect.Ptr)
	run1(&res)
	return
}

func run1(v any) {
	ref := reflect.TypeOf(v)
	fmt.Println("vvv", ref.Name(), ref.String(), ref.Kind() == reflect.Ptr)
}
