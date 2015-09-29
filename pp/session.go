package pp

import (
	"errors"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	client http.Client
}

var badCredentials *xpath.Expression = xpath.Compile("//*[text()='Login Incorrect.']")

// New returns a new Session for the specified credentials.
// The session will hold on to some cookies for persistence.
func New(username, password string) (*Session, error) {

	jar, _ := cookiejar.New(nil)
	s := &Session{client: http.Client{Jar: jar}}

	v := url.Values{}
	v.Set("u", username)
	v.Set("p", password)
	v.Set("go", "login")

	resp, err := s.client.PostForm("http://pocapoint.com/pp/boa/", v)
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

	badNode, err := root.Root().Search(badCredentials)
	if err != nil {
		return nil, err
	}

	if badNode != nil {
		return nil, errors.New("Incorrect credentials")
	}

	// session cookies are in the jar
	return s, nil
}
