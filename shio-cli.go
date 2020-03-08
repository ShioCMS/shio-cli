package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/ShioCMS/shio-cli/site"
)

func main() {
	siteName := ""
	fmt.Println("Starting the application...")

	if strings.Compare(os.Args[1], "new") == 0 {
		xsrfToken := ""
		siteName = os.Args[2]

		config := site.ReadConfig()
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
			xsrfToken = site.GetCookieByName(resp.Cookies(), "XSRF-TOKEN")
		}

		if resp.StatusCode == 200 {
			fmt.Println("Connected!")

		} else {
			fmt.Println("Authentication Failed!")
		}

		site.NewSite(client, config, xsrfToken, siteName)
	}
}
