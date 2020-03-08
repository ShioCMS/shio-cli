package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alexcesaro/log/stdlog"
)

// Config ...
type Config struct {
	Login    string
	Password string
	Server   string
}

// ReadConfig reads info from config file
func ReadConfig() Config {
	logger := stdlog.GetFromFlags()
	var configfile = "shio-cli.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	logger.Debug(config)
	return config
}

// GetCookieByName ...
func GetCookieByName(cookies []*http.Cookie, name string) string {
	logger := stdlog.GetFromFlags()
	result := ""
	for _, cookie := range cookies {
		if cookie.Name == name {
			result = cookie.Value
		}
		logger.Debug("Found a cookie named:", cookie.Name)
		logger.Debug("Found a cookie value:", cookie.Value)
	}
	return result
}

// NewSite ...
func NewSite(client *http.Client, config Config, xsrfToken string, siteName string) {
	logger := stdlog.GetFromFlags()
	jsonData := map[string]string{"objectType": "SITE", "name": siteName}
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", config.Server+"/api/v2/site", bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-XSRF-TOKEN", xsrfToken)
	resp, err := client.Do(req)

	if err != nil && resp.StatusCode == 200 {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Printf("%s site was created.\n", siteName)
		data, _ := ioutil.ReadAll(resp.Body)
		logger.Debug("Body: " + string(data))
	}

}
