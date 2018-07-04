package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func requestGet(url string) (string, error) {
	log.Info("requestGet: ", url)
	resp, err := http.Get(url)
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

func getTmallCommentRequest(chnl chan string, url string, page int, totalPage int) {
	response, err := requestGet(url)
	if err != nil {
		log.Error(fmt.Sprintf("error in page: %d, %s", page, err))
		if page == totalPage {
			close(chnl)
		}
	}

	if strings.Contains(response, "rateList") {
		log.Info(fmt.Sprintf("ok in page: %d", page))
	}
	chnl <- strconv.Itoa(page)
	if page == totalPage {
		close(chnl)
	}
}

// 获取天猫评论
func getTmallComment(itemID string, sellerID string) {
	chnl := make(chan string)

	baseURL := "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&order=3&currentPage=%d&callback=jsonp"
	url := fmt.Sprintf(baseURL, itemID, sellerID, 1)
	log.Info(fmt.Sprintf("Tmall: %s", url))

	// 第一次请求获取评论页码数量
	response, err := requestGet(url)
	if err != nil {
		log.Error(fmt.Sprintf("get page error: %s", err))
		return
	}
	reg := regexp.MustCompile(`"page":(\d)}`)
	matches := reg.FindStringSubmatch(response)
	if len(matches) < 1 {
		log.Error(fmt.Sprintf("page 1 error: %s", response))
		return
	}
	totalPage, err := strconv.Atoi(matches[1])
	if err != nil {
	}

	for page := 1; page <= 5; page++ {
		go getTmallCommentRequest(chnl, fmt.Sprintf(baseURL, itemID, sellerID, page), page, totalPage)
	}

	for msg := range chnl {
		fmt.Println("Received ", msg)
	}

	fmt.Println("Finished")
}

func main() {
	itemID := "538232353890"
	sellerID := "1862759827"

	getTmallComment(itemID, sellerID)
}

// 	const json = `{"name":[{"wang": 123}, {"wang": 456}]}`

// 	value := gjson.Get(json, "name.#.wang")
// 	println(value.String())
// }
