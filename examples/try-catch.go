package main

import (
	"fmt"
	"github.com/ShindouMihou/go-bun/bun"
)

func main() {
	bun.Try(func() {
		file := bun.File("test")
		fmt.Println(file.Text())
	}).Catch(func(err any) {
		fmt.Println("failed to open test file: ", err)
	}).Run()
}
