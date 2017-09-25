package goerr

import (
	"log"

	"github.com/pkg/errors"
)

//Error Custom error type to handle errors gracefully
//
//After stashing one error this acts as a no-op.
//This asure that a function returns the first occured error
//in all circumstances.
type Error struct {
	err error
	Ret *error
}

//Stash with capacity 1
func (ew *Error) Stash(err error, msgs ...string) {
	if ew.err != nil {
		return
	}
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		ew.err = errors.Wrap(err, str)
	}
}

//StashFatal Panics when error occures
//Only first param of msgs is used each following variadic arg is dropped
func (ew *Error) StashFatal(err error, msgs ...string) {
	if ew.err != nil {
		return
	}

	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		ew.err = errors.Wrap(err, str)
		panic(ew.err)
	}
}

//ReportError logs a error directly to console or file
//Use this at the top level to bublish errors
func ReportError(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		log.Printf("%+v", errors.Wrap(err, str))
	}
}

//Recover the first error
//If unhandled errors occure
//Must be used at the top of the function
//defer ew.Recover()
func (ew *Error) Recover() {
	if r := recover(); r != nil {
		if ew.err == nil {
			ew.err = errors.Wrap(
				r.(error),
				"Recovered a unhandled error, most likely a memory corruption")
		}
		if ew.Ret != nil {
			*ew.Ret = errors.Wrap(ew.err, "")
		} else {
			ReportError(ew.err) // Report when error can't bubble up
		}
	}
}
