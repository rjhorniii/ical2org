# ical2org-go
Convert a calendar in ICal format (e.g., .ics) into org-mode structure.


Usage: `ical2org [-d=<duplicates>] [-o=output] [-a=append]
       [--inactive] [--active]
       [--deadline] [--scheduled] input files`

The input files can be either URLs ("http://....") or local files.

The resulting org formatted events will either
* replace the file specified with -o output,
* be appended to the file specified with -a output
* or be sent to stdout if -o and -a are not present.

Converted events have
* a headline with the Summary from the event with a timestamp from event start time.  It will be active by default.
The option `inactive` can be used to make it an inactive timestamp.
* an optional scheduling line: if `--scheduled` it will contain `SCHEDULED`; if `--deadline` it will contain `DEADLINE`. 
* a drawer ICALCONTENTS with potentially useful information extracted from the event
* a subheading Description with the description field, and
* a subheading Location with the Location information

Example output:

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


### Duplicates Management

If the duplicates option is present and provides a file, that
file contains events that are already in org format.  Duplicates of
these events should not be generated as output.  The duplicates file
is scanned for ORGUIDs inside an ICALCONTENTS drawer.  When the input
files are processed the generated event ORGUIDs are compared with these
ORGUIDs.  Events with matching ORGUIDs are not output.

The matching rule permits you to decide whether modified events are
treated as duplicates or not.  If you remove the ORGUID from the
drawer, a modified event will be treated as different, and a new org
headline generated.

ical2org does not compare event contents.  It only looks at the
ORGUIDs.

This means that by you can change task status to DONE, etc., and the
modified event will still be considered a duplicate.  This way you
need not ensure that all Ical sources are also updated when you change
the org file.  This should allow repeated conversion of old email
attachments and calendars that show historical events.  New events
will not be created in the org file.

The duplicates file can also be an orgmode file that is being manually
mainatined. In the manually maintained org drawer you list the ORGUIDs
that should be considered duplicate:

```
* Dummy event headline
  :ICALCONTENTS:
  :ORGUID: skdivndlkjf123
  :ORGUID: skdivndlkjf42354
  :ORGUID: skdivndlkjf987
  :END:
```

The duplicates processing permits a construction like:

```ical2org -d=events.org -a=events.org https://calendar.google.com/more-stuff```

to be run regularly, even hourly, without filling the `events.org`
file with duplicates.  The internal logic waits for duplicates file
processing to complete before it begins event generation, so it is
safe to have the output file and duplicates file be the same file.

ical2org depends upon the forked library in
https://github.com/rjhorniii/ics-golang.  Switch to the ical2org
branch.