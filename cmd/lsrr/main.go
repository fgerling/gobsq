package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
	obs "github.com/fgerling/gobs"
	"github.com/fgerling/gobsq/pkg/config"
)

func main() {
	var config_file = flag.String("conf", "", "Set the config file.")
	var group = flag.String("group", "", "Set the group to search for.")
	var user = flag.String("user", "", "Set the obs user.")
	var password = flag.String("password", "", "Set the password.")

	flag.Parse()
	if *config_file == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("%+v", err)
			panic(err)
		}
		*config_file = filepath.Join(home, ".gobs.toml")
	}
	var conf config.Config
	dat, err := ioutil.ReadFile(*config_file)
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
	}
	tmpl, err := template.New("list-requests").Parse("{{if eq .Request.Priority \"important\"}}!{{else}} {{end}} https://maintenance.suse.de/request/{{.Request.Id}} ({{.Summary}})\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, request := range rrs {
		patchinfo, err := client.GetPatchinfo(request)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, TmplStruct{Request: request, Summary: patchinfo.Summary})
		if err != nil {
			log.Fatal(err)
		}
	}
}
