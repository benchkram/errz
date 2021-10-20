// Package errz provides highlevel error handling helpers
package errz

import (
	"fmt"

	"github.com/pkg/errors"
)

// Fatal panics on error
// First parameter of msgs is used each following variadic arg is dropped
func Fatal(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		panic(errors.Wrap(err, str))
	}
}

// Recover recovers a panic introduced by Fatal, any other function which calls panics
// 				or a memory corruption. Logs the error when called without args.
//
// Must be used at the top of the function defered
// defer Recover()
// or
// defer Recover(&err)
func Recover(errs ...*error) {
	var e *error
	for _, err := range errs {
		e = err
		break
	}

	// handle panic
	if r := recover(); r != nil {
		var errmsg error
		// Preserve error which might have happend before panic/recover
		// Check if a error ptr was passed + a error occured
		if e != nil && *e != nil {
			// When error occured before panic then Wrap panic error around it
			errmsg = errors.Wrap(*e, r.(error).Error())
		} else {
			// No error occured just add a stacktrace
			errmsg = errors.Wrap(r.(error), "")
		}

		// If error cant't bubble up -> Log it
		if e != nil {
			*e = errmsg
		} else {
			log(errmsg)
		}
	}
}

// Log logs an error + stack trace directly to console or file
// Use this at the top level to publish errors without creating a new error object
func Log(err error, msgs ...string) {
	if err != nil {
		log(err, msgs...)
	}
}

// log splits the error in it's wrapped part and the original error message
// and writes it to the underlying logger.
func log(err error, msgs ...string) {
	var str string
	for _, msg := range msgs {
		str = msg
		break
	}
	sum := fmt.Errorf("%s", fmt.Sprintf("%+v", err))
	logger.Error(errors.Unwrap(err), str, logKeyStack, sum.Error())
}
