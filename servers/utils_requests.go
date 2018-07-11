package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// RequestOptions request options 
type RequestOptions struct {
	proxy string
}

// NewClient can create the http request client
func NewClient(options RequestOptions) *http.Client {
	tr := &http.Transport{}

	fmt.Println(options.proxy)
	if options.proxy != "" {
		proxyURL, _ := url.Parse(options.proxy)
		tr.Proxy = http.ProxyURL(proxyURL)
		fmt.Println("===")
	}

	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Println(tr)

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 600,
	}

	return client
}

func requestGet(url string, options RequestOptions) (string, error) {
	log.Info("requestGet: ", url)

	client := NewClient(options)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
