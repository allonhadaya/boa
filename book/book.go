package book

import (
	"github.com/allonhadaya/boa/pp"
	"time"
)

type Source interface {
	GetBetween(start, end time.Time) ([]Booking, error)
}

type Booking interface {
	Practitioner() pp.Practitioner
	Range() (start, end time.Time)
	SatisfyWith(appointment pp.Appointment) error
	NotSatisfied()
}

func Book(source Source, start, end time.Time, session *pp.Session) (booked []pp.Appointment, err error) {

	bookings, err := source.GetBetween(start, end)
	if err != nil {
		return
	}

	for _, booking := range bookings {

		practitioner := booking.Practitioner()
		start, end = booking.Range()

		blockDomain := daysBetween(start, end)
		targets := appointmentTimes(start, end)

		for _, day := range blockDomain {

			block, blockerr := session.GetBlock(day, practitioner)
			if blockerr != nil {
				err = blockerr
				return
			}

			for _, target := range targets {
				for _, appointment := range block {

					if appointment.Status == pp.Taken {
						continue
					}

					if appointment.Timestamp != target {
						continue
					}

					booked = append(booked, appointment)

					// try to book it
					if err = appointment.Book(); err != nil {
						return
					}

					// try to mark it as booked
					if err = booking.SatisfyWith(appointment); err != nil {
						// consider wrapping err to point out that the systems are out of sync now
						// which may lead to double booking in the future
						return
					}

					goto BOOKED
				}
			}
		}
		booking.NotSatisfied()
	BOOKED:
	}

	return
}

// daysBetween enumerates the inclusive dates between start and end
func daysBetween(start, end time.Time) []time.Time {

	// same day at midnight, local time.
	dateOf := func(value time.Time) time.Time {
		y, m, d := value.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, value.Location())
	}

	start = dateOf(start)
	end = dateOf(end)

	var days []time.Time
	for next := start; next.Before(end); next = next.AddDate(0, 0, 1) {
		days = append(days, next)
	}
	days = append(days, end)
	return days
}

// appointmentTimes enumerates the ten minute intervals between start and end
func appointmentTimes(start, end time.Time) []time.Time {
	var times []time.Time
	for next := start; start.Before(end); start = start.Add(time.Minute * 10) {
		times = append(times, next)
	}
	times = append(times, end)
	return times
}
