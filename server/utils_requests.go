package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
)

// RequestOptions request options
type RequestOptions struct {
	retry          int
	proxy          string
	redisClient    RedisClient
	useRandomProxy bool
}

// Request class
type Request struct {
	client    *http.Client
	transport *http.Transport
	options   RequestOptions
}

// NewRequest create the Request
func NewRequest(options RequestOptions) Request {
	tr := &http.Transport{}

	if options.proxy != "" {
		proxyURL, _ := url.Parse(options.proxy)
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	_client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 300,
	}

	return Request{options: options, transport: tr, client: _client}
}

// 验证是否能在这里切换代理，毕竟是个指针，应该可以
// 自动重试，当返回200的时候就不重试了
func (request *Request) get(url string) (string, error) {
	log.Info("request.get: ", url)

	tries := 0
	var resp *http.Response
	var err error
	for {
		if tries > request.options.retry {
			return "", errors.New("retries count is max")
		}
		tries = tries + 1

		if request.options.useRandomProxy {
			log.Info("useRandomProxy")
			request.changeProxy()
		}
		resp, err = request.client.Get(url)
		if err != nil {
			log.Warning(err)
			log.Warning(resp)
			log.Info("retry request get")
			continue
		}
		defer resp.Body.Close()

		log.Info(resp.Status)
		if resp.StatusCode == 200 {
			break
		}
	}

	contentType := resp.Header.Get("Content-Type")
	utf8reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		return "", err
	}

	text, err := ioutil.ReadAll(utf8reader)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func (request *Request) changeProxy() {
	if request.options.useRandomProxy {
		proxy, err := request.getRandomProxy()
		if err != nil {
		}

		request.transport.Proxy = http.ProxyURL(proxy)
	}
}

func (request *Request) getRandomProxy() (*url.URL, error) {
	if !request.options.redisClient.isConnected() {
		return &url.URL{}, errors.New("redis not connectetd")
	}

	proxy := ""
	err := errors.New("getting proxy failed")
	count := 0
	rand.Seed(time.Now().Unix())
	for {
		if count > 10 {
			log.Error("can not get proxy from redis")
		}
		count = count + 1

		var proxies []string
		proxies, err = request.options.redisClient.zrange(request.options.redisClient.KeyProxy, int64(1000), int64(10000))
		if err != nil {
			log.Warn(err)
			continue
		}
		proxy = proxies[rand.Intn(len(proxies))]
		if strings.Contains(proxy, "socks") || strings.Contains(proxy, "https") {
			log.Warn(err)
			continue
		}
		log.Info(fmt.Sprintf("get proxy: %s", proxy))
		break
	}

	proxyURL, _ := url.Parse(proxy)
	return proxyURL, nil
}
