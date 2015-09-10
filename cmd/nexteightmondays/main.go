package main

import (
	"fmt"
	"github.com/allonhadaya/boa/pp"
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"time"
)

const notFoundLayout = "Jan 2"

func main() {

	usage := `nexteightmondays is a command line client for setting acupuncture appointments.

Usage:
  nexteightmondays
  nexteightmondays --help

Options:
  -h --help    Show this screen.
  PP_USERNAME  Required environment variable containing the pocapoint username.
  PP_PASSWORD  Required environment variable containing the pocapoint password.`

	_, err := docopt.Parse(usage, nil, true, "nexteightmondays V1", false)
	if err != nil {
		log.Fatal(err)
	}

	username := os.Getenv("PP_USERNAME")
	if username == "" {
		log.Fatal("expected username to be defined in the environment variable PP_USERNAME")
	}

	password := os.Getenv("PP_PASSWORD")
	if password == "" {
		log.Fatal("expected password to be defined in the environment variable PP_PASSWORD")
	}

	session, err := pp.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	for _, day := range nextEightMondays() {

		block, err := session.GetBlock(day, pp.Elizabeth)
		if err != nil {
			log.Fatal(err)
		}

		for _, desired := range [...]time.Time{timeOnDay(day, "4:40pm"), timeOnDay(day, "4:50pm"), timeOnDay(day, "5:00pm")} {
			for _, appt := range block {
				if appt.Status == pp.Taken {
					continue
				}
				if appt.Timestamp == desired {
					if err := appt.Book(); err != nil {
						log.Fatalf("attempted booking (%s) but encountered error: %s", &appt, err)
					}
					log.Print(&appt)
					goto BOOKED
				}
			}
		}
		log.Printf("could not book an appointment for %s", day.Format(notFoundLayout))
	BOOKED:
	}
}

func nextEightMondays() []time.Time {

	// find next monday
	now := time.Now()
	y, m, d := now.Date()
	daysUntilMonday := int((7 + time.Monday - now.Weekday()) % 7)
	nextmonday := time.Date(y, m, d+daysUntilMonday, 0, 0, 0, 0, now.Location())

	// append seven weeks
	mondays := []time.Time{nextmonday}
	for i := 0; i < 7; i++ {
		mondays = append(mondays, mondays[i].AddDate(0, 0, 7))
	}
	return mondays
}

func timeOnDay(day time.Time, timeOfDay string) time.Time {
	const layout = "2006-1-2 3:04pm MST"
	y, m, d := day.Date()
	zone, _ := day.Zone()
	result, err := time.Parse(layout, fmt.Sprintf("%d-%d-%d %s %s", y, m, d, timeOfDay, zone))
	if err != nil {
		panic(err)
	}
	return result
}
