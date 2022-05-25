package main

import (
	json "encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var config = Config{}

// Config holds configuration values.
type Config struct {
	Host  string
	Port  string
	Users []user
}

type user struct {
	Username  string
	Password  string
	isAdmin   bool
	MaxUpload int64
	FileTypes []string
}

// getConf reads config values from a file.
func (c *Config) getConf() *Config {

	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panic("Failed to read config file:	#%v ", err)
	}

	err = json.Unmarshal(jsonFile, c)
	if err != nil {
		log.Fatalf("failed to Unmarshal: %v", err)
	}

	return c
}

func (c *Config) getbyUsername(username string) *user {
	for _, u := range c.Users {
		if u.Username == username {
			return &u
		}
	}
	return nil
}
