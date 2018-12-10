package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rjhorniii/ics-golang"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type args struct {
	dupfile   string
	appfile   string
	outfile   string
	afterfile string
	label     string
	dupflag   bool
	sched     bool
	active    bool
	inactive  bool
	repeats   bool
	dead      bool
	count     bool
	after     bool
	args      []string
}

func main() {

	var a args
	var dupfilePtr *string
	var appPtr *string
	var outPtr *string
	var afterPtr *string
	var labelPtr *string

	// define flags
	dupfilePtr = flag.String("d", "", "Filename for duplicate removal")
	appPtr = flag.String("a", "", "Filename to append new events")
	outPtr = flag.String("o", "", "Filename for event output, default stdout")
	afterPtr = flag.String("after", "", "Only use events at and after this date")
	labelPtr = flag.String("label", "", "Label word in drawer to identify conversion")
	flag.BoolVar(&a.sched, "scheduled", false, "Event time should be scheduled")
	flag.BoolVar(&a.dead, "deadline", false, "Event time should be deadline")
	flag.BoolVar(&a.active, "active", false, "Headline timestamp should be active")
	flag.BoolVar(&a.inactive, "inactive", false, "Headline timestamp should be inactive")
	flag.BoolVar(&a.repeats, "repeats", true, "Generate an event per repeat")
	flag.BoolVar(&a.dupflag, "dupinput", false, "Do not generate duplicates from input")
	flag.BoolVar(&a.count, "count", false, "Report number of new events found on stdout")

	// parse flags and arguments
	flag.Parse()

	a.dupfile = *dupfilePtr
	a.appfile = *appPtr
	a.outfile = *outPtr
	a.afterfile = *afterPtr
	a.label = *labelPtr
	a.args = flag.Args()

	process(a)
}

func process(a args) {

	var after bool

	afterTime := time.Now()
	if a.afterfile != "" {
		if strings.HasPrefix(a.afterfile, "-") {
			d, err := time.ParseDuration(a.afterfile)
			if err != nil {
				fmt.Printf("After option, duration format error %s\n", err)
				return
			} else {
				afterTime = afterTime.Add(d)
				after = true
			}
		} else {
			t, err := time.Parse("2006-01-02", a.afterfile)
			if err != nil {
				fmt.Printf("After option, time format error %s\n", err)
				return
			} else {
				afterTime = t
				after = true
			}
		}
	}

	// Verify that input and output are present
	if len(a.args) == 0 {
		fmt.Println("At least one input argument is required.")
		return
	}

	if a.appfile != "" && a.outfile != "" {
		fmt.Println("can't have both output and append files")
		os.Exit(1)
	}
	ics.RepeatRuleApply = a.repeats

	// Collect duplicate IDs before parsing inputs

	dupIDs := map[string]bool{"": true}
	if a.dupfile != "" {
		dupIDs = dups(a.dupfile)
		a.dupflag = true
	}

	//  create new parser
	parser := ics.New()

	// get the input chan
	inputChan := parser.GetInputChan()

	// send referenced arguments
	for _, url := range a.args {
		// fmt.Printf( "send filename: %s\n", url)
		parser.Add() // another calendar to wait for
		inputChan <- url
	}

	close(inputChan)
	//  wait for the final calendar to be parsed
	parser.Wait()

	// get all calendars in this parser
	cal, err := parser.GetCalendars()
	eventsSaved := 0

	//  check for errors
	if err == nil {
		// set output file when there are events
		var f *os.File
		if a.outfile != "" {
			f, err = os.OpenFile(a.outfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
		} else if a.appfile != "" {
			f, err = os.OpenFile(a.appfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			f = os.Stdout
		}

		for _, calendar := range cal {
			allEvents := calendar.GetEvents()
			for _, event := range allEvents {
				// eliminate duplicates
				if dupIDs[event.GetID()] {
					continue
				}
				if a.dupflag {
					dupIDs[event.GetID()] = true
					// fmt.Printf(" duplicate ID added %s\n", event.GetID())
				}
				// eliminate events before after
				if after {
					if event.GetStart().Before(afterTime) {
						continue
					}
				}
				eventsSaved++
				// print the event
				// choose active or inactive timestamp
				summary := strings.Replace(event.GetSummary(), `\,`, ",", -1)
				start := event.GetStart().Format("2006-01-02 Mon 15:04")
				switch {
				case a.inactive:
					fmt.Fprintf(f, "* %s [%s]\n", summary, start)
				case a.active:
					fmt.Fprintf(f, "* %s <%s>\n", summary, start)
				default:
					fmt.Fprintf(f, "* %s\n", summary)
				}

				// Scheduled, Deadline, or nothing depending upon switches
				switch {
				case a.dead:
					fmt.Fprintf(f, "    DEADLINE: <%s-%s>\n", event.GetStart().Format("2006-01-02 Mon 15:04"), event.GetEnd().Format("15:04"))
				case a.sched:
					fmt.Fprintf(f, "    SCHEDULED: <%s-%s>\n", event.GetStart().Format("2006-01-02 Mon 15:04"), event.GetEnd().Format("15:04"))
				default:
				}
				// Print drawer contents
				fmt.Fprintln(f, "  :ICALCONTENTS:")
				fmt.Fprintf(f, "  :ORGUID: %s\n", event.GetID())
				if a.label != "" {
					fmt.Fprintf(f, "  :CONVERTLABEL: %s\n", a.label)
				}
				fmt.Fprintf(f, "  :ORIGINAL-UID: %s\n", event.GetImportedID())
				fmt.Fprintf(f, "  :DTSTART: %s\n", event.GetStart().Format("2006-01-02 15:04"))
				fmt.Fprintf(f, "  :DTEND: %s\n", event.GetEnd().Format("2006-01-02 15:04"))
				fmt.Fprintf(f, "  :DTSTAMP: %s\n", event.GetDTStamp().Format("2006-01-02 15:04"))
				for _, attendee := range event.GetAttendees() {
					fmt.Fprintf(f, "  :ATTENDEE: %v\n", attendee)
				}
				fmt.Fprintf(f, "  :ORGANIZER: %s\n", event.GetOrganizer())
				if event.GetGeo() != nil {
					fmt.Fprintf(f, "  :GEO: %v, \n", event.GetGeo())
				}
				tzids := ""
				for _, tz := range event.GetDTZID() {
					if !strings.Contains(tzids, tz) {
						tzids = tzids + tz
					}
				}
				fmt.Fprintf(f, "  :TZIDS: %s\n", tzids)
				fmt.Fprintf(f, "  :RRULE: %s\n", event.GetRRule())
				fmt.Fprintln(f, "  :END:")
				// Print Description and location

				fmt.Fprintf(f, "** Description\n\n")
				for _, line := range strings.Split(event.GetDescription(), `\n`) {
					fmt.Fprintf(f, "  %s\n", strings.Replace(line, `\,`, ",", -1)) //remove escape from commas (a CSV thing)
				}
				if event.GetLocation() != "" {
					fmt.Fprintf(f, "** Location %s \n", event.GetLocation())
				}
			}
		}
		if a.count {
			fmt.Fprintf(os.Stdout, " New events written: %d\n", eventsSaved)
			errors, _ := parser.GetErrors()
			if len(errors) != 0 {
				fmt.Printf("errors occurred %v\n", errors)
			}
		}
	} else {
		// error
		fmt.Fprintln(os.Stderr, err)
	}

}

/*  make this parallel later

type Duplicates struct {
	inputChan chan string
	outputChan chan map[string] bool
}
*/

func dups(dupname string) map[string]bool {
	// Basic state machine to find ORGID in ICALCONTENTS drawer.  It
	// accepts org-mode full syntax, but takes lots of shortcuts to combine
	// and ignore many irrelevant fields.  It will tolerate incorrect syntax,
	// although it might not get recognition right.

	// It uses a state machine that processes one line at a time.
	//
	// States are:
	//    Body - somewhere in body.  Only a headline will depart this state
	//    Head - somewhere in headline material.  Looking for drawers.
	//    Drawer - in a drawer that is not ICALCONTENTS.  Loofing for end.
	//    Contents - in the ICALCONTENTS drawer, looking for ORGID
	//
	const Body = 1
	const Head = 2
	const Drawer = 3
	const Contents = 4
	state := Body // State begins in Body

	// Pattern matches

	rHeadline, _ := regexp.Compile(`^\*`)                   // first character is "*"
	rBody, _ := regexp.Compile(`^[^:]*$`)                   // no colon anywhere on the line (also catches blank lines)
	rContents, _ := regexp.Compile(`^\s*:ICALCONTENTS:`)    //start of content drawer
	rOrgID, _ := regexp.Compile(`^\s*:ORGUID:\s*(\S*)\s*$`) // the orgID
	rOther, _ := regexp.Compile(`^\s*:\S*:`)                // start of another drawer
	rEnd, _ := regexp.Compile(`^\s*:END:`)                  // end of any drawer

	found := make(map[string]bool)

	// read lines until the end
	dupfile, err := os.Open(dupname)
	if err != nil {
		if os.IsNotExist(err) {
			return (found)
		} else {
			log.Fatal(err)
		}
	}
	defer dupfile.Close()
	scanner := bufio.NewScanner(dupfile)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case err != nil:
		case rHeadline.MatchString(line):
			state = Head
		case rBody.MatchString(line):
			state = Body
		case state == Body:
			// break.  Rest only apply if in header
		case state == Head && rContents.MatchString(line):
			state = Contents
		case state == Contents && rOrgID.MatchString(line):
			// extract UID add to map
			found[rOrgID.FindStringSubmatch(line)[1]] = true // extract the word after :ORGID:
		case state == Head && rOther.MatchString(line):
			state = Drawer
		case state == Drawer && rEnd.MatchString(line):
			state = Head
		case state == Contents && rEnd.MatchString(line):
			state = Head
		}
	}
	return (found)
}
