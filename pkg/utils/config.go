package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Group - A organization, namespace
type Group struct {
	Slug string
}

// User -
type User struct {
	Username string
	Password string
	Groups   []string
}

// Config -
type Config struct {
	Groups []Group
	Users  []User
	Port   int
	Host   string
}

// Options -
var Options Config

// LoadConfig -
func LoadConfig() {
	Options = Config{}

	dat, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	yamlConfig := string(dat)
	fmt.Println(yamlConfig)

	err = yaml.Unmarshal(dat, &Options)
	if err != nil {
		panic(err)
	}
}
