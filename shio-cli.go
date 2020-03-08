package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

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

func main() {
	siteName := ""
	fmt.Println("Starting the application...")

	if strings.Compare(os.Args[1], "new") == 0 {
		xsrfToken := ""
		siteName = os.Args[2]

		config := ReadConfig()
		cookieJar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar: cookieJar,
		}
		req, err := http.NewRequest("GET", config.Server+"/api/v2", nil)
		req.SetBasicAuth(config.Login, config.Password)
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		} else {
			xsrfToken = getCookieByName(resp.Cookies(), "XSRF-TOKEN")
		}

		if resp.StatusCode == 200 {
			fmt.Println("Connected!")

		} else {
			fmt.Println("Authentication Failed!")
		}

		newSite(client, config, xsrfToken, siteName)
	}
}
func getCookieByName(cookies []*http.Cookie, name string) string {
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

func newSite(client *http.Client, config Config, xsrfToken string, siteName string) {
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
