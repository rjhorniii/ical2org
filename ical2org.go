package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/PuloV/ics-golang"
)

func main() {

	//  create new parser
	parser := ics.New()

	// get the input chan
	inputChan := parser.GetInputChan()

	// send referenced arguments
	for _, url := range os.Args[1:] {
		inputChan <- url
	}

	//  wait for the calendar to be parsed
	parser.Wait()

	// get all calendars in this parser
	cal, err := parser.GetCalendars()

	//  check for errors
	if err == nil {

	   for _, calendar := range cal {

		allEvents := calendar.GetEventsByDates()
	       for _, event := range allEvents {

		// print the event
		       fmt.Printf("* %s <%s>\n", strings.Replace( event[0].GetSummary(), `\,`, ",", -1), event[0].GetStart().Format("2006-01-02 15:04"))
		       fmt.Printf( "    SCHEDULED: <%s-%s>\n", event[0].GetStart().Format("2006-01-02 15:04"), event[0].GetEnd().Format("15:04"))
		       fmt.Println("  :PROPERTIES:")
		       fmt.Printf("  :UID: %s\n", event[0].GetID())
		       fmt.Printf("  :ORIGINAL-UID: %s\n", event[0].GetImportedID())
		       fmt.Printf("  :DTSTART: %s\n", event[0].GetStart().Format("2006-01-02 15:04"))
		       fmt.Printf("  :DTEND: %s\n", event[0].GetEnd().Format("2006-01-02 15:04"))
		       fmt.Printf("  :DTSTAMP: %s\n", event[0].GetDTStamp().Format("2006-01-02 15:04"))
		       for _, attendee := range event[0].GetAttendees() {
			       fmt.Printf("  :ATTENDEE: %v\n", attendee)
		       }
		       fmt.Printf("  :ORGANIZER: %s\n", event[0].GetOrganizer())
		       if event[0].GetGeo() != nil {
			       fmt.Printf("  :GEO: %v, \n", event[0].GetGeo())
		       }
		       if event[0].GetLocation() != "" {
			       fmt.Printf("  :LOCATION: %s \n", event[0].GetLocation())
		       }
		       tzids := "" 
		       for _, tz := range event[0].GetDTZID() {
			       if( !strings.Contains( tzids, tz)) {
				       tzids = tzids + tz
			       }
		       }
		       fmt.Printf("  :TZIDS: %s\n", tzids)
		       fmt.Println("  :END:")
		       fmt.Println("** Description\n")
		       for _, line := range strings.Split(event[0].GetDescription(), `\n`) {
			       fmt.Printf("  %s\n", strings.Replace( line, `\,`, ",", -1))  //remove escape from commas (a CSV thing)
		       }
		   }
	    }				
	} else {
		// error
		fmt.Println(err)
	}

}
