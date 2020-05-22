package store

import (
	"io/ioutil"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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
	Groups        []Group
	Users         []User
	Host          string
	Port          int
	JobTimeout    int    `yaml:"job_timeout"`
	WebPort       int    `yaml:"web_port"`
	DelayPolicy   string `yaml:"delay_policy"`
	LogLevel      string `yaml:"log_level"`
	TestMode      bool   `yaml:"test_mode"`
	AllowDispatch bool   `yaml:"allow_dispatch"`
}

// Options -
var Options Config

// LoadConfigFromFile -
func LoadConfigFromFile(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	LoadConfigFromString(string(data))
}

// LoadConfigFromString -
func LoadConfigFromString(yamlConfig string) {
	Options = Config{}
	err := yaml.Unmarshal([]byte(yamlConfig), &Options)
	if err != nil {
		panic(err)
	}
	if Options.Host == "" {
		Options.Host = "0.0.0.0"
	}
	if Options.Port == 0 {
		Options.Port = 8962
	}
	if Options.WebPort == 0 {
		Options.WebPort = 8952
	}
	if !Options.TestMode {
		Options.TestMode = os.Getenv("TEST_MODE") != ""
	}
	if !Options.AllowDispatch {
		Options.AllowDispatch = os.Getenv("ALLOW_DISPATCH") != ""
	}
	var logLevel = log.InfoLevel
	switch strings.ToLower(Options.LogLevel) {
	case "trace":
		logLevel = log.TraceLevel
	case "debug":
		logLevel = log.DebugLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "fatal":
		logLevel = log.FatalLevel
	case "panic":
		logLevel = log.PanicLevel
	}
	log.SetLevel(logLevel)
}
