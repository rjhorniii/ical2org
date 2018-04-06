package main

import (
	//	"fmt"
	// "github.com/rjhorniii/ics-golang"
	//	"github.com/davecgh/go-spew/spew"
	"testing"
)


func TestMain(t *testing.T) {

	var a args
	
	a.outfile = "tests/xx91596.org"
	a.args = append( a.args, "tests/xx91596.ics")

	process(a)

}
