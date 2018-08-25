//Package errz provides highlevel error handling helpers
//
//
package errz

import (
	"log"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/pkg/errors"
)

//errzLogger is a surrogate to logging.
//it enable using different logging packages.
type errzLogger struct {
	useZeroLog bool
}

func (el errzLogger) Printf(format string, v ...interface{}) {
	if el.useZeroLog {
		zlog.Error().Msgf(format, v)
		return
	}

	log.Printf(format, v)
}

var (
	logger errzLogger = errzLogger{useZeroLog: false}
)

func UseZeroLog(zl zerolog.Logger) {
	zlog.Logger = zl
	logger.useZeroLog = true
}

//Fatal panics on error
//First parameter of msgs is used each following variadic arg is dropped
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

//Recover recovers a panic introduced by Fatal, any other function which calls panics
//				or a memory corruption. Logs the error when called without args.
//
//Must be used at the top of the function defered
//defer Recover()
//or
//defer Recover(&err)
func Recover(errs ...*error) {
	var e *error
	for _, err := range errs {
		e = err
		break
	}

	//handle panic
	if r := recover(); r != nil {
		var errmsg error
		//Preserve error which might have happend before panic/recover
		//Check if a error ptr was passed + a error occured
		if e != nil && *e != nil {
			//When error occured before panic then Wrap panic error around it
			errmsg = errors.Wrap(*e, r.(error).Error())
		} else {
			//No error occured just add a stacktrace
			errmsg = errors.Wrap(r.(error), "")
		}

		//If error cant't bubble up -> Log it
		if e != nil {
			*e = errmsg
		} else {
			logger.Printf("%+v", errmsg)
		}
	}
}

//Log logs an error + stack trace directly to console or file
//Use this at the top level to publish errors without creating a new error object
func Log(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		logger.Printf("%+v", errors.Wrap(err, str))
	}
}
