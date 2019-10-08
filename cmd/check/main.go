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

	dat, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		panic(err)
	}
	var conf config.Config

	_, err = toml.Decode(string(dat), &conf)
	if err != nil {
		panic(err)
	}

	var group = flag.String("group", "", "Set the group to search for.")
	var user = flag.String("user", "", "Set the obs user.")
	var password = flag.String("password", "", "Set the password.")
	flag.Parse()
	if *user == "" {
		user = &conf.Username
		log.Printf("User: %q\n", *user)
	}
	if *password == "" {
		password = &conf.Password
	}
	if *group == "" {
		group = &conf.Group
		log.Printf("Group: %q\n", *group)
	}
	var c obs.Collection
	c, err = obs.GetRRByGroup(*user, *password, *group)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Matches: %v\n", c.Matches)
	var when string
	for _, request := range c.ReleaseRequests {
		for _, review := range request.Reviews {
			if review.By_group == *group {
				when = review.When
			}
		}
		flag := ' '
		if request.Priority != "" {
			flag = '!'
		}
		fmt.Printf("%c %v %v\n", flag, when, obs.GetRepo(request))
	}
}
