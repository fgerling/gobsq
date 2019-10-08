package obs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Person struct {
	Name string `xml:"name,attr"`
	Role string `xml:"role,attr"` //optional
}
type Grouped struct {
	Id string `xml:"id,attr"`
}
type Acceptinfo struct {
	Rev      string `xml:"rev,attr"`
	Srcmd5   string `xml:"srcmd5,attr"`
	Osrcmd5  string `xml:"osrcmd5,attr"`
	Oproject string `xml:"oproject,attr"` //optional
	Opackage string `xml:"opackage,attr"` //optional
	Xsrcmd5  string `xml:"xsrcmd5,attr"`  //optional
	Oxsrcmd5 string `xml:"oxsrcmd5,attr"` //optional

}
type Options struct {
	Sourceupdate    string `xml:"sourceupdate"`    //optional
	Updatelink      string `xml:"updatelink"`      //optional
	Makeoriginolder string `xml:"makeoriginolder"` //optional
}
type Group struct {
	Name string `xml:"name,attr"`
	Role string `xml:"role,attr"` //optional
}
type Target struct {
	Project        string `xml:"project,attr"`
	Package        string `xml:"package,attr"`        //optional
	Releaseproject string `xml:"releaseproject,attr"` //optional
	Repository     string `xml:"repository,attr"`     //optional
}
type Source struct {
	Project string `xml:"project,attr"`
	Package string `xml:"package,attr"` //optional
	Rev     string `xml:"rev,attr"`     //optional
}
type Action struct {
	Type       string     `xml:"type,attr"`
	Source     Source     `xml:"source"`     //optional
	Target     Target     `xml:"target"`     //optional
	Person     Person     `xml:"person"`     //optional
	Group      Group      `xml:"group"`      //optional
	Grouped    []Grouped  `xml:"grouped"`    //optional-oneOrMore
	Options    Options    `xml:"options"`    //optional
	Acceptinfo Acceptinfo `xml:"acceptinfo"` //optional
}
type State struct {
	Name    string `xml:"name,attr"`
	Who     string `xml:"who,attr"`
	When    string `xml:"when,attr"`
	Comment string `xml:"comment"`
}
type Review struct {
	State      string `xml:"state,attr"`
	By_user    string `xml:"by_user,attr"`
	By_group   string `xml:"by_group,attr"`
	By_project string `xml:"by_project,attr"`
	By_package string `xml:"by_package,attr"`
	Who        string `xml:"who,attr"`
	When       string `xml:"when,attr"`
	Comment    string `xml:"comment"`
}
type History struct {
	Who         string `xml:"who,attr"`
	When        string `xml:"when,attr"`
	Description string `xml:"description"`
	Comment     string `xml:"comment"` //optional
}
type ReleaseRequest struct {
	Id          string    `xml:"id,attr"`      //optional
	Creator     string    `xml:"creator,attr"` //optional
	Actions     []Action  `xml:"action"`       //oneOrMore
	State       State     `xml:"state"`        //optional
	Description string    `xml:"description"`  //optional
	Priority    string    `xml:"priority"`     //optional ref:obs-ratings
	Reviews     []Review  `xml:"review"`       //zeroOrMore
	Histories   []History `xml:"history"`      //zeroOrMore
	Title       string    `xml:"title"`        //optional
	Accept_at   string    `xml:"accept_at"`    //optional
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
