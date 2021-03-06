package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fgerling/gobsq/pkg/config"

	"github.com/BurntSushi/toml"
	obs "github.com/fgerling/gobs"
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
	var when string

	for _, request := range rrs {
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
