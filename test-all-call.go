package main

import (
	"fmt"
	"github.com/PuloV/ics-golang"
	//	"github.com/davecgh/go-spew/spew"
)

func main() {

	//  create new parser
	parser := ics.New()

	// get the input chan
	inputChan := parser.GetInputChan()

	//  send a local ics file
	inputChan <- "data/xx03614.ics"

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
		fmt.Println(err)
	}

}
