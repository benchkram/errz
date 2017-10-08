package goerr

import (
	"testing"

	"github.com/pkg/errors"
)

const (
	errorMsg   string = "An intended error occured"
	tier1Error string = "Tier 1 error msg"
	tier2Error string = "Tier 2 error msg"
)

func Tier1() (err error) {
	defer Recover(&err)

	err = Tier2()
	Fatal(err)

	//Should never be reached when tier2 returns err
	err = errors.New(tier1Error)

	return
}

func Tier2() (err error) {
	err = errors.New(tier2Error)
	return
}

//TestMultiTierErrors test if original error after multiple function calls is returned
func TestMultiTierErrors(t *testing.T) {
	err := Tier1()

	if errors.Cause(err).Error() != tier2Error {
		t.Errorf("Expecting Tier2 error")
	}
}

//Creates a memory corruption which is handled by recover
func TestNullPointerDereference(t *testing.T) {
	var err error
	defer func() {
		if err == nil {
			t.Errorf("No error occured, expecting : runtime error: invalid memory address or nil pointer dereference")
		}
	}() //Called after recover
	defer Recover(&err)

	//Creating memory corruption leading to panic
	var numberptr *int
	number := 1 + *numberptr

	//Dead Code
	t.Errorf("This point must not be reached: %d", number)
}

//Test handling od a common error
func TestCommonError(t *testing.T) {
	var err error
	defer func() {
		if err == nil {
			t.Errorf("No error was thrown")
		}
	}()
	defer Recover(&err)

	err = errors.New("Thrown a error")
}

//Cannot fail, must not print anything to console
func TestEmptyRecover(t *testing.T) {
	defer Recover()
}

// //This should report a error directly to console
// func TestEmptyRecoverNilPointerDereference(t *testing.T) {
// 	defer Recover()a
//
// 	var numberptr *int
// 	number := 1 + *numberptr
//
// 	//Dead Code
// 	t.Errorf("This point must not be reached: %d", number)
// }

func TestFatal(t *testing.T) {
	var err error
	defer func() {
		if err == nil {
			t.Errorf("Error did not occure")
		}
	}()
	defer Recover(&err)

	err = errors.New("Thrown a error")
	Fatal(err, "Panic on error")

	t.Errorf("This point must not be reached")
}
