package utils

import (
	"fmt"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

// Group - A organization, namespace
type Group struct {
	Slug string
}

// User -
type User struct {
	Username     string
	Password     string
	Groups       []string
	Addr         string
	CurrentGroup Group
	Conn         net.Conn
	Status       string
}

// Config -
type Config struct {
	Groups  []Group
	Users   []User
	Host    string
	Port    int
	WebPort int
}

// Options -
var Options Config

// LoadConfigFromFile -
func LoadConfigFromFile(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &Options)
	if err != nil {
		panic(err)
	}
	fmt.Println(Options)
	if Options.Host == "" {
		Options.Host = "0.0.0.0"
	}
	if Options.Port == 0 {
		Options.Port = 8962
	}
	if Options.WebPort == 0 {
		Options.WebPort = 8952
	}
}

// LoadConfigFromString -
func LoadConfigFromString(yamlConfig string) {
	Options = Config{}
	err := yaml.Unmarshal([]byte(yamlConfig), &Options)
	if err != nil {
		panic(err)
	}
}
