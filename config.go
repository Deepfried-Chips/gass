package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
)

var config = Config{}

// Config holds configuration values.
type Config struct {
	Host            string
	Port            string
	PostgreLocation string
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
		fmt.Println(fmt.Errorf("Failed to read config file:	#%v ", err))
	}

	err = json.Unmarshal(jsonFile, c)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to Unmarshal: %v", err))
	}

	return c
}

func (c *Config) getPostgreConfig(location string) *sql.DB {
	db, err := sql.Open("postgres", location)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS USERS (USERNAME TEXT UNIQUE PRIMARY KEY, PASSWORD TEXT, ISADMIN BOOLEAN, MAXUPLOAD INT, ACTIVE BOOL)")
	if err != nil {
		return nil
	}
	return db
}

func (c *Config) getUserValidity(username string, db *sql.DB) bool {
	data, err := db.Query("SELECT * FROM USERS WHERE USERNAME = $1 AND ACTIVE = true", username)
	if err != nil {
		return false
	}
	if data != nil {
		err := data.Close()
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (c *Config) getUserAdmin(username string, db *sql.DB) bool {
	data, err := db.Query("SELECT isAdmin FROM USERS WHERE USERNAME = $1 AND ACTIVE = true", username)
	if err != nil {
		return false
	}
	defer func(data *sql.Rows) {
		err := data.Close()
		if err != nil {
		}
	}(data)
	for data.Next() {
		var isAdmin bool
		err = data.Scan(&isAdmin)
		if err != nil {
			return false
		}
		return isAdmin
	}
	return false
}

func (c *Config) getPassHash(username string, db *sql.DB) string {
	data, err := db.Query("SELECT * FROM USERS WHERE USERNAME = $1 AND ACTIVE = true", username)
	if err != nil {
		return ""
	}
	defer func(data *sql.Rows) {
		err := data.Close()
		if err != nil {

		}
	}(data)
	for data.Next() {
		var password string
		err = data.Scan(&password)
		if err != nil {
			return ""
		}
		return password
	}
	return ""
}

func (c *Config) getMaxUpload(username string, db *sql.DB) int {
	data, err := db.Query("SELECT * FROM USERS WHERE USERNAME = $1 AND ACTIVE = true", username)
	if err != nil {
		return 0
	}
	defer func(data *sql.Rows) {
		err := data.Close()
		if err != nil {

		}
	}(data)
	for data.Next() {
		var maxUpload int
		err = data.Scan(&maxUpload)
		if err != nil {
			return 0
		}
		return maxUpload
	}
	return 0
}
