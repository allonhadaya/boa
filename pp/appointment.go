package pp

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Appointment struct {
	// There's an interesting quirk here: the unix timestamp is always one hour
	// earlier than the appointment date. This is most likely an implementation
	// bug related to timezones.
	timestamp    time.Time
	practitioner int
}

func (a Appointment) String() string {
	const layout = "on Jan 2 at 3:04 pm"
	return fmt.Sprintf(
		"meeting with %d %s",
		a.practitioner,
		a.timestamp.Format(layout))
}

// ListAppointments fetches and parses the last 100 appointments
// as seen on the user's "My Account" page.
func (s *Session) List() ([]Appointment, error) {

	body, err := s.loadApps()
	if err != nil {
		return nil, err
	}

	divs, err := parseDivs(body)
	if err != nil {
		return nil, err
	}

	return parseApps(divs)
}

func (s *Session) loadApps() ([]byte, error) {

	v := url.Values{}
	v.Set("start", "0")
	v.Set("limit", "100")
	v.Set("destinationurl", "apps_pat.php")

	resp, err := (*s).client.PostForm("http://pocapoint.com/pp/_request/index.php", v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

var divsPath *xpath.Expression = xpath.Compile("//body/div")

func parseDivs(body []byte) ([]xml.Node, error) {

	root, err := gokogiri.ParseHtml(body)
	if err != nil {
		return nil, err
	}

	return root.Root().Search(divsPath)
}

func parseApps(divs []xml.Node) ([]Appointment, error) {

	apps := make([]Appointment, len(divs))
	for i, div := range divs {

		app, err := parseId(div.Attr("id"))
		if err != nil {
			return nil, err
		}

		apps[i] = *app
	}

	return apps, nil
}

var idPattern *regexp.Regexp = regexp.MustCompile("^app-([0-9]+)-([0-9]+)$")

func parseId(id string) (*Appointment, error) {

	values := idPattern.FindStringSubmatch(id)

	ts, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp := time.Unix(ts, 0)

	practitioner, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Appointment{timestamp: timestamp, practitioner: int(practitioner)}, nil
}

// Book books
func (s *Session) Book(a Appointment) error {
	//         'yyyy-mm-dd hh:MM:SS'         10                75 (liz)    1438265400 (seconds)
	//                 v                      v                   v          v
	//         x - date_time - locationid - duration - doubl - pracid - pp - ut
	//         ^                   ^                     ^               ^
	//   must be scraped           25                    0              boa

	// SETAPP:
	// GET /pp/boa/yyyy-mm-dd
	// parse token with id: app-{ut}-{pracid}
	// look for a reserve button
	//   no? return err already booked
	//   yes. scrape out the x:
	//     ./[Text contains YES]/a/get onclick text .. regex: '\((\d+),' ... assumed first arg. fuck this programming
	//     POST /_request/index.php
	//     ContentType: x-formurl-encode or some shit
	//     x={x}&date_time:{^^}
	return nil
}

func (a Appointment) pathDate() string {
	return a.timestamp.Format("2006-01-02")
}

func (a Appointment) htmlId() string {
	return fmt.Sprintf("app-%d-%d", a.timestamp.Unix(), a.practitioner)
}
