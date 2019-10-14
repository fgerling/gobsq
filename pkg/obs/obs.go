package obs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseUrl string = "api.suse.de"

type Client struct {
	BaseURL    *url.URL
	username   string
	password   string
	httpClient *http.Client
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
