[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/benchkram/errz)
[![Go Report Card](https://goreportcard.com/badge/github.com/equanox/gotron)](https://goreportcard.com/report/github.com/equanox/gotron)
# errz
Package errz allows you to handle errors in one line of code.

##### You should use errz when you want to..
  * ..easily detect & find memory leaks like a nil pointer dereference
  * ..give context to your errors trough stack traces
  * ..avoid the verbose but idiomatic error handling pattern `if err != nil { return err }`


##### When should i **not** use errz?
If you write algorithms, modules or libraries with high performance in mind.

--

The name of this package was inspired by [Marcel van Lohuizen](https://github.com/mpvl).  
It has a dependency on [github.com/pkg/errors ](https://github.com/pkg/errors)

## Using errz
#### Logging
When you want to log/report an error.

```go
err := foo()
errz.Log(err)
```
This is useful during development and mostly used at the top level of an application

#### Return on error
It means stopping further code execution and returning to the calling function. panic/recover is used to stop code execution when a error occurred.
```go
func bar() (err error){
  defer errz.Recover(&err) //recover all panics

  err = foo()
  errz.Fatal(err) //panics on error

  //Some pieces of code.
  //Not executed when foo() returned a error
  //..
  //..

  return err
}
```
This pattern will not just handle errors. It will also handle memory corruptions, turns them into an error, adds a stack trace to it and returns the error. This can prevent you from extensive debugging sessions.

# License
MIT
