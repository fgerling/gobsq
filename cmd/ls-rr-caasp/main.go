package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"gitlab.suse.de/fgerling/qam-caasp-concourse-resource/pkg/config"
	"gitlab.suse.de/fgerling/qam-caasp-concourse-resource/pkg/obs"
	"io/ioutil"
	"log"
)

func main() {
	var config_file = flag.String("conf", "", "Set the config file.")
	var group = flag.String("group", "", "Set the group to search for.")
	var user = flag.String("user", "", "Set the obs user.")
	var password = flag.String("password", "", "Set the password.")

	flag.Parse()
	if *config_file == "" {
		*config_file = "./config.toml"
		log.Printf("Config file: %q\n", *config_file)
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
	log.Printf("User: %q\n", *user)
	if *password == "" {
		password = &conf.Password
	}
	if *group == "" {
		group = &conf.Group
	}
	log.Printf("Group: %q\n", *group)

	var rrs []obs.ReleaseRequest
	client := obs.NewClient(*user, *password)
	rrs, err = client.GetReleaseRequests(*group, "new,review")
	if err != nil {
		log.Fatal(err)
	}

	for _, request := range rrs {
		flag := ' '
		if request.Priority != "" {
			flag = '!'
		}
		requestLink := fmt.Sprintf("https://maintenance.suse.de/request/%v", request.Id)
		patchinfo, err := client.GetPatchinfo(request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%c %v (%v) \n", flag, requestLink, patchinfo.Summary)
	}
}
