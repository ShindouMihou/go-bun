package main

import (
	"fmt"
	"github.com/ShindouMihou/go-bun/bun"
)

func main() {
	bun.Try(func() {
		file := bun.File("examples/test-overwrite")
		file.Overwrite("hello")

		file = bun.File("examples/test-write")
		file.Write("hello\n")

		fmt.Println("test-write: ", file.Text())

		type Hello struct {
			World string `json:"world"`
		}

		file = bun.File("examples/test-write.json")
		file.Write(Hello{World: "test"})

		var hello Hello
		file.Json(&hello)

		fmt.Println("json: ", hello.World)
	}).Catch(func(err any) {
		fmt.Println("failed to perform write operations: ", err)
	}).Run()
}
