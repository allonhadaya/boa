package pp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

var (
	appointmentTaken = errors.New("this appointment cannot be booked because it is taken")
)

// Book books
func (a *Appointment) Book() error {
	switch a.Status {
	case Booked:
		return nil
	case Taken:
		return appointmentTaken
	case Available:
		return a.reallyBook()
	default:
		panic(fmt.Errorf("Undefined booking behavior for status: %s", a.Status))
	}
}

func (a *Appointment) reallyBook() error {

	v := url.Values{}
	v.Set("x", a.blockIndex)
	v.Set("date_time", a.makeAppDateFormat())
	v.Set("locationid", "25")
	v.Set("duration", "10")
	v.Set("doubl", "0")
	v.Set("pracid", strconv.Itoa(int(a.Practitioner)))
	v.Set("pp", "boa")
	v.Set("ut", strconv.FormatInt(a.Timestamp.Unix(), 10))
	v.Set("destinationurl", "app.php")

	resp, err := a.session.client.PostForm("http://pocapoint.com/pp/_request/index.php", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.Status = Booked

	return nil
}

func (a *Appointment) makeAppDateFormat() string {
	const layout = "2006-01-02 15:04:05"
	return a.Timestamp.Format(layout)
}
