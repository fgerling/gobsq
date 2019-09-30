package obs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Target struct {
	Project string `xml:"project,attr"`
	Package string `xml:"package,attr"`
}
type Source struct {
	Project string `xml:"project,attr"`
	Package string `xml:"package,attr"`
}
type Action struct {
	Type   string `xml:"type,attr"`
	Source Source `xml:"source"`
	Target Target `xml:"target"`
}
type State struct {
	Name    string `xml:"name,attr"`
	Who     string `xml:"who,attr"`
	When    string `xml:"when,attr"`
	Comment string `xml:"comment"`
}
type Review struct {
	State    string `xml:"state,attr"`
	Who      string `xml:"who,attr"`
	When     string `xml:"when,attr"`
	By_user  string `xml:"by_user,attr"`
	By_group string `xml:"by_group,attr"`
	Comment  string `xml:"comment"`
}
type ReleaseRequest struct {
	Id          string   `xml:"id,attr"`
	Creator     string   `xml:"creator,attr"`
	Actions     []Action `xml:"action"`
	State       State    `xml:"state"`
	Review      Review   `xml:"review"`
	Description string   `xml:"description"`
}
type Collection struct {
	Matches         string           `xml:"matches,attr"`
	ReleaseRequests []ReleaseRequest `xml:"request"`
}

func GetRepo(rr ReleaseRequest) string {
	trgPrjStr := strings.Replace(rr.Actions[0].Target.Project, ":", "_", -1)
	srcPrjStr := strings.Replace(rr.Actions[0].Source.Project, ":", ":/", -1)
	repo := fmt.Sprintf("http://download.suse.de/ibs/%s/%s/%s.repo", srcPrjStr, trgPrjStr, rr.Actions[0].Source.Project)
	return repo
}

func GetRRByGroup(username string, password string, group string) (Collection, error) {
	v := Collection{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.suse.de/request", nil)
	req.SetBasicAuth(username, password)
	q := req.URL.Query()
	q.Add("view", "collection")
	q.Add("group", group)
	q.Add("states", "new,review")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return v, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return v, errors.New(fmt.Sprintf("Got status code: %v\n", resp.StatusCode))
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, err
	}

	err = xml.Unmarshal([]byte(bodyText), &v)
	if err != nil {
		return v, err
	}
	return v, err
}
