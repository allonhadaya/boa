package pp

import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
)

// The maximum number of appointments to be loaded by the List function.
const appListLimit int = 100

// ListAppointments fetches and parses the last appListLimit appointments
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

	apps := make([]Appointment, len(divs))
	for i, div := range divs {
		app, err := parseListApp(div)
		if err != nil {
			return nil, err
		}
		apps[i] = *app
	}

	return apps, nil
}

func (s *Session) loadApps() ([]byte, error) {

	v := url.Values{}
	v.Set("start", "0")
	v.Set("limit", strconv.Itoa(appListLimit))
	v.Set("destinationurl", "apps_pat.php")

	resp, err := s.client.PostForm("http://pocapoint.com/pp/_request/index.php", v)
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

var idListPattern *regexp.Regexp = regexp.MustCompile("^app-([0-9]+)-([0-9]+)$")

func parseListApp(div xml.Node) (*Appointment, error) {

	values := idListPattern.FindStringSubmatch(div.Attr("id"))

	timestamp, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return nil, err
	}

	practitioner, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Appointment{
		session:      nil,
		pptimestamp:  timestamp,
		practitioner: Practitioner(practitioner),
		status:       Booked,
	}, nil
}
