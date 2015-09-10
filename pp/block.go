package pp

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type Block []Appointment

// paths maps the different kind of appointments to xpaths that select their respective divs.
// They must be formatted with practitioner before being valid xpath expressions.
var paths = map[Status]string{
	Available: "//div[substring(@id, string-length(@id) - string-length('%[1]d') +1) = '%[1]d']/a[text() = 'Reserve']/..",
	Booked:    "//div[substring(@id, string-length(@id) - string-length('%[1]d') +1) = '%[1]d']/a[text() = 'Cancel']/..",
	Taken:     "//div[substring(@id, string-length(@id) - string-length('%[1]d') +1) = '%[1]d' and not(a)]",
}

// GetBlock GETs and parses practitioner's appointments on date
// along with any information needed to book available appointments.
func (s *Session) GetBlock(date time.Time, practitioner Practitioner) (Block, error) {

	root, err := s.loadBlock(date)
	if err != nil {
		return nil, err
	}

	var result Block

	for status, path := range paths {

		divs, err := root.Search(xpath.Compile(fmt.Sprintf(path, practitioner)))
		if err != nil {
			return nil, err
		}

		for _, div := range divs {

			timestamp, blockIndex, err := parseAppDiv(div)
			if err != nil {
				return nil, err
			}

			result = append(result, Appointment{
				session:      s,
				Timestamp:    time.Unix(timestamp, 0),
				Practitioner: practitioner,
				Status:       status,
				blockIndex:   blockIndex,
			})
		}
	}
	return result, nil
}

// GET and parse the schedule for the given date.
func (s *Session) loadBlock(date time.Time) (*html.HtmlDocument, error) {

	address := fmt.Sprintf("http://pocapoint.com/pp/boa/%s/25", date.Format("2006-01-02"))
	resp, err := s.client.Get(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	root, err := gokogiri.ParseHtml(body)
	if err != nil {
		return nil, err
	}

	return root, nil
}

var (
	idBlockPattern    *regexp.Regexp = regexp.MustCompile("^app-([0-9]+)-[0-9]+$")
	blockIndexPattern *regexp.Regexp = regexp.MustCompile("makeApp\\(([0-9]+)")
)

// parseAppDiv extracts timestamp and blockindex from an appointment div
func parseAppDiv(div xml.Node) (timestamp int64, blockIndex string, err error) {

	idValues := idBlockPattern.FindStringSubmatch(div.Attr("id"))
	timestamp, err = strconv.ParseInt(idValues[1], 10, 64)
	if err != nil {
		return
	}

	blockIndexValues := blockIndexPattern.FindStringSubmatch(div.Content())
	if len(blockIndexValues) == 1 {
		blockIndex = blockIndexValues[0]
	}
	return
}
