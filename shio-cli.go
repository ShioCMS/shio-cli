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
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	Login    string
	Password string
	Server   string
}

// ReadConfig: Reads info from config file
func ReadConfig() Config {
	var configfile = "shio-cli.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	log.Print(config)
	return config
}

func main() {
	fmt.Println("Starting the application...")
	config := ReadConfig()
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Server+"/api/v2", nil)
	req.SetBasicAuth(config.Login, config.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//	bodyText, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println("Connected!")

	} else {
		fmt.Println("Authentication Failed!")
	}

	newSite()
}
func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}

func newSite() {
	token1 := ""
	config := ReadConfig()
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	req, err := http.NewRequest("GET", config.Server+"/api/v2", nil)
	req.SetBasicAuth(config.Login, config.Password)
	resp0, err0 := client.Do(req)

	if err0 == nil {
		for _, cookie := range resp0.Cookies() {
			if cookie.Name == "XSRF-TOKEN" {
				token1 = cookie.Value
			}
			fmt.Println("Found a cookie named:", cookie.Name)
			fmt.Println("Found a cookie value:", cookie.Value)
		}
	}

	jsonData1 := map[string]string{"objectType": "SITE", "name": "Teste10"}
	jsonValue1, _ := json.Marshal(jsonData1)
	req1, err1 := http.NewRequest("POST", config.Server+"/api/v2/site", bytes.NewBuffer(jsonValue1))

	fmt.Println("Err1: " + token1)

	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-XSRF-TOKEN", token1)
	resp1, err1 := client.Do(req1)

	for _, cookie := range req1.Cookies() {
		fmt.Println("Found a cookie2 named:", cookie.Name)
		fmt.Println("Found a cookie2 value:", cookie.Value)
	}
	if err1 != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Printf("Create site was executed: " + strconv.Itoa(resp1.StatusCode))
		data, _ := ioutil.ReadAll(resp1.Body)
		fmt.Println("Body" + string(data))
	}

}
