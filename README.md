# errz
Package errz provides simple error handling primitives to create more readable code and handle logging.

##### You should use errz when you want to..
  * ..easily detect & find memory leaks like a nil pointer dereference
  * ..give context to your errors trough stack traces
  * ..avoid the verbose but idiomatic error handling pattern

```go
err := someFunc()
if err != nil {
  log.Println(err)
  return err
}
```

##### When should i **not** use errz?
If you write algorithms, modules or libraries with high performance in mind.

-

The name of this package was inspired by [Marcel van Lohuizen](https://github.com/mpvl).  
It has a dependency on [github.com/pkg/errors ](https://github.com/pkg/errors)

## Using errz

#### Logging
When you just want to log/report an error just do

```go
err := foo()
errz.Log(err)
//or
errz.Log(err, "Additional error description")
```
This is useful during development and mostly used at the top level of an application

#### Return an error, add context, log and stop execution
A common use case is to return an error an stop further code execution in the currently executed function. To avoid the ideomatic error handling pattern we use panic/recover to stop code execution after the error occurred and pass the error to the calling function using named return values.

```go
func bar() (err error){
  defer errz.Recover(&err) //recover all panics

  err = foo()
  errz.Fatal(err) //panics on error

  //some awesome pieces of code
  //which is not executed when a error occurred at foo()
  //..
  return err
}
```
This pattern will not just handle errors produced through errz.Fatal(). It will also handle memory corruptions, turns them into an error, adds a stack trace to it and returns the error. This can prevent you from extensive debugging sessions.

# License
MIT
