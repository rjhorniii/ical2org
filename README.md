# ical2org-go
Convert a calendar in ICal format (e.g., .ics) into org-mode structure.


Usage: `ical2org [-d=<duplicates>] [-o=output] input files`

The input files can be either URLs ("http://....") or local files.

The resulting org formatted events will either
* replace the file specified with -o output,
* be appended to the file specified with -a output
* or be sent to stdout.

Converted events have
* a headline with the Summary from the event,
* a scheduling line
* a drawer ICALCONTENTS with potentially useful information extracted from the event
* a subheading Description with the description field, and
* a subheading Location with the Location information

```
* DECISION MEETING: ITI Planning Monthly Call (11am CT, 12pm ET, 6pm CET) Host key: xxxxxx <2027-10-19 11:00>
    SCHEDULED: <2027-10-19 11:00-12:30>
  :ICALCONTENTS:
  :ORGUID: 1e23c027675290d3a25f3f8441bd91b2
  :ORIGINAL-UID: 040000008200E00074C5B7101A82E00800000000C0DE4F324F47D301000000000000000010000000ED14C4A947673341969029C8BE8EDA02
  :DTSTART: 2027-10-19 11:00
  :DTEND: 2027-10-19 12:30
  :DTSTAMP: 2017-10-17 19:04
  :ORGANIZER: <nil>
  :TZIDS: TZID="Central Standard Time":
  :END:
** Description

  Jeremiah Myers invites you to an online meeting using WebEx.
  
  Meeting Number: ccc ccc ccc
  Meeting Password: xxxxxxxxx
  
  -------------------------------------------------------
  To join this meeting
  -------------------------------------------------------
  1. Go to https://himss.webex.com/himss/j.php?MTID=m3d541
  ...

```


*future* If the duplicates option is present and provides a file, that
file contains events that are already in org format.  Duplicates of
these events should not be generated as output.  The duplicates file
is scanned for ORGUIDs inside an ICALCONTENTS drawer.  When the input
files are processed the generated event ORGUIDs are compared with these
ORGUIDs.  Events with matching ORGUIDs are not output.

The matching rule permits you to decide whether modified events are
treated as duplicates or not.  If you remove the ORGUID from the
drawer, a modified event will be treated as different.  ical2org does
not attempt to match on event contents.  It only looks at the ORGUIDs.

The duplicates file can be an orgmode events file that is being
mainatined. It can also be a manually maintained org drawer:

```
* Dummy event headline
  :ICALCONTENTS:
  :ORGUID: skdivndlkjf123
  :ORGUID: skdivndlkjf42354
  :ORGUID: skdivndlkjf987
  :END:
```

This permits a construction like:

`ical2org -d=events.org https://calendar.google.com/more-stuff  >>events.org`

That can be run regularly, even hourly, without filling the events.org
file with duplicates.  The internal logic waits for duplicate
processing to complete before it begins event generation, so it is
safe to have the output file and duplicates file be the same file.

ical2org depends upon the forked library in
https://github.com/rjhorniii/ics-golang.  Switch to the ical2org
branch.