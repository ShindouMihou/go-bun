# go-bun

> **Warning!**
> 
> go-bun is intended to be a little fun joke. do not ever use it in production, although you could try 
> copying the [`files.go`](bun/files.go) which is really clean imo, but just change the panic handling to use 
> the traditional go error handling.

some little fun experimentation and remake of random stuff from other languages. initially, this was to remake the cool  
interfaces that bun.sh had for files, but then expanded to figuring out how to remake stuff from other languages.

what does go-bun have:
- [`promises`](examples/promises.go): literally javascript promises, but dumbed down and made into a linked-list structure.
- [`try-catch`](examples/try-catch.go): literally java and javscript try-catch, but dumbed down.
- [`files`](examples/files.go): literally bun.sh's files interface, dumbed for golang.

## error-handling.

in order to imitate the simple look of bun.sh's interfaces, we also need to let go of golang's error-handling, and that's 
where go-bun's try-catch system comes into play.

```go
func main() {
	bun.Try(func() {
		file := bun.File("examples/test-overwrite")
		file.Overwrite("hello")
	}).Catch(func(err any) {
		fmt.Println("failed to perform write operations: ", err)
	}).Run()
}
```

literally... try-catch!

## promises

we can't anything javascript without promises, right? hell yes, we have promises and it's built with [craziness](bun/future.go), it's 
like traditional promises, but crazy:

```go
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
```