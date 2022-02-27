package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlutil "net/url"
	"os"
	"strings"
)

func GetLoginUrl() (string, string) {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	res := string(body)
	url := strings.Split(res, "'")[1]
	queryString := urlutil.QueryEscape(strings.Split(url, "?")[1])
	return url, queryString
}

func GetCookie(url string) *http.Cookie {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	cmdargs := os.Args
	argslen := len(cmdargs)
	if argslen != 4 && argslen != 5 || argslen == 5 && cmdargs[4] != "auto" {
		fmt.Println("Usage: " + cmdargs[0] + " <username> <password> <internet|local> [auto]")
		os.Exit(1)
	}
	username := cmdargs[1]
	password := cmdargs[2]
	service := cmdargs[3]

	if service != "internet" && service != "local" {
		fmt.Println("Please use legal service")
		fmt.Println("Usage: " + cmdargs[0] + " <username> <password> <internet|local> [auto]")
		os.Exit(1)
	}

	url, queryString := GetLoginUrl()
	cookie := GetCookie(url)
	res := Login(url, queryString, username, password, service, cookie)

	if len(strings.Split(res, "\"result\":\"success\"")) == 2 {
		fmt.Println("Login success!")
	} else {
		fmt.Println("Login fail!")
		fmt.Println(res)
		os.Exit(1)
	}

	if argslen == 5 {
		userIndex := strings.Split(res, "\"")[3]
		res := RegisterMAC(url, userIndex, cookie)
		fmt.Println(res)
	}
}
