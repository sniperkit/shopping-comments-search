package main

import (
	"errors"
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
	retry int
	proxy string
	redisClient RedisClient
	useRandomProxy bool
}

// Request class
type Request struct {
	client  *http.Client
	transport *http.Transport
	options RequestOptions
}

// NewRequest create the Request
func NewRequest(options RequestOptions) Request {
	tr := &http.Transport{}

	fmt.Println(options.proxy)
	if options.proxy != "" {
		proxyURL, _ := url.Parse(options.proxy)
		tr.Proxy = http.ProxyURL(proxyURL)
		fmt.Println("===")
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	_client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 600,
	}

	return Request{options: options, transport: tr, client: _client}
}

func (request *Request) get (url string) (string, error) {
	log.Info("requestGet: ", url)

	// TODO: 验证是否能在这里切换代理

	if request.options.useRandomProxy {
		request.changeProxy()
	}

	fmt.Println(url)
	resp, err := request.client.Get(url)
	if err != nil { return "", err}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {return "", err}

	return string(body), nil
}

func (request *Request) changeProxy() {
	if request.options.useRandomProxy {
		proxy, err := request.getRandomProxy()
		if err != nil {}

		request.transport.Proxy = http.ProxyURL(proxy)
	}
}

func (request *Request) getRandomProxy() (*url.URL, error) {
	if !request.options.redisClient.isConnected() {
		return &url.URL{}, errors.New("redis not connectetd")
	}

	proxy, err := request.options.redisClient.srandmember(request.options.redisClient.KeyProxy)
	if err!= nil {}

	fmt.Println(proxy)
	proxyURL, _ := url.Parse(proxy)
	return proxyURL, nil
}
