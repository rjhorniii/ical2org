package main

import (
	"fmt"
	"github.com/rjhorniii/ics-golang"
	//	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestMain(t *testing.T) {

	//  create new parser
	parser := ics.New()

	// get the input chan
	inputChan := parser.GetInputChan()

	//  send a local ics file
	inputChan <- "tests/xx91596.ics"

	//  wait for the calendar to be parsed
	parser.Wait()

	// get all calendars in this parser
	cal, err := parser.GetCalendars()

	//  check for errors
	if err == nil {

		for _, calendar := range cal {
			allEvents := calendar.GetEventsByDates()
			for i, event := range allEvents {

				// print the event
				fmt.Printf("%+v \n", i)
				fmt.Printf("%+v \n", event)
			}
			//		spew.Dump( allEvents)
		}
	} else {
		// error
		t.Fatal(err)
	}

}
