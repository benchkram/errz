# errz
Package errz provides simple error handling primitives to create more readable code and handle logging.

Inspired from [error handling by Dave Chaney](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) and based on [github.com/pkg/errors ](https://github.com/pkg/errors)

##### You should use errz when you want..
  * ..easily detect & find memory leaks like a nil pointer dereference etc.
  * ..to give context to your errors trough stack traces
  * ..to avoid the verbose but idiomatic error handling pattern
```go
err := someFunc()
if err != nil {
  log.Println(err)
  return
}
```

##### When should i **not** use errz?
If you write algorithms, modules or libraries with high performance in mind.

## Error handling



#### Return an error, add context, log and stop execution
A common use case is to return an error an stop further code execution in the currently executed function. To avoid the 3-4 lines we use a panic/recover pattern to stop code execution after the error occured and pass the error to the calling funtion
```go
func bar() (err error){
  defer errz.Recover(&err)

  err = foo()
  errz.Fatal(err) //will panic on error

  //some awesome pieces of code
  //which is not executed when a error occurred at foo()
  //..
}
```
This pattern will not just handle errors produced through errz.Fatal(). It will also handle memory corruptions, turns them into an error and add a stack trace to it. This can prevent you from extensive debugging sessions.

#### Logging
When you just want to log/report an error just do
```go
err := foo()
errz.Log(err)
```
This is useful during development and mostly used at the top level of an application

# License
MIT
