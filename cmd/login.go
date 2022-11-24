/*
Copyright © 2022 a76yyyy q981331502@163.com
*/

package cmd

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var register bool

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Hust web auth only once",
	Long:  `Hust web auth only once.`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := Login()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Cmd.PersistentFlags().String("foo", "", "A help for foo")
	loginCmd.PersistentFlags().BoolVarP(&register, "register", "r", false, "Register Mac address")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Get cookie of the auth page
func GetCookie(url string) (*http.Cookie, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cookie := resp.Cookies()[0]
	return cookie, err
}

// Login to auth the network
func login(url string, queryString string, account string, password string, service string, cookie *http.Cookie) (string, error) {
	trueurl := strings.Split(url, "/eportal/")[0] + "/eportal/InterFace.do?method=login"

	client := &http.Client{}
	var req *http.Request
	data := "userId=" + account +
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
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), err
}

// RegisterMAC register the mac address, only for the first time
func RegisterMAC(url string, userIndex string, cookie *http.Cookie) (string, error) {
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
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// Hust web auth once.
func Login() (res string, err error) {
	url, queryString, connected, err := GetLoginUrl()
	if err != nil {
		return "", err
	}
	if connected {
		if logConnected {
			return "The network is connected, no authentication required", nil
		}
		return "", nil
	}

	cookie, err := GetCookie(url)
	if err != nil {
		return "", err
	}

	res, err = login(url, queryString, account, password, service, cookie)
	if err != nil {
		return "", err
	}
	if len(strings.Split(res, "\"result\":\"success\"")) == 2 {
		res = "Login success!"
	} else {
		return "", errors.New("Login fail: " + res)
	}

	if register {
		userIndex := strings.Split(res, "\"")[3]
		res, err := RegisterMAC(url, userIndex, cookie)
		if err != nil {
			register = false
			return "", err
		}
		return res, nil
	}
	return res, nil
}
