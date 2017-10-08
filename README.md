# goerr
Error handling primitives to create more readable code and handle logging.
Inspire from [error handling by Dave Chaney](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) and based on [github.com/pkg/errors ](https://github.com/pkg/errors)

##### When should i use goerr
You want..
  * ..easily detect & find memory leaks like a nil pointer dereference etc.
  * ..to give context to your errors trough stack traces
  * ..to avoid the verbose standard error handling pattern (SEHP)
```go
    err := someFunc()
    if err != nil {
      log.Println(err)
      return
    }
```

##### When should i **not** use goerr
If you write algorithms, modules or libraries with performance in mind.

## Error handling



#### Return an error, add context, log and stop execution
A common usecase is to return an error an stop further code execution in the currently executed function. To avoid the 3-4 line SEHP we use a panic/recover pattern to stop code execution after the error occured and pass the error to the calling funtion
```
func bar() (err error){
  defer goerr.Recover(&err)

  err = foo()
  goerr.Fatal(err) //will panic on error

  //some awesome pieces of code
  //which are not executed when a error occurred
  //..
}
```
This pattern will not just handle errors produced through goerr.Fatal(). It will also handle memory corruptions, turns them into an error and add a stack trace to it. This can prevent you from extensive debugging sessions.

#### Logging
When you just want to log/report an error just do
```go
  err := foo()
  goerr.Log(err)
```
This is useful during development and mostly used at the top level of an application

# License
MIT
