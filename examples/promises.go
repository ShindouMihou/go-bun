package main

import (
	"fmt"
	"github.com/ShindouMihou/go-bun/bun"
	"strings"
)

func main() {
	bun.NewPromise(func(res any) any {
		file := bun.File("test")
		return file.Text()
	}).Exceptionally(func(res any) any {
		fmt.Println("failed to read file", res)
		return nil
	}).Then(func(res any) any {
		return strings.ReplaceAll(res.(string), "s", "r")
	}).Then(func(res any) any {
		fmt.Println(res)
		return nil
	}).Await()
}
