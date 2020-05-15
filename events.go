package main

import (
	"context"
	"time"

	"google.golang.org/api/calendar/v3"
)

const timeFormat = "2006-01-02"

func processEvents(opts options) {
	l := func(s string, args ...interface{}) {
		if opts.debug != nil {
			opts.debug.Printf(s, args...)
		}
	}

	l("Create calendar client")
	srv, err := calendar.New(opts.client)
	if err != nil {
		opts.stderr.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	start, end, err := formatTimes(opts)
	if err != nil {
		opts.stderr.Fatalf("Error parsing start and/or end time: %w", err)
	}

	l("Getting all calendar events from %s to %s", start, end)
	eventList := srv.Events.List(opts.CalendarID).
		ShowDeleted(false).SingleEvents(true).OrderBy("startTime").TimeMin(start).TimeMax(end)

	eventList.Pages(context.Background(), func(events *calendar.Events) error {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}

			opts.stdout.Printf("Deleting %v (%v) ...\n", item.Summary, date)
			if !opts.DryRun {
				start := time.Now()
				if err := srv.Events.Delete(opts.CalendarID, item.Id).Do(); err != nil {
					opts.stderr.Printf("Error deleting event: %s", err.Error())
					return err
				}
				end := time.Now()

				// The Google Calendar API has rate limiting. One second per deletion
				// seems to keep the app under the radar. Subtract the run time from
				// one second and sleep for that duration. If the delete took longer
				// than a second, sleep for a whole second to slow things down and
				// stay under the radar.
				execTime := end.Sub(start)
				remainder := time.Duration(time.Second) - execTime
				if remainder < 0 {
					remainder = time.Duration(time.Second)
				}

				time.Sleep(remainder)
			}
		}

		return nil
	})

	_, err = eventList.Do()
	if err != nil {
		opts.stderr.Fatalf("Error deleting events: %w\n\n", err)
	}
}

func formatTimes(opts options) (string, string, error) {
	var st, et time.Time
	if opts.StartDate == "" {
		st = time.Unix(0, 0)
	} else {
		tmp, err := time.Parse(timeFormat, opts.StartDate)
		if err != nil {
			return "", "", err
		}

		st = time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 999999999, time.Local)
	}

	if opts.EndDate == "" {
		now := time.Now()
		et = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.Local)
	} else {
		tmp, err := time.Parse(timeFormat, opts.EndDate)
		if err != nil {
			return "", "", err
		}

		et = time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 999999999, time.Local)
	}

	return st.Format(time.RFC3339), et.Format(time.RFC3339), nil
}
