package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/allonhadaya/boa/book"
	"github.com/allonhadaya/boa/pp"
	"github.com/docopt/docopt-go"
	"os"
	"time"
)

const argDateLayout string = "2006-01-02"

func main() {

	usage := `book sets acupuncture appointments based on a calendar.

Usage:
  book [--start=<start>] [--end=<end>]
  book -h | --help
  book --version

Options:
  -h --help        Show this screen.
  --version        Show version.
  --start=<start>  The earliest date for which we should try to book an appointment in the format, yyyy-mm-dd [default: now].
  --end=<end>      The latest date for which we should try to book an appointment in the format, yyyy-mm-dd   [default: 8 weeks after start].
  PP_USERNAME      Required environment variable containing the pocapoint username.
  PP_PASSWORD      Required environment variable containing the pocapoint password.`

	args, err := docopt.Parse(usage, nil, true, "book V1", false)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	if startArg := args["--start"].(string); startArg != "now" {
		start, err = time.Parse(argDateLayout, startArg)
		if err != nil {
			log.
				WithError(err).
				WithField("--start", startArg).
				Fatal("could not parse argument")
		}
	}

	end := start.AddDate(0, 0, 8*7)
	if endArg := args["--end"].(string); endArg != "8 weeks after start" {
		end, err = time.Parse(argDateLayout, endArg)
		if err != nil {
			log.
				WithError(err).
				WithField("--end", endArg).
				Fatal("could not parse argument")
		}
	}

	if start.After(end) {
		log.
			WithField("start", start.String()).
			WithField("end", end.String()).
			Fatal("start cannot be after end")
	}

	username := os.Getenv("PP_USERNAME")
	if username == "" {
		log.Fatal("expected username to be defined in the environment variable PP_USERNAME")
	}

	password := os.Getenv("PP_PASSWORD")
	if password == "" {
		log.Fatal("expected password to be defined in the environment variable PP_PASSWORD")
	}

	log.
		WithField("start", start.String()).
		WithField("end", end.String()).
		WithField("username", username).
		Info("booking monday afternoons")

	session, err := pp.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	if booked, err := book.Book(mondayAfternoons{session}, start, end, session); err != nil {
		log.
			WithError(err).
			WithField("start", start.String()).
			WithField("end", end.String()).
			WithField("username", username).
			Fatal("Failed while booking monday afternoons")
	} else {
		for _, appointment := range booked {
			log.
				WithField("ts", appointment.Timestamp.String()).
				WithField("practitioner", appointment.Practitioner.String()).
				WithField("status", appointment.Status.String()).
				Info("booking appointment")
		}
	}
}

type mondayAfternoons struct {
	session *pp.Session
}

func (source mondayAfternoons) GetBetween(start, end time.Time) ([]book.Booking, error) {

	// find next monday
	y, m, d := start.Date()
	daysUntilMonday := int((7 + time.Monday - start.Weekday()) % 7)
	nextmonday := time.Date(y, m, d+daysUntilMonday, 0, 0, 0, 0, start.Location())

	// build a map of booked dates
	bookedSet := make(map[time.Time]bool)
	if booked, err := source.session.List(); err != nil {
		return nil, err
	} else {
		for _, appointment := range booked {
			y, m, d = appointment.Timestamp.Date()
			date := time.Date(y, m, d, 0, 0, 0, 0, start.Location())
			bookedSet[date] = true
		}
	}

	// add a week at a time
	var bookings []book.Booking
	for !nextmonday.After(end) {
		// take days where nothing has been booked yet
		if _, exists := bookedSet[nextmonday]; !exists {
			bookings = append(bookings, lizafternoon{nextmonday})
		}
		nextmonday = nextmonday.AddDate(0, 0, 7)
	}
	return bookings, nil
}

type lizafternoon struct {
	day time.Time
}

func (a lizafternoon) Range() (time.Time, time.Time) {

	at := func(timeOfDay string) time.Time {
		const layout = "2006-1-2 3:04pm MST"
		y, m, d := a.day.Date()
		zone, _ := a.day.Zone()
		result, err := time.Parse(layout, fmt.Sprintf("%d-%d-%d %s %s", y, m, d, timeOfDay, zone))
		if err != nil {
			panic(err)
		}
		return result
	}

	return at("4:40pm"), at("5:00pm")
}

func (_ lizafternoon) Practitioner() pp.Practitioner {
	return pp.Elizabeth
}

func (_ lizafternoon) SatisfyWith(appointment pp.Appointment) error {
	// nothing... unless you can get this booking to never get listed again
	return nil
}

func (_ lizafternoon) NotSatisfied() {
	// nothing
}
