package errz

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

var g_err error

const (
	tier1Error string = "Tier 1 error msg"
	tier2Error string = "Tier 2 error msg"
)

func Tier1() (err error) {
	defer Recover(&err)

	err = Tier2()
	Fatal(err)

	// Should never be reached when tier2 returns err
	err = errors.New(tier1Error)

	return
}

func Tier2() (err error) {
	err = errors.New(tier2Error)
	return
}

// TestMultiTierErrors test if tier2 error is returned (after multiple function calls)
func TestMultiTierErrors(t *testing.T) {
	err := Tier1()

	if errors.Cause(err).Error() != tier2Error {
		t.Errorf("Expecting Tier2 error")
	}
}

// Creates a memory corruption which is handled by recover
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

func CreateTier2NullPointerDereference() {
	var numberptr *int
	number := 1 + *numberptr

	//Dead Code
	fmt.Printf("This point must not be reached: %d", number)
}

// Creates a memory corruption which is handled by recover
func TestTier2NullPointerDereference(t *testing.T) {
	var err error
	defer func() {
		if err == nil {
			t.Errorf("No error occured, expecting : runtime error: invalid memory address or nil pointer dereference")
		}
	}() // Called after recover
	defer Recover(&err)

	// Creating memory corruption leading to panic
	CreateTier2NullPointerDereference()

	// Dead Code
	t.Errorf("This point must not be reached")
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
//just here to assure the api contract
func TestEmptyRecover(t *testing.T) {
	defer Recover()
}

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

//****** Benchmarks ********

func CheckLogRaw() {
	if g_err != nil {
		fmt.Println(g_err)
	}
}

func CheckLog() {
	Log(g_err)
}

func CheckFatal() {
	defer Recover(&g_err)
	Fatal(g_err)
}

func BenchmarkCheckLogRaw(b *testing.B) {
	g_err = nil
	for n := 0; n < b.N; n++ {
		CheckLogRaw()
	}
	g_err = nil
}

func BenchmarkCheckLog(b *testing.B) {
	g_err = nil
	for n := 0; n < b.N; n++ {
		CheckLog()
	}
	g_err = nil
}

func BenchmarkCheckFatal(b *testing.B) {
	g_err = nil
	for n := 0; n < b.N; n++ {
		CheckFatal()
	}
	g_err = nil
}
