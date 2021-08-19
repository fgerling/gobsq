package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	obs "github.com/fgerling/gobs"
	"github.com/fgerling/gobsq/pkg/config"
)

func main() {
	var configFile = flag.String("conf", "", "Set the config file.")
	var group = flag.String("group", "", "Set the group to search for.")
	var user = flag.String("user", "", "Set the obs user.")
	var password = flag.String("password", "", "Set the password.")

	flag.Parse()
	if *configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("%+v", err)
			panic(err)
		}
		*configFile = filepath.Join(home, ".gobs.toml")
	}
	var conf config.Config
	dat, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("%+v", err)
	} else {
		_, err = toml.Decode(string(dat), &conf)
		if err != nil {
			panic(err)
		}
	}
	if *user == "" {
		user = &conf.Username
	}
	if *password == "" {
		password = &conf.Password
	}
	if *group == "" {
		group = &conf.Group
	}

	var rrs []obs.ReleaseRequest
	client := obs.NewClient(*user, *password)
	rrs, err = client.GetReleaseRequests(*group, "new,review")
	if err != nil {
		log.Fatal(err)
	}
	type TmplStruct struct {
		Request obs.ReleaseRequest
		Summary string
		Src     string
	}
	tmpl, err := template.New("list-requests").Parse("{{.Src}}:{{.Request.Id}} ({{.Summary}})\n    https://maintenance.suse.de/request/{{.Request.Id}}\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, request := range rrs {
		patchinfo, err := client.GetPatchinfo(request)
		if err != nil {
			log.Fatal(err)
		}
		srcPrjStr := strings.Replace(request.Actions[0].Source.Project, "SUSE", "S", -1)
		srcPrjStr = strings.Replace(srcPrjStr, "Maintenance", "M", -1)
		err = tmpl.Execute(os.Stdout, TmplStruct{Request: request, Summary: patchinfo.Summary, Src: srcPrjStr})
		if err != nil {
			log.Fatal(err)
		}
	}
}
