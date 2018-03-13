# ical2org-go
Convert a calendar in ICal format (e.g., .ics) into org-mode structure.


Usage: `ical2org [-d=<duplicates>] [-o=output] [-a=append]
       [--inactive] [--active]
       [--deadline] [--scheduled]
       [--repeats] [-dupinput] [-count]
       input files`

The input files can be URLs ("http://...."), local files, or stdin.  If the filename given is "-" then stdin is read.

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

### Dependencies

ical2org depends upon the forked library in
https://github.com/rjhorniii/ics-golang.  When the fork/pull request
are processed this will be changed.

### Duplicates Management

If the duplicates option is present and provides a file, that file
contains events that are already in org format.  Duplicates of these
events will not be generated as output.  The duplicates file is
scanned for ORGUIDs inside an ICALCONTENTS drawer.  When the input
files are processed, their event ORGUIDs are compared with the
ORGUIDs from the duplicates file.  Input events with matching ORGUIDs are
not output.

This matching rule permits you to decide whether modified events are
treated as duplicates or not.  If you remove the ORGUID from the
drawer, a modified event will be treated as different, and a new org
headline generated when that Ical event is processed again.  If the
ORGUID matches, no matter what other changes are made, a new org
headline will not be generated.

This means that by you can change task status to DONE, etc., and the
modified event will still be considered a duplicate.  You need not
ensure that all Ical sources are also updated when you change the org
file.  This allows later conversion of old email attachments
and calendars that contain historical events.  New events will not be
created in the org file for those duplicates.

The duplicates file can also be an orgmode file that is manually
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

The duplicates option also eliminates duplicates that may occur in the
inputs, so that if the same event is in multiple sources only one
event will be generated.  The ```dupinput``` flag can be used to
eliminate duplicates in the input files in situations where a
duplicates file is not used.

### Repeating events

There are good reasons to convert a repeating event into a single org
headline, and good reasons to replicate the repeating event as
multiple separate headlines.  The flag ```repeats``` is used to
control this.  If missing, or present with ```repeats=true``` multiple
org headlines will be generated.  If ```repeats=false```, only one
headline will be generated.  In either case, the ```ICALCONTENTS```
will contain the repeat rule, e.g., ```:RRULE:
FREQ=WEEKLY;UNTIL=20180325T035959Z;BYDAY=SA```.

The ORGUID will be the same for all of the headlines.  This enables
later efforts to find them if the repeating events need to be modified
or rescheduled.

### Time Zones

Org-mode does not use time zone tagged timestamps.  There are a
variety of good reasons for this.  As soon as you deal with
significant travel and teleconferences that originate in various parts
of the world you hit complex edge conditions.  The person involved can
establish the right thing to do fairly easily.  I've found no software
calendar that handles this properly.  For example, if I plan to be in
New York on Monday and Tuesday, Chicago on Wednesday, and Berlin on
Thursday and Friday, what timezone should be attached to which events?
I personally use the rule: local time there and then.  So I want
Monday and Tuesday to be America/NewYork, Wednesday to be
America/Chicago, and Thursday/Friday to be Europe/Berlin.  Note that I
did not use UTC offsets.  I want the times to be then local time.  I
don't want to worry about whether at that time and location it is
summer time or not.  I find org agendas are most useful this way, even
though different days and events are different time zones.

This program converts everything assuming that the event creator local
time zone is the one that matches my rule of "then/there time".  This
is most often, but not always, correct.  It is especially a problem for
teleconferences across time zones.  In the property ```TZIDS``` are
the creator specified time zone(s) for the event.  This allows a human
to know what was sent.  I then assume that between the description,
the zones, and the times, a person can decide whether and how to
adjust the timestamps in the org file.

### Systemd example

The following files specify using `ical2org` to fetch a calendar from
Google and update `events.org` at specified times.

Google-fetch.sh
```
#!/bin/bash
cd ~/org-directory
export https_proxy=https://192.168.1.1:8000
/home/rjhorniii/bin/ical2org -d=events-g.org -a=events-g.org https://calendar.google.com/google-stuff
```

google-fetch.service
```
[Unit]
Description=Fetch events from my google calendar into events-g.org

[Service]
Type=oneshot
ExecStart=/home/rjhorniii/org-directory/google-fetch.sh
```

google-fetch.timer, this example shows an irregular polling interval
```
[Unit]
Description=Fetch events timer

[Timer]
OnCalendar=*-*-* 6,8,10,12,13,14,16,18,23:10
Unit=google-fetch.service
Persistent=true

[Install]
WantedBy=timers.target
```

Put the `service` and `timer` in `~/.config/systemd/user`. To run it immediately
and then at the usual schedule use the commands:

```
systemctl --user daemon-reload
systemctl --user start google-fetch.timer
systemctl --user enable google-fetch.timer
```

### Mu4e integration

The following additions to your `.emacs` simplifies extracting schedule
from email.  In the email view, the keystrokes `A`,`s` will start
processing and ask for the attachment number to process.  The result
will be a buffer indicating the number of events captured.

```
;; define a pipe to parse appointments from mu4e

(defun process-ical-appointments (msg attachnum)
  "schedule appointments onto /home/rjhorniii/org/events-g.org"
  (mu4e-view-pipe-attachment msg attachnum
  	"/home/rjhorniii/bin/ical2org -count -d=/home/rjhorniii/org/events-g.org -a=/home/rjhorniii/org/events-g.org -"))
;; define 's' as the shortcut
(add-to-list 'mu4e-view-attachment-actions '("schedule appointment" . process-ical-appointments) t)
```

The `-count` indicates that the generated event count be sent to
stdout.  The file `events-g.org` will be updated. The use of `-`
indicates that input will be on stdin.  Mu4e deals with extracting the
attachment and sending it to the indicated command, and taking the
output and showing it as an emacs buffer.