package pp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

var (
	appointmentTaken = errors.New("This appointment cannot be booked.")
)

func (a *Appointment) IsTaken() bool {
	return a.status == Taken
}

// Book books
func (a *Appointment) Book() error {
	switch a.status {
	case Booked:
		return nil
	case Taken:
		return appointmentTaken
	case Available:
		return a.reallyBook()
	default:
		panic(fmt.Errorf("Undefined booking behavior for status: %s", a.status))
	}
}

func (a *Appointment) reallyBook() error {

	v := url.Values{}
	v.Set("x", a.blockIndex)
	v.Set("date_time", a.makeAppDateFormat())
	v.Set("locationid", "25")
	v.Set("duration", "10")
	v.Set("doubl", "0")
	v.Set("pracid", strconv.Itoa(int(a.practitioner)))
	v.Set("pp", "boa")
	v.Set("ut", strconv.FormatInt(a.pptimestamp, 10))
	v.Set("destinationurl", "app.php")

	resp, err := a.session.client.PostForm("http://pocapoint.com/pp/_request/index.php", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// assume success because the server does not implement meaningful status codes

	a.status = Booked

	return nil
}

func (a *Appointment) makeAppDateFormat() string {
	const layout = "2006-01-02 15:04:05"
	return a.RealTimestamp().Format(layout)
}
