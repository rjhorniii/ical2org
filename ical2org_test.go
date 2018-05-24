package main

import (
	// "fmt"
	"io/ioutil"
	"os"
	// "github.com/rjhorniii/ics-golang"
	//	"github.com/davecgh/go-spew/spew"
	"testing"
	//"github.com/davecgh/go-spew/spew"
)

func TestMultiple(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", count: true, args: []string{"tests/xx91596.ics", "tests/test-vcal-3.vcs", "tests/wg-29.ics"}}

	process(a)
	// order is unpredicatable so no comparison
}

func TestX91596(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-correct output
	if compareFiles(a.outfile, "tests/xx91596.org-correct", t) == false {
		t.Fail()
	}
}
func TestDeadline(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", dead: true, args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-dead
	if compareFiles(a.outfile, "tests/xx91596.org-dead", t) == false {
		t.Fail()
	}
}
func TestSchedule(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", sched: true, args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-scheduled
	if compareFiles(a.outfile, "tests/xx91596.org-scheduled", t) == false {
		t.Fail()
	}
}
func TestActive(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", active: true, args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-correct
	if compareFiles(a.outfile, "tests/xx91596.org-correct", t) == false {
		t.Fail()
	}
}

func TestInactive(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", inactive: true, args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-inactive
	if compareFiles(a.outfile, "tests/xx91596.org-inactive", t) == false {
		t.Fail()
	}
}

func TestDupflag(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", dupflag: true, args: []string{"tests/xx91596.ics", "tests/xx91596a.ics"}}

	process(a)
	// compare with org-correct
	if compareFiles(a.outfile, "tests/xx91596.org-correct", t) == false {
		t.Fail()
	}
}

func TestDual(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", count: true, args: []string{"tests/xx91596.ics", "tests/xx91596a.ics"}}

	process(a)
	// compare with org-dual
	if compareFiles(a.outfile, "tests/xx91596.org-dual", t) == false {
		t.FailNow()
	}
}

func TestAfterDate(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", afterfile: "2030-01-01", args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with empty
	finfo, err := os.Stat(a.outfile)
	if err != nil {
		t.Fail()
	}
	if finfo.Size() != int64(0) {
		t.Fail()
	}
}
func TestAfterDuration(t *testing.T) {
	// remove tests/xx91596.org
	a := args{appfile: "tests/xx91596.org", afterfile: "-36h", args: []string{"tests/xx91596.ics"}}
	err := os.Remove(a.appfile)
	if !os.IsNotExist(err) && err != nil {
		t.Error(err)
	}
	process(a)
	// compare with org-correct
	if compareFiles(a.appfile, "tests/xx91596.org-correct", t) == false {
		t.Fail()
	}
	process(a)
	// compare with org-dual
	if compareFiles(a.appfile, "tests/xx91596.org-dual", t) == false {
		t.Fail()
	}
}

func TestLabel(t *testing.T) {
	a := args{outfile: "tests/xx91596.org", label: "test-label", args: []string{"tests/xx91596.ics"}}

	process(a)
	// compare with org-inactive
	if compareFiles(a.outfile, "tests/xx91596.org-labeled", t) == false {
		t.Fail()
	}
}
//
// file comparison function.  This assumes files are small and memory is large.
// It reads the whole files into memory and then compares.

func compareFiles(fname1 string, fname2 string, t *testing.T) bool {
	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f1, err1 := ioutil.ReadFile(fname1)

	if err1 != nil {
		t.Error(err1)
	}

	f2, err2 := ioutil.ReadFile(fname2)

	if err2 != nil {
		t.Error(err2)
	}

	str1 := string(f1)
	str2 := string(f2)

	if str1 != str2 {
		//		spew.Printf("test file: %v \n\ncomparison file %v\n", str1, str2)
		return false
	}
	return true
}
