//go:build windows || plan9

package main

import (
	"io"
	"log"
	"net/http"
	urlutil "net/url"
	"os"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

func GetLoginUrl() (string, string) {
	pinger, err := ping.NewPinger("202.114.0.131")
	if err != nil {
		log.Fatal("Error ", err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Duration(3 * time.Second)
	pinger.SetPrivileged(true)
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Fatal("Error ", err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	if stats.PacketLoss < 100.0 {
		log.Println("The network is connected, no authentication required")
		os.Exit(0)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("http://123.123.123.123")
	if err != nil {
		log.Fatal("Error ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error ", err)
	}
	res := string(body)
	url := strings.Split(res, "'")[1]
	queryString := urlutil.QueryEscape(strings.Split(url, "?")[1])
	return url, queryString
}

func GetCookie(url string) *http.Cookie {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error ", err)
	}
	defer resp.Body.Close()
	cookie := resp.Cookies()[0]
	return cookie
}

func Login(url string, queryString string, username string, password string, service string, cookie *http.Cookie) string {
	trueurl := strings.Split(url, "/eportal/")[0] + "/eportal/InterFace.do?method=login"

	client := &http.Client{}
	var req *http.Request
	data := "userId=" + username +
		"&password=" + password +
		"&service=" + service +
		"&queryString=" + queryString +
		"&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false"
	req, _ = http.NewRequest("POST", trueurl, strings.NewReader(data))
	req.AddCookie(cookie)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error ", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func RegisterMAC(url string, userIndex string, cookie *http.Cookie) string {
	trueurl := strings.Split(url, "/eportal/")[0] + "/eportal/InterFace.do?method=registerMac"
	client := &http.Client{}
	var req *http.Request
	data := "mac=&userIndex=" + userIndex
	req, _ = http.NewRequest("POST", trueurl, strings.NewReader(data))
	req.AddCookie(cookie)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error ", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func main() {
	cmdargs := os.Args
	argslen := len(cmdargs)
	if argslen != 4 && argslen != 5 || argslen == 5 && cmdargs[4] != "auto" {
		log.Fatal("Usage: " + cmdargs[0] + " <username> <password> <internet|local> [auto]")
	}
	username := cmdargs[1]
	password := cmdargs[2]
	service := cmdargs[3]

	if service != "internet" && service != "local" {
		log.Fatal("Please use legal service\n", "Usage: "+cmdargs[0]+" <username> <password> <internet|local> [auto]")
	}

	url, queryString := GetLoginUrl()
	cookie := GetCookie(url)
	res := Login(url, queryString, username, password, service, cookie)
	if len(strings.Split(res, "\"result\":\"success\"")) == 2 {
		log.Println("Login success!")
	} else {
		log.Fatal("Login fail!\n", res)
	}

	if argslen == 5 {
		userIndex := strings.Split(res, "\"")[3]
		res := RegisterMAC(url, userIndex, cookie)
		log.Println(res)
	}
}
