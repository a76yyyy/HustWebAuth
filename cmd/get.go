/*
Copyright © 2022 a76yyyy q981331502@163.com

*/
package cmd

import (
	"io"
	"log"
	"net/http"
	urlutil "net/url"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the login url from the redirect url",
	Long:  `If the specified IP fails to be pinged for more than the specified counts, get the login_url from the redirect_url`,
	Run: func(cmd *cobra.Command, args []string) {
		url, queryString, connected, err := GetLoginUrl()
		if err != nil {
			log.Fatal(err.Error())
		}
		if connected {
			log.Println("The network is connected, no authentication required")
		} else {
			log.Println("The login url is: ", url)
			log.Println("The query string is: ", queryString)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetLoginUrl() (string, string, bool, error) {
	pinger, err := ping.NewPinger(pingIP)
	if err != nil {
		return "", "", false, err
	}
	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.SetPrivileged(pingPrivilege)
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return "", "", false, err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	if stats.PacketLoss < 100.0 {
		return "", "", true, nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(redirectURL)
	if err != nil {
		return "", "", false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", "", false, err
	}
	res := string(body)
	url := strings.Split(res, "'")[1]
	queryString := urlutil.QueryEscape(strings.Split(url, "?")[1])
	return url, queryString, false, nil
}
