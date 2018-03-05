package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"github.com/rjhorniii/ics-golang"
	"strings"
)

func main() {

	var dead bool
	var sched bool
	var active bool
	var inactive bool
	
	// define flags
	dupsPtr := flag.String( "d", "", "Filename for duplicate removal")
	appPtr := flag.String("a", "", "Filename to append new events")
	outPtr := flag.String("o", "", "Filename for event output, default stdout")
	flag.BoolVar(&sched, "scheduled", false, "Event time should be scheduled")
	flag.BoolVar(&dead, "deadline", false, "Event time should be scheduled")
	flag.BoolVar(&active, "active", true, "Event time should be scheduled")
	flag.BoolVar(&inactive, "inactive", false, "Event time should be scheduled")
	
	// parse flags and arguments
	flag.Parse()
	
	if len(flag.Args()) == 0 {
		fmt.Println("At least one input argument is required.\n")
		return
	}

	if *appPtr != "" && *outPtr != *outPtr {
		fmt.Println("can't have both output and append files")
		os.Exit(1)
	}
	// Collect duplicate IDs before parsing inputs
	var dupIDs string
	if *dupsPtr != "" {
		dupIDs = dups( *dupsPtr)
	}
	
	//  create new parser
	parser := ics.New()

	// get the input chan
	inputChan := parser.GetInputChan()


	// send referenced arguments
	for _, url := range flag.Args() {
		inputChan <- url
	}

	//  wait for the calendar to be parsed
	parser.Wait()

	// get all calendars in this parser
	cal, err := parser.GetCalendars()

	//  check for errors
	if err == nil {
		// set output file when there are events
		var f *os.File
		if *outPtr != "" {
			f, err = os.OpenFile(*outPtr, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
		} else
		if *appPtr != "" {
			f, err = os.OpenFile(*appPtr, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			f = os.Stdout
		}

		for _, calendar := range cal {
			allEvents := calendar.GetEventsByDates()
			for _, event := range allEvents {
				// eliminate duplicates
				if (event[0].GetID() == dupIDs) {
					continue
				}
				
				// print the event
				// choose active or inactive timestamp
				format :=  "* %s <%s>\n"
				switch {
				case inactive:
					format = "* %s [%s]\n"
				case active:
					format = "* %s <%s>\n"
				}					

				fmt.Fprintf(f, format, strings.Replace(event[0].GetSummary(), `\,`, ",", -1), event[0].GetStart().Format("2006-01-02 15:04"))
				// Scheduled, Deadline, or nothing depending upon switches
				switch {
					case dead:
					fmt.Fprintf(f, "    DEADLINE: <%s-%s>\n", event[0].GetStart().Format("2006-01-02 15:04"), event[0].GetEnd().Format("15:04"))
					case sched:
					fmt.Fprintf(f, "    SCHEDULED: <%s-%s>\n", event[0].GetStart().Format("2006-01-02 15:04"), event[0].GetEnd().Format("15:04"))
					default:
				}
				// Print drawer contents
				fmt.Fprintln(f, "  :ICALCONTENTS:")
				fmt.Fprintf(f, "  :ORGUID: %s\n", event[0].GetID())
				fmt.Fprintf(f, "  :ORIGINAL-UID: %s\n", event[0].GetImportedID())
				fmt.Fprintf(f, "  :DTSTART: %s\n", event[0].GetStart().Format("2006-01-02 15:04"))
				fmt.Fprintf(f, "  :DTEND: %s\n", event[0].GetEnd().Format("2006-01-02 15:04"))
				fmt.Fprintf(f, "  :DTSTAMP: %s\n", event[0].GetDTStamp().Format("2006-01-02 15:04"))
				for _, attendee := range event[0].GetAttendees() {
					fmt.Fprintf(f, "  :ATTENDEE: %v\n", attendee)
				}
				fmt.Fprintf(f, "  :ORGANIZER: %s\n", event[0].GetOrganizer())
				if event[0].GetGeo() != nil {
					fmt.Fprintf(f, "  :GEO: %v, \n", event[0].GetGeo())
				}
				tzids := ""
				for _, tz := range event[0].GetDTZID() {
					if !strings.Contains(tzids, tz) {
						tzids = tzids + tz
					}
				}
				fmt.Fprintf(f, "  :TZIDS: %s\n", tzids)
				fmt.Fprintln(f, "  :END:")
				// Print Description and location
				
				fmt.Fprintln(f, "** Description\n")
				for _, line := range strings.Split(event[0].GetDescription(), `\n`) {
					fmt.Fprintf(f, "  %s\n", strings.Replace(line, `\,`, ",", -1)) //remove escape from commas (a CSV thing)
				}
				if event[0].GetLocation() != "" {
					fmt.Fprintf(f, "** Location %s \n", event[0].GetLocation())
				}
			}
		}
	} else {
		// error
		fmt.Fprintln(os.Stderr, err)
	}

}

func dups( dup string) string {
	return ("Duplicate processing not implemented yet")
}
